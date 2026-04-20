import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { When, Then } = createBdd();

When('I click the {string} link', async ({ page }, text: string) => {
  // Scoped to the auth cross-link footer to avoid colliding with the
  // header's "Sign In" link (case-insensitive role-name matching).
  await page.locator('.auth-alt').getByRole('link', { name: text }).click();
});

Then('I should be on the sign-up page', async ({ page }) => {
  await expect(page).toHaveURL('/signup');
});

Then('I should be on the login page', async ({ page }) => {
  await expect(page).toHaveURL('/login');
});
