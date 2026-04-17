import { test, expect } from '@playwright/test';

test('unauthenticated user sees upcoming Columbus Crew matches', async ({ page }) => {
  await page.goto('/matches');

  await expect(page.getByRole('heading', { name: /upcoming matches/i })).toBeVisible();
  await expect(page.getByText(/Columbus Crew/i).first()).toBeVisible();
  await expect(page.getByTestId('match-card').first()).toBeVisible();
});
