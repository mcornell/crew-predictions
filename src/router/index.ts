import { createRouter, createWebHistory } from 'vue-router'
import MatchesView from '../views/MatchesView.vue'
import LoginView from '../views/LoginView.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: MatchesView },
    { path: '/matches', component: MatchesView },
    { path: '/login', component: LoginView },
  ],
})
