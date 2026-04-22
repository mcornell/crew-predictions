import { createBdd } from 'playwright-bdd';

const { Before } = createBdd();

Before({ tags: 'not @no-reset' }, async ({ request }) => {
  await request.delete('/admin/reset');
});
