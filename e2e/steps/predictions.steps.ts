import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { When, Then } = createBdd();

When('I enter a home score of {int} and away score of {int} for the first match', async ({ page }, home: number, away: number) => {
  const card = page.locator('[data-testid="match-card"]').first();
  await card.locator('input[name="home_goals"]').fill(String(home));
  await card.locator('input[name="away_goals"]').fill(String(away));
});

When('I click {string}', async ({ page }, label: string) => {
  await page.getByRole('button', { name: label }).first().click();
});

Then('I should see my prediction of {string} on the first match card', async ({ page }, score: string) => {
  const card = page.locator('[data-testid="match-card"]').first();
  await expect(card.getByText(score)).toBeVisible();
});
