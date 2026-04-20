import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { When, Then } = createBdd();

When('I visit the sign-up page', async ({ page }) => {
  await page.goto('/signup');
  await page.waitForSelector('form[data-testid="signup-form"]', { timeout: 5000 });
});

When('I sign up with email {string} and password {string}', async ({ page }, email: string, password: string) => {
  await page.fill('input[type="email"]', email);
  await page.fill('input[type="password"]', password);
  await page.click('button[type="submit"]');
});

Then('I should stay on the sign-up page', async ({ page }) => {
  await expect(page).toHaveURL('/signup');
});

Then('I should see the error {string}', async ({ page }, message: string) => {
  await expect(page.locator('.form-error')).toHaveText(message);
});
