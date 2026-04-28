import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Given, When, Then } = createBdd();

Given('season {string} has been archived with {string} at {int} Aces Radio points', async ({ request }, seasonID: string, handle: string, points: number) => {
  await request.post('/admin/seed-season', {
    form: { season_id: seasonID, entry_handle: handle, entry_aces: String(points) }
  })
});

When('I visit the historical leaderboard for season {string}', async ({ page }, season: string) => {
  await page.goto(`/leaderboard/${season}`)
});

Then('I should see a season selector on the leaderboard page', async ({ page }) => {
  await expect(page.getByTestId('season-selector')).toBeVisible()
});
