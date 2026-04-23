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
