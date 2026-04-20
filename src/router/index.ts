import { createRouter, createWebHistory } from 'vue-router'
import MatchesView from '../views/MatchesView.vue'
import LoginView from '../views/LoginView.vue'
import SignupView from '../views/SignupView.vue'
import LeaderboardView from '../views/LeaderboardView.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: MatchesView },
    { path: '/matches', component: MatchesView },
    { path: '/login', component: LoginView },
    { path: '/signup', component: SignupView },
    { path: '/leaderboard', component: LeaderboardView },
  ],
})
