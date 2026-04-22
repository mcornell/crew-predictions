import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { When, Then } = createBdd();

Then('the display name field should contain {string}', async ({ page }, value: string) => {
  await expect(page.getByTestId('display-name-input')).toHaveValue(value);
});

Then('the location field should contain {string}', async ({ page }, value: string) => {
  await expect(page.getByTestId('location-input')).toHaveValue(value);
});

When('I visit my profile page', async ({ page, context }) => {
  const cookies = await context.cookies()
  const session = cookies.find(c => c.name === '__session')
  if (!session) throw new Error('No session cookie — user must be logged in first')
  const { userID } = JSON.parse(Buffer.from(session.value, 'base64').toString())
  await page.goto(`/profile/${userID}`)
  await page.waitForSelector('[data-testid="prediction-count"], form[data-testid="profile-form"]', { timeout: 5000 })
});

When('I set my display name to {string}', async ({ page }, name: string) => {
  await page.fill('input[data-testid="display-name-input"]', name);
});

When('I set my location to {string}', async ({ page }, loc: string) => {
  await page.fill('input[data-testid="location-input"]', loc);
});

When('I save my profile', async ({ page }) => {
  await page.click('button[type="submit"]');
  await page.waitForURL('/matches', { timeout: 10000 });
});

Then('I should see my prediction count as {int}', async ({ page }, count: number) => {
  await expect(page.getByTestId('prediction-count')).toHaveText(String(count));
});

Then('I should see my Aces Radio points', async ({ page }) => {
  await expect(page.getByTestId('aces-radio-points')).toBeVisible();
});

When('I click the handle {string} on the leaderboard', async ({ page }, handle: string) => {
  await page.locator('[data-testid="leaderboard-row"] a').filter({ hasText: handle }).first().click()
  await page.waitForURL(/\/profile\//, { timeout: 5000 })
});

Then('I should be on the profile page for that user', async ({ page }) => {
  await expect(page).toHaveURL(/\/profile\//)
});

Then('I should not see the profile edit form', async ({ page }) => {
  await expect(page.locator('form[data-testid="profile-form"]')).toHaveCount(0)
});
