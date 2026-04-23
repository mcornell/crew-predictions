import { createBdd } from 'playwright-bdd';
import { expect } from '@playwright/test';

const { Given, When, Then } = createBdd();

Given('I am not logged in', async ({ page }) => {
  await page.context().clearCookies();
});

When('I visit the matches page', async ({ page }) => {
  await page.goto('/matches');
});

When('I revisit the matches page', async ({ page }) => {
  await page.goto('/matches');
});

Then('I should see the {string} heading', async ({ page }, heading: string) => {
  await expect(page.getByRole('heading', { name: heading })).toBeVisible();
});

Then('I should see at least one Columbus Crew match card', async ({ page }) => {
  await expect(page.getByText(/Columbus Crew/i).first()).toBeVisible();
  await expect(page.getByTestId('match-card').first()).toBeVisible();
});

When('the admin triggers a match refresh', async ({ request }) => {
  const resp = await request.post('/admin/refresh-matches');
  expect(resp.status()).toBe(204);
});

When('the admin triggers a score poll', async ({ request }) => {
  const resp = await request.post('/admin/poll-scores');
  expect(resp.status()).toBe(204);
});

Then('the matches API includes match {string}', async ({ request }, matchId: string) => {
  const resp = await request.get('/api/matches');
  const body = await resp.json();
  const ids = (body.matches ?? []).map((m: any) => m.id);
  expect(ids).toContain(matchId);
});

Then('I should see a LIVE indicator on the match card', async ({ page }) => {
  await expect(page.locator('[data-testid="live-indicator"]').first()).toBeVisible();
});

Then('I should see a countdown on the match card', async ({ page }) => {
  await expect(page.locator('[data-testid="match-countdown"]').first()).toBeVisible();
});

Then('I should see a DELAYED indicator on the match card', async ({ page }) => {
  await expect(page.locator('[data-testid="delayed-indicator"]').first()).toBeVisible();
});

Given('the following matches are seeded in order:', async ({ request }, table: any) => {
  for (const row of table.hashes()) {
    const offset = parseInt(row.kickoffOffset ?? '24', 10);
    const d = new Date();
    d.setHours(d.getHours() + offset);
    await request.post('/admin/seed-match', {
      form: {
        id: row.id,
        home_team: row.homeTeam,
        away_team: row.awayTeam,
        kickoff: d.toISOString(),
        status: row.status,
        state: row.state ?? '',
        home_score: row.homeScore ?? '',
        away_score: row.awayScore ?? '',
      },
    });
  }
});

Then('match {string} should appear before match {string}', async ({ page }, firstId: string, secondId: string) => {
  const cards = page.locator('[data-testid="match-card"]');
  const cardIds = await cards.evaluateAll((els) =>
    els.map((el) => el.getAttribute('data-match-id'))
  );
  const firstIdx = cardIds.indexOf(firstId);
  const secondIdx = cardIds.indexOf(secondId);
  expect(firstIdx).toBeGreaterThanOrEqual(0);
  expect(secondIdx).toBeGreaterThanOrEqual(0);
  expect(firstIdx).toBeLessThan(secondIdx);
});

Given('the match {string} has already kicked off', async ({ request }, matchId: string) => {
  const d = new Date();
  d.setMinutes(d.getMinutes() - 5);
  await request.post('/admin/seed-match', {
    form: {
      id: matchId,
      home_team: 'Columbus Crew',
      away_team: 'FC Dallas',
      kickoff: d.toISOString(),
      status: 'STATUS_SCHEDULED',
      state: 'pre',
      home_score: '',
      away_score: '',
    },
  });
});
