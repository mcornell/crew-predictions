import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Given, When, Then } = createBdd();

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
  expect.soft(overflows).toBe(false);
});

Then('each match card should be at most 260px tall', async ({ page }) => {
  await expect.soft(page.getByTestId('match-card').first()).toBeVisible();
  const cards = page.getByTestId('match-card');
  const count = await cards.count();
  expect.soft(count).toBeGreaterThan(0);
  for (let i = 0; i < count; i++) {
    const box = await cards.nth(i).boundingBox();
    expect.soft(box!.height).toBeLessThanOrEqual(260);
  }
});

Then('the site header should be at most 64px tall', async ({ page }) => {
  const header = page.locator('.site-header');
  await expect.soft(header).toBeVisible();
  const box = await header.boundingBox();
  expect.soft(box!.height).toBeLessThanOrEqual(64);
});

Then('the Predict button should be at least 44px tall', async ({ page }) => {
  const btn = page.getByRole('button', { name: 'Predict' }).first();
  await expect.soft(btn).toBeVisible();
  const box = await btn.boundingBox();
  expect.soft(box!.height).toBeGreaterThanOrEqual(44);
});

Then('team names should not be clipped on any match card', async ({ page }) => {
  await expect.soft(page.getByTestId('match-card').first()).toBeVisible()
  const clipped = await page.evaluate(() => {
    const names = document.querySelectorAll('.team-name')
    return Array.from(names).some(el => el.scrollWidth > el.clientWidth)
  })
  expect.soft(clipped).toBe(false)
})

When('I tap the hamburger menu', async ({ page }) => {
  await page.getByTestId('hamburger').click();
});

Then('I should see the mobile navigation drawer', async ({ page }) => {
  await expect.soft(page.getByTestId('mobile-drawer')).toBeVisible();
});

When('I tap outside the drawer', async ({ page }) => {
  await page.locator('.drawer-backdrop').click();
});

Then('the mobile navigation drawer should be closed', async ({ page }) => {
  await expect.soft(page.getByTestId('mobile-drawer')).not.toBeVisible();
});

When('I tap the Leaderboard link in the drawer', async ({ page }) => {
  await page.getByTestId('drawer-lb-toggle').click();
  await page.getByTestId('mobile-drawer').getByRole('link', { name: 'Current Season' }).click();
});

Then('I should be on the leaderboard page', async ({ page }) => {
  await expect.soft(page).toHaveURL('/leaderboard');
});
