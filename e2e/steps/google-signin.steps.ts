import { createBdd } from 'playwright-bdd';

const { When } = createBdd();

When('I sign in with Google as {string}', async ({ page }, email: string) => {
  const popupPromise = page.waitForEvent('popup');
  await page.click('button[data-testid="google-signin"]');
  const popup = await popupPromise;
  await popup.waitForLoadState();

  // Firebase Auth emulator OAuth handler: click "Add new account" to reveal
  // the email form, fill the email, and confirm. Material Design Components
  // styling hides the real <input>/<button> so force-click/fill.
  await popup.getByText('Add new account').click();
  await popup.fill('input[id="email-input"]', email, { force: true });
  await popup.locator('button#sign-in').click({ force: true });
});
