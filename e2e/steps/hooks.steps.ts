import { createBdd } from 'playwright-bdd';

const { Before } = createBdd();

Before({ tags: '@reset' }, async ({ request }) => {
  await request.delete('/admin/reset');
});
