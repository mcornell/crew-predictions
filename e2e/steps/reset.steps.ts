import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { When, Then } = createBdd();

When('I visit the reset page', async ({ page }) => {
  await page.goto('/reset');
  await page.waitForSelector('form[data-testid="reset-form"]', { timeout: 5000 });
});

When('I enter {string} in the reset email field', async ({ page }, email: string) => {
  await page.fill('input[type="email"]', email);
});

When('I submit the reset form', async ({ page }) => {
  await page.click('button[type="submit"]');
});

Then('I should see a reset confirmation message', async ({ page }) => {
  await expect(page.getByTestId('reset-confirmation')).toBeVisible({ timeout: 5000 });
});

Then('I should see a {string} link pointing to {string}', async ({ page }, text: string, href: string) => {
  const link = page.getByRole('link', { name: text });
  await expect(link).toBeVisible();
  await expect(link).toHaveAttribute('href', href);
});
