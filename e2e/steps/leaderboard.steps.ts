import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Given, When, Then } = createBdd();

Given('{string} predicted {int}-{int} for match {string}', async ({ request }, handle: string, home: number, away: number, matchID: string) => {
  await request.post('/predictions', {
    form: { match_id: matchID, home_goals: String(home), away_goals: String(away) },
    headers: {
      Cookie: `session=${Buffer.from(JSON.stringify({ handle })).toString('base64')}`,
    },
  });
});

Given('the final score for match {string} was {int}-{int}', async ({ request }, matchID: string, home: number, away: number) => {
  await request.post('/admin/results', {
    form: { match_id: matchID, home_goals: String(home), away_goals: String(away) },
  });
});

When('I visit the leaderboard', async ({ page }) => {
  await page.goto('/leaderboard');
});

Then('I should see {string} with {int} points', async ({ page }, handle: string, points: number) => {
  const row = page.locator('[data-testid="leaderboard-row"]').filter({ hasText: handle });
  await expect(row.getByText(String(points))).toBeVisible();
});
