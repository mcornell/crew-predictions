import { createRouter, createWebHistory } from 'vue-router'
import MatchesView from '../views/MatchesView.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: MatchesView },
    { path: '/matches', component: MatchesView },
  ],
})
