import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Given, Then } = createBdd();

Given('I am viewing on an iPhone 15', async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
});

Given('I am viewing on a Galaxy S24', async ({ page }) => {
  await page.setViewportSize({ width: 360, height: 780 });
});

Then('the page should not overflow horizontally', async ({ page }) => {
  const overflows = await page.evaluate(
    () => document.documentElement.scrollWidth > window.innerWidth
  );
  expect(overflows).toBe(false);
});

Then('each match card should be at most 180px tall', async ({ page }) => {
  await expect(page.getByTestId('match-card').first()).toBeVisible();
  const cards = page.getByTestId('match-card');
  const count = await cards.count();
  expect(count).toBeGreaterThan(0);
  for (let i = 0; i < count; i++) {
    const box = await cards.nth(i).boundingBox();
    expect(box!.height).toBeLessThanOrEqual(160);
  }
});

Then('the site header should be at most 64px tall', async ({ page }) => {
  const header = page.locator('.site-header');
  await expect(header).toBeVisible();
  const box = await header.boundingBox();
  expect(box!.height).toBeLessThanOrEqual(64);
});

Then('the Predict button should be at least 44px tall', async ({ page }) => {
  const btn = page.getByRole('button', { name: 'Predict' }).first();
  await expect(btn).toBeVisible();
  const box = await btn.boundingBox();
  expect(box!.height).toBeGreaterThanOrEqual(44);
});
