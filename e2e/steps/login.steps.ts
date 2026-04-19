import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Given, When, Then } = createBdd();

const AUTH_EMULATOR = process.env.FIREBASE_AUTH_EMULATOR_HOST || 'localhost:9099';
const PROJECT_ID = process.env.GOOGLE_CLOUD_PROJECT || 'crew-predictions';

Given('a test user exists with email {string} and password {string}', async ({ request }, email: string, password: string) => {
  // Create user via Auth emulator REST API (no-op if already exists)
  await request.post(
    `http://${AUTH_EMULATOR}/identitytoolkit.googleapis.com/v1/accounts:signUp?key=fake-key`,
    {
      data: { email, password, returnSecureToken: false },
      headers: { 'Content-Type': 'application/json' },
    }
  );
});

When('I visit the login page', async ({ page }) => {
  await page.goto('/login');
  // Wait for FirebaseUI provider buttons to render
  await page.waitForSelector('.firebaseui-idp-button', { timeout: 10000 });
});

When('I sign in with email {string} and password {string}', async ({ page }, email: string, password: string) => {
  // Click "Sign in with email" provider button
  await page.click('.firebaseui-idp-password');
  await page.waitForSelector('.firebaseui-id-email', { timeout: 5000 });
  await page.fill('.firebaseui-id-email', email);
  await page.click('.firebaseui-id-submit');
  await page.waitForSelector('.firebaseui-id-password', { timeout: 5000 });
  await page.fill('.firebaseui-id-password', password);
  await page.click('.firebaseui-id-submit');
});

Then('I should be on the matches page', async ({ page }) => {
  await expect(page).toHaveURL('/matches', { timeout: 10000 });
});
