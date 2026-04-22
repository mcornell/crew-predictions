import { createRouter, createWebHistory } from 'vue-router'
import MatchesView from '../views/MatchesView.vue'
import LoginView from '../views/LoginView.vue'
import SignupView from '../views/SignupView.vue'
import ResetView from '../views/ResetView.vue'
import ProfileView from '../views/ProfileView.vue'
import LeaderboardView from '../views/LeaderboardView.vue'
import NotFoundView from '../views/NotFoundView.vue'
import RulesView from '../views/RulesView.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: MatchesView },
    { path: '/matches', component: MatchesView },
    { path: '/login', component: LoginView },
    { path: '/signup', component: SignupView },
    { path: '/reset', component: ResetView },
    { path: '/profile', redirect: '/matches' },
    { path: '/profile/:userID', component: ProfileView },
    { path: '/leaderboard', component: LeaderboardView },
    { path: '/rules', component: RulesView },
    { path: '/:pathMatch(.*)*', component: NotFoundView },
  ],
})
