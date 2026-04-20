import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Given, When, Then } = createBdd();

const AUTH_EMULATOR = process.env.FIREBASE_AUTH_EMULATOR_HOST || 'localhost:9099';

Given('a test user exists with email {string} and password {string}', async ({ request }, email: string, password: string) => {
  const resp = await request.post(
    `http://${AUTH_EMULATOR}/identitytoolkit.googleapis.com/v1/accounts:signUp?key=fake-key`,
    {
      data: { email, password, returnSecureToken: true },
      headers: { 'Content-Type': 'application/json' },
    }
  );
  if (!resp.ok()) {
    throw new Error(`Failed to create test user: HTTP ${resp.status()} — ${await resp.text()}`);
  }
});

When('I visit the login page', async ({ page }) => {
  await page.goto('/login');
  await page.waitForSelector('form[data-testid="login-form"]', { timeout: 5000 });
});

When('I sign in with email {string} and password {string}', async ({ page }, email: string, password: string) => {
  const consoleMessages: string[] = [];
  const firebaseRequests: string[] = [];
  page.on('console', msg => consoleMessages.push(`[${msg.type()}] ${msg.text()}`));
  page.on('request', req => {
    const url = req.url();
    if (url.includes('identitytoolkit') || url.includes('9099')) {
      firebaseRequests.push(`${req.method()} ${url}`);
    }
  });
  (page as any).__consoleMessages = consoleMessages;
  (page as any).__firebaseRequests = firebaseRequests;

  await page.fill('input[type="email"]', email);
  await page.fill('input[type="password"]', password);
  await page.click('button[type="submit"]');
});

Then('I should be on the matches page', async ({ page }) => {
  const consoleMessages: string[] = (page as any).__consoleMessages ?? [];
  const firebaseRequests: string[] = (page as any).__firebaseRequests ?? [];
  await expect(
    page,
    `firebase requests: ${JSON.stringify(firebaseRequests)} | console: ${JSON.stringify(consoleMessages)}`
  ).toHaveURL('/matches', { timeout: 10000 });
});
