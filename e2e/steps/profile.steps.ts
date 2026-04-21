import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { When, Then } = createBdd();

Then('the display name field should contain {string}', async ({ page }, value: string) => {
  await expect(page.getByTestId('display-name-input')).toHaveValue(value);
});

When('I visit the profile page', async ({ page }) => {
  await page.goto('/profile');
  await page.waitForSelector('form[data-testid="profile-form"]', { timeout: 5000 });
});

When('I set my display name to {string}', async ({ page }, name: string) => {
  await page.fill('input[data-testid="display-name-input"]', name);
});

When('I save my profile', async ({ page }) => {
  await page.click('button[type="submit"]');
  await page.waitForURL('/matches', { timeout: 10000 });
});
