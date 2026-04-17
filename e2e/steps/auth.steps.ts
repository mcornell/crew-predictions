import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Then } = createBdd();

Then('I should see a {string} link in the header', async ({ page }, text: string) => {
  await expect(page.getByRole('banner').getByRole('link', { name: text })).toBeVisible();
});

Then('the {string} link should point to {string}', async ({ page }, text: string, href: string) => {
  const link = page.getByRole('banner').getByRole('link', { name: text });
  await expect(link).toHaveAttribute('href', href);
});
