import { createBdd } from 'playwright-bdd';

const { When } = createBdd();

When('I visit the sign-up page', async ({ page }) => {
  await page.goto('/signup');
  await page.waitForSelector('form[data-testid="signup-form"]', { timeout: 5000 });
});

When('I sign up with email {string} and password {string}', async ({ page }, email: string, password: string) => {
  await page.fill('input[type="email"]', email);
  await page.fill('input[type="password"]', password);
  await page.click('button[type="submit"]');
});
