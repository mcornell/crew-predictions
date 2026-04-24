import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Then } = createBdd();

Then('the predictions for match {string} should include TwoOneBot with {int}-{int}', async ({ request }, matchID: string, home: number, away: number) => {
  const resp = await request.get(`/api/matches/${matchID}`);
  expect(resp.status()).toBe(200);
  const body = await resp.json();
  const predictions: any[] = body.predictions ?? [];
  const bot = predictions.find((p: any) => p.handle === "Upper 90 Club's TwoOneBot");
  expect(bot, 'TwoOneBot prediction not found').toBeTruthy();
  expect(bot.homeGoals).toBe(home);
  expect(bot.awayGoals).toBe(away);
});

Then('the predictions for match {string} should not include TwoOneBot', async ({ request }, matchID: string) => {
  const resp = await request.get(`/api/matches/${matchID}`);
  expect(resp.status()).toBe(200);
  const body = await resp.json();
  const predictions: any[] = body.predictions ?? [];
  const bot = predictions.find((p: any) => p.handle === "Upper 90 Club's TwoOneBot");
  expect(bot).toBeUndefined();
});

Then('I should see {string} on the leaderboard', async ({ page }, text: string) => {
  await expect(page.getByText(text).first()).toBeVisible();
});
