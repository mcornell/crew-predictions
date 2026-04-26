import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { When, Then } = createBdd();

When('I click on the result card for match {string}', async ({ page }, matchID: string) => {
  await page.locator(`[data-testid="result-card"][data-match-id="${matchID}"]`).click();
});

Then('I should be on the match detail page for {string}', async ({ page }, matchID: string) => {
  await expect(page).toHaveURL(new RegExp(`/matches/${matchID}`));
});

Then('I should see the match header with {string} vs {string}', async ({ page }, homeTeam: string, awayTeam: string) => {
  await expect(page.getByText(homeTeam).first()).toBeVisible();
  await expect(page.getByText(awayTeam).first()).toBeVisible();
});

Then('I should see {string} in the predictions table', async ({ page }, handle: string) => {
  await expect(page.locator('[data-testid="prediction-row"]').filter({ hasText: handle })).toBeVisible();
});

Then('{string} should have more points than {string}', async ({ page }, handle1: string, handle2: string) => {
  const rows = page.locator('[data-testid="prediction-row"]');
  const handle1Index = await rows.evaluateAll((els, h1) =>
    els.findIndex(el => el.textContent?.includes(h1)), handle1);
  const handle2Index = await rows.evaluateAll((els, h2) =>
    els.findIndex(el => el.textContent?.includes(h2)), handle2);
  expect(handle1Index).toBeGreaterThanOrEqual(0);
  expect(handle2Index).toBeGreaterThanOrEqual(0);
  expect(handle1Index).toBeLessThan(handle2Index);
});

Then('the result card for match {string} should link to {string}', async ({ page }, matchID: string, expectedPath: string) => {
  const card = page.locator(`[data-testid="result-card"][data-match-id="${matchID}"]`);
  await expect(card).toBeVisible();
  const href = await card.getAttribute('href');
  expect(href).toContain(expectedPath);
});

Then('the upcoming card for match {string} should not have a detail link', async ({ page }, matchID: string) => {
  const card = page.locator(`[data-testid="match-card"][data-match-id="${matchID}"]`);
  await expect(card).toBeVisible();
  const tagName = await card.evaluate(el => el.tagName.toLowerCase());
  expect(tagName).not.toBe('a');
});

When('I visit the match detail page for {string}', async ({ page }, matchID: string) => {
  await page.goto(`/matches/${matchID}`);
});

Then('I should see {string}', async ({ page }, text: string) => {
  await expect(page.getByText(text).first()).toBeVisible();
});

Then('I should see the Grouchy column header in the predictions table', async ({ page }) => {
  await expect(page.locator('[data-testid="sort-grouchy"]')).toBeVisible();
});

Then('{string} should have {int} Grouchy point in the predictions table', async ({ page }, handle: string, points: number) => {
  const row = page.locator('[data-testid="prediction-row"]').filter({ hasText: handle }).first();
  await expect(row.locator('[data-testid="prediction-grouchy-points"]')).toHaveText(String(points));
});

Then('I should see the LIVE indicator on the match detail page', async ({ page }) => {
  await expect(page.locator('[data-testid="live-indicator-detail"]')).toBeVisible();
});

Then('the match detail header should show score {string} to {string}', async ({ page }, home: string, away: string) => {
  const score = page.locator('[data-testid="match-score"]');
  await expect(score.locator('.inline-score').nth(0)).toHaveText(home);
  await expect(score.locator('.inline-score').nth(1)).toHaveText(away);
});

Then('the projected points label should be visible', async ({ page }) => {
  await expect(page.locator('[data-testid="projected-label"]')).toBeVisible();
});

Then('{string} should have projected points greater than {string}', async ({ page }, handle1: string, handle2: string) => {
  const rows = page.locator('[data-testid="prediction-row"]');
  const idx1 = await rows.evaluateAll((els, h) => els.findIndex(el => el.textContent?.includes(h)), handle1);
  const idx2 = await rows.evaluateAll((els, h) => els.findIndex(el => el.textContent?.includes(h)), handle2);
  expect(idx1).toBeGreaterThanOrEqual(0);
  expect(idx2).toBeGreaterThanOrEqual(0);
  expect(idx1).toBeLessThan(idx2);
});

Then('the now playing card for match {string} should link to {string}', async ({ page }, matchID: string, expectedPath: string) => {
  const card = page.locator(`[data-testid="now-playing-card"][data-match-id="${matchID}"]`);
  await expect(card).toBeVisible();
  const href = await card.getAttribute('href');
  expect(href).toContain(expectedPath);
});
