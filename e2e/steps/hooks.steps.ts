import { createBdd } from 'playwright-bdd';

const { Before } = createBdd();

Before(async ({ request }) => {
  await request.delete('/admin/reset');
});
