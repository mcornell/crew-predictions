import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Given, When, Then } = createBdd();

Given('I am not logged in', async ({ page }) => {
  await page.context().clearCookies();
});

When('I visit the matches page', async ({ page }) => {
  await page.goto('/matches');
});

Then('I should see the {string} heading', async ({ page }, heading: string) => {
  await expect(page.getByRole('heading', { name: heading })).toBeVisible();
});

Then('I should see at least one Columbus Crew match card', async ({ page }) => {
  await expect(page.getByText(/Columbus Crew/i).first()).toBeVisible();
  await expect(page.getByTestId('match-card').first()).toBeVisible();
});
