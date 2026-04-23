import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Given, When, Then } = createBdd();

Given('{string} predicted {int}-{int} for match {string}', async ({ request }, handle: string, home: number, away: number, matchID: string) => {
  await request.post('/admin/seed-prediction', {
    form: { match_id: matchID, user_id: `google:${handle}`, handle, home_goals: String(home), away_goals: String(away) },
  });
});

Given('the final score for match {string} was {int}-{int} with Columbus away', async ({ request }, matchID: string, home: number, away: number) => {
  await request.post('/admin/results', {
    form: { match_id: matchID, home_team: 'Portland Timbers', away_team: 'Columbus Crew', home_goals: String(home), away_goals: String(away) },
  });
});

Given('the final score for match {string} was {int}-{int} with Columbus home', async ({ request }, matchID: string, home: number, away: number) => {
  await request.post('/admin/results', {
    form: { match_id: matchID, home_team: 'Columbus Crew', away_team: 'FC Dallas', home_goals: String(home), away_goals: String(away) },
  });
});

When('I visit the leaderboard', async ({ page }) => {
  await page.goto('/leaderboard');
});

Then('I should see {string} with {int} Aces Radio points', async ({ page }, handle: string, points: number) => {
  const row = page.locator('[data-testid="leaderboard-row"]').filter({ hasText: handle }).first();
  await expect(row.locator('[data-testid="leaderboard-points"]')).toHaveText(String(points));
});

Then('I should see {string} with {int} Upper90Club points', async ({ page }, handle: string, points: number) => {
  const row = page.locator('[data-testid="leaderboard-row"]').filter({ hasText: handle }).last();
  await expect(row.locator('[data-testid="leaderboard-points"]')).toHaveText(String(points));
});
