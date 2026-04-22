import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Given, Then } = createBdd();

// Seeds a prediction for the currently logged-in user by reading their session cookie.
Given('I have a seeded prediction of {int}-{int} for match {string}', async ({ context, request }, home: number, away: number, matchID: string) => {
  const cookies = await context.cookies();
  const session = cookies.find(c => c.name === '__session');
  if (!session) throw new Error('No session cookie — user must be logged in first');
  const { userID, handle } = JSON.parse(Buffer.from(session.value, 'base64').toString());
  await request.post('/admin/seed-prediction', {
    form: { match_id: matchID, user_id: userID, handle, home_goals: String(home), away_goals: String(away) },
  });
});

Then('I should not see {string} on the leaderboard', async ({ page }, handle: string) => {
  await expect(page.locator('[data-testid="leaderboard-row"]').filter({ hasText: handle })).toHaveCount(0);
});
