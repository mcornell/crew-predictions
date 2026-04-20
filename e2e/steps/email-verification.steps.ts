import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Given, Then } = createBdd();

Given('I am logged in as a verified user', async ({ context }) => {
  await context.addCookies([{
    name: 'session',
    value: Buffer.from(JSON.stringify({ userID: 'google:verified', handle: 'VerifiedFan', emailVerified: true })).toString('base64'),
    domain: 'localhost',
    path: '/',
  }]);
});

Then('I should see an email verification banner', async ({ page }) => {
  await expect(page.getByTestId('email-verification-banner')).toBeVisible();
});

Then('I should not see an email verification banner', async ({ page }) => {
  await expect(page.getByTestId('email-verification-banner')).not.toBeVisible();
});
