import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Given, When, Then } = createBdd();

let lastPredictionStatus = 0;

function kickoffForStatus(status: string): string {
  const d = new Date();
  const isScheduled = status === 'STATUS_SCHEDULED' || status === 'STATUS_IN_PROGRESS';
  d.setHours(d.getHours() + (isScheduled ? 24 : -24));
  return d.toISOString();
}

Given('the following matches are seeded:', async ({ request }, table: any) => {
  for (const row of table.hashes()) {
    await request.post('/admin/seed-match', {
      form: {
        id: row.id,
        home_team: row.homeTeam,
        away_team: row.awayTeam,
        kickoff: kickoffForStatus(row.status),
        status: row.status,
      },
    });
  }
});

When('I enter a home score of {int} and away score of {int} for the first match', async ({ page }, home: number, away: number) => {
  const card = page.locator('[data-testid="match-card"]').first();
  await card.locator('input[name="home_goals"]').fill(String(home));
  await card.locator('input[name="away_goals"]').fill(String(away));
});

When('I click {string}', async ({ page }, label: string) => {
  await page.getByRole('button', { name: label, exact: true }).first().click();
});

Then('I should see my prediction of {string} on the first match card', async ({ page }, score: string) => {
  const card = page.locator('[data-testid="match-card"]').first();
  await expect(card.locator('[data-testid="matchup"]').getByText(score)).toBeVisible();
});

When('I submit a prediction via API for match {string}', async ({ page }, matchId: string) => {
  const resp = await page.request.post('/api/predictions', {
    form: { match_id: matchId, home_goals: '2', away_goals: '1' },
  });
  lastPredictionStatus = resp.status();
});

Then('I should see a {string} button', async ({ page }, label: string) => {
  await expect(page.getByRole('button', { name: label, exact: true }).first()).toBeVisible();
});

Then('I should not see a {string} button', async ({ page }, label: string) => {
  await expect(page.getByRole('button', { name: label, exact: true })).not.toBeVisible();
});

Then('I should see a disabled {string} button', async ({ page }, label: string) => {
  const btn = page.getByRole('button', { name: label, exact: true }).first();
  await expect(btn).toBeVisible();
  await expect(btn).toBeDisabled();
});

Then('the server should reject it with 403', async () => {
  expect(lastPredictionStatus).toBe(403);
});
