import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Given, When, Then } = createBdd();

Given('I am not logged in', async ({ page }) => {
  await page.context().clearCookies();
});

When('I visit the matches page', async ({ page }) => {
  await page.goto('/matches');
});

When('I revisit the matches page', async ({ page }) => {
  await page.goto('/matches');
});

Then('I should see the {string} heading', async ({ page }, heading: string) => {
  await expect(page.getByRole('heading', { name: heading })).toBeVisible();
});

Then('I should see at least one Columbus Crew match card', async ({ page }) => {
  await expect(page.getByText(/Columbus Crew/i).first()).toBeVisible();
  await expect(page.getByTestId('match-card').first()).toBeVisible();
});

When('the admin triggers a match refresh', async ({ request }) => {
  const resp = await request.post('/admin/refresh-matches');
  expect(resp.status()).toBe(204);
});

When('the admin triggers a score poll', async ({ request }) => {
  const resp = await request.post('/admin/poll-scores');
  expect(resp.status()).toBe(204);
});

Then('the matches API includes match {string}', async ({ request }, matchId: string) => {
  const resp = await request.get('/api/matches');
  const body = await resp.json();
  const ids = (body.matches ?? []).map((m: any) => m.id);
  expect(ids).toContain(matchId);
});

Then('I should see a LIVE indicator on the match card', async ({ page }) => {
  await expect(page.locator('[data-testid="live-indicator"]').first()).toBeVisible();
});

Then('I should see a countdown on the match card', async ({ page }) => {
  await expect(page.locator('[data-testid="match-countdown"]').first()).toBeVisible();
});
