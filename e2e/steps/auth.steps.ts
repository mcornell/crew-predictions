import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Given, Then } = createBdd();

Given('I am logged in as {string}', async ({ page, context }, handle: string) => {
  // Seed a session cookie directly — bypasses OAuth for testing
  await context.addCookies([{
    name: 'session',
    value: Buffer.from(JSON.stringify({ handle })).toString('base64'),
    domain: 'localhost',
    path: '/',
  }]);
});

Then('I should see a {string} link in the header', async ({ page }, text: string) => {
  await expect(page.getByRole('banner').getByRole('link', { name: text })).toBeVisible();
});

Then('the {string} link should point to {string}', async ({ page }, text: string, href: string) => {
  const link = page.getByRole('banner').getByRole('link', { name: text });
  await expect(link).toHaveAttribute('href', href);
});

Then('I should see {string} in the header', async ({ page }, text: string) => {
  await expect(page.getByRole('banner').getByText(text)).toBeVisible();
});

Then('I should not see a {string} link', async ({ page }, text: string) => {
  await expect(page.getByRole('banner').getByRole('link', { name: text })).not.toBeVisible();
});
