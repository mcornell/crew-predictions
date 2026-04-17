import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Then } = createBdd();

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
