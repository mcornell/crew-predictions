import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { When, Then } = createBdd();

let lastPredictionStatus = 0;

When('I enter a home score of {int} and away score of {int} for the first match', async ({ page }, home: number, away: number) => {
  const card = page.locator('[data-testid="match-card"]').first();
  await card.locator('input[name="home_goals"]').fill(String(home));
  await card.locator('input[name="away_goals"]').fill(String(away));
});

When('I click {string}', async ({ page }, label: string) => {
  await page.getByRole('button', { name: label }).first().click();
});

Then('I should see my prediction of {string} on the first match card', async ({ page }, score: string) => {
  const card = page.locator('[data-testid="match-card"]').first();
  await expect(card.locator('[data-testid="matchup"]').getByText(score)).toBeVisible();
});

When('I submit a prediction via API for a match that has already kicked off', async ({ page }) => {
  const matchResp = await page.request.get('/api/matches');
  const data = await matchResp.json();
  const completed = data.matches.find((m: any) =>
    m.status !== 'STATUS_SCHEDULED' && m.status !== 'STATUS_IN_PROGRESS'
  );
  if (!completed) throw new Error('No completed match found to test locking against');
  const resp = await page.request.post('/api/predictions', {
    form: { match_id: completed.id, home_goals: '2', away_goals: '1' },
  });
  lastPredictionStatus = resp.status();
});

Then('the server should reject it with 403', async () => {
  expect(lastPredictionStatus).toBe(403);
});
