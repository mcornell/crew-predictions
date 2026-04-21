import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { When, Then } = createBdd();

When('I visit an unknown page', async ({ page }) => {
  await page.goto('/this-page-does-not-exist');
});

Then('I should see a not-found message', async ({ page }) => {
  await expect(page.getByTestId('not-found')).toBeVisible();
});

Then('I should see a link home', async ({ page }) => {
  await expect(page.getByRole('link', { name: /home|matches/i })).toBeVisible();
});

Then('the page title should be {string}', async ({ page }, title: string) => {
  await expect(page).toHaveTitle(title);
});

Then('I should see a site header with {string}', async ({ page }, text: string) => {
  await expect(page.getByRole('banner').getByText(text)).toBeVisible();
});

Then('the page should load HTMX', async ({ page }) => {
  const htmxLoaded = await page.evaluate(() => typeof (window as any).htmx !== 'undefined');
  expect(htmxLoaded).toBe(true);
});
