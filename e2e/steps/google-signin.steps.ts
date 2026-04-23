import { createBdd } from 'playwright-bdd';

const { When } = createBdd();

When('I sign in with Google as {string}', async ({ page }, email: string) => {
  await page.click('button[data-testid="google-signin"]');
  // signInWithRedirect navigates the main page to the emulator OAuth handler
  await page.waitForURL(/9099.*handler/, { timeout: 10000 });

  // Same emulator UI as popup but now full-page
  await page.getByText('Add new account').click();
  await page.fill('input[id="email-input"]', email, { force: true });
  await page.locator('button#sign-in').click({ force: true });

  // Wait for redirect back to the app and session to complete
  await page.waitForURL(/localhost/, { timeout: 10000 });
});
