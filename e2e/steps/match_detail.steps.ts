import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { When, Then } = createBdd();

When('I click on the result card for match {string}', async ({ page }, matchID: string) => {
  await page.locator(`[data-testid="result-card"][data-match-id="${matchID}"]`).click();
});

Then('I should be on the match detail page for {string}', async ({ page }, matchID: string) => {
  await expect.soft(page).toHaveURL(new RegExp(`/matches/${matchID}`));
});

Then('I should see the match header with {string} vs {string}', async ({ page }, homeTeam: string, awayTeam: string) => {
  await expect.soft(page.getByText(homeTeam).first()).toBeVisible();
  await expect.soft(page.getByText(awayTeam).first()).toBeVisible();
});

Then('I should see {string} in the predictions table', async ({ page }, handle: string) => {
  await expect.soft(page.locator('[data-testid="prediction-row"]').filter({ hasText: handle })).toBeVisible();
});

Then('{string} should have more points than {string}', async ({ page }, handle1: string, handle2: string) => {
  const rows = page.locator('[data-testid="prediction-row"]');
  const handle1Index = await rows.evaluateAll((els, h1) =>
    els.findIndex(el => el.textContent?.includes(h1)), handle1);
  const handle2Index = await rows.evaluateAll((els, h2) =>
    els.findIndex(el => el.textContent?.includes(h2)), handle2);
  expect.soft(handle1Index).toBeGreaterThanOrEqual(0);
  expect.soft(handle2Index).toBeGreaterThanOrEqual(0);
  expect.soft(handle1Index).toBeLessThan(handle2Index);
});

Then('the result card for match {string} should link to {string}', async ({ page }, matchID: string, expectedPath: string) => {
  const card = page.locator(`[data-testid="result-card"][data-match-id="${matchID}"]`);
  await expect.soft(card).toBeVisible();
  const href = await card.getAttribute('href');
  expect.soft(href).toContain(expectedPath);
});

Then('the upcoming card for match {string} should not have a detail link', async ({ page }, matchID: string) => {
  const card = page.locator(`[data-testid="match-card"][data-match-id="${matchID}"]`);
  await expect.soft(card).toBeVisible();
  const tagName = await card.evaluate(el => el.tagName.toLowerCase());
  expect.soft(tagName).not.toBe('a');
});

When('I visit the match detail page for {string}', async ({ page }, matchID: string) => {
  await page.goto(`/matches/${matchID}`);
});

Then('I should see {string}', async ({ page }, text: string) => {
  await expect.soft(page.getByText(text).first()).toBeVisible();
});

Then('I should see the Grouchy column header in the predictions table', async ({ page }) => {
  await expect.soft(page.locator('[data-testid="sort-grouchy"]')).toBeVisible();
});

Then('{string} should have {int} Grouchy point in the predictions table', async ({ page }, handle: string, points: number) => {
  const row = page.locator('[data-testid="prediction-row"]').filter({ hasText: handle }).first();
  await expect.soft(row.locator('[data-testid="prediction-grouchy-points"]')).toHaveText(String(points));
});

Then('I should see the LIVE indicator on the match detail page', async ({ page }) => {
  await expect.soft(page.locator('[data-testid="live-indicator-detail"]')).toBeVisible();
});

Then('the match detail header should show score {string} to {string}', async ({ page }, home: string, away: string) => {
  const score = page.locator('[data-testid="match-score"]');
  await expect.soft(score.locator('.inline-score').nth(0)).toHaveText(home);
  await expect.soft(score.locator('.inline-score').nth(1)).toHaveText(away);
});

Then('the projected points label should be visible', async ({ page }) => {
  await expect.soft(page.locator('[data-testid="projected-label"]')).toBeVisible();
});

Then('{string} should have projected points greater than {string}', async ({ page }, handle1: string, handle2: string) => {
  const rows = page.locator('[data-testid="prediction-row"]');
  const idx1 = await rows.evaluateAll((els, h) => els.findIndex(el => el.textContent?.includes(h)), handle1);
  const idx2 = await rows.evaluateAll((els, h) => els.findIndex(el => el.textContent?.includes(h)), handle2);
  expect.soft(idx1).toBeGreaterThanOrEqual(0);
  expect.soft(idx2).toBeGreaterThanOrEqual(0);
  expect.soft(idx1).toBeLessThan(idx2);
});

Then('the now playing card for match {string} should link to {string}', async ({ page }, matchID: string, expectedPath: string) => {
  const card = page.locator(`[data-testid="now-playing-card"][data-match-id="${matchID}"]`);
  await expect.soft(card).toBeVisible();
  const href = await card.getAttribute('href');
  expect.soft(href).toContain(expectedPath);
});

Then('I should see the venue {string} on the match detail page', async ({ page }, venue: string) => {
  await expect.soft(page.locator('[data-testid="match-detail-venue"]')).toHaveText(venue);
});

Then('I should see an ESPN link for match {string}', async ({ page }, matchId: string) => {
  const link = page.locator('[data-testid="espn-link"]');
  await expect.soft(link).toBeVisible();
  const href = await link.getAttribute('href');
  expect.soft(href).toContain(`gameId/${matchId}`);
});

Then('I should see home record {string} on the match detail page', async ({ page }, record: string) => {
  await expect.soft(page.locator('[data-testid="home-record"]')).toHaveText(record);
});

Then('I should see home form {string} on the match detail page', async ({ page }, form: string) => {
  await expect.soft(page.locator('[data-testid="home-form"]')).toHaveText(form);
});

Then('I should see the attendance {string} on the match detail page', async ({ page }, attendance: string) => {
  await expect.soft(page.locator('[data-testid="match-detail-attendance"]')).toHaveText(attendance);
});

Then('I should see the event timeline on the match detail page', async ({ page }) => {
  await expect.soft(page.locator('[data-testid="match-events"]')).toBeVisible();
});

Then('I should see at least one event in the timeline', async ({ page }) => {
  await expect.soft(page.locator('[data-testid="match-event"]').first()).toBeVisible();
});

Then('I should see the home team logo on the match detail page', async ({ page }) => {
  await expect.soft(page.locator('[data-testid="home-logo"]')).toBeVisible();
});

Then('I should see the away team logo on the match detail page', async ({ page }) => {
  await expect.soft(page.locator('[data-testid="away-logo"]')).toBeVisible();
});

Then('I should see the referee on the match detail page', async ({ page }) => {
  const ref = page.locator('[data-testid="match-referee"]');
  await expect.soft(ref).toBeVisible();
  await expect.soft(ref).not.toHaveText('');
});
