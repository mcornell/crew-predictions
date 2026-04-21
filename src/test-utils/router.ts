import { createRouter, createMemoryHistory } from 'vue-router'

const stub = { template: '<div />' }

export function makeRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: stub },
      { path: '/matches', component: stub },
      { path: '/login', component: stub },
      { path: '/signup', component: stub },
      { path: '/reset', component: stub },
      { path: '/leaderboard', component: stub },
      { path: '/profile', component: stub },
      { path: '/rules', component: stub },
      { path: '/:pathMatch(.*)*', component: stub },
    ],
  })
}
