<template>
  <header class="site-header">
    <a class="brand" href="/">Crew Predictions</a>
    <nav class="site-nav">
      <div class="nav-lb" ref="navLbRef">
        <button
          class="btn-ghost nav-lb-trigger"
          data-testid="season-selector"
          @click="leaderboardOpen = !leaderboardOpen"
        >Leaderboard</button>
        <div v-if="leaderboardOpen" class="season-flyout" data-testid="season-flyout">
          <a href="/leaderboard" class="season-flyout-item season-flyout-item--current" @click="leaderboardOpen = false">Current Season</a>
          <a
            v-for="s in seasons"
            :key="s.id"
            :href="`/leaderboard/${s.id}`"
            class="season-flyout-item"
            @click="leaderboardOpen = false"
          >{{ s.name }}</a>
        </div>
      </div>
      <a class="btn-ghost" href="/rules">Rules</a>
      <template v-if="user">
        <a class="btn-ghost" :href="`/profile/${user.userID}`">{{ user.handle }}</a>
        <a class="btn-ghost" href="/auth/logout">Sign out</a>
      </template>
      <a v-else-if="!loading" class="btn-primary" href="/login">Sign In</a>
    </nav>
    <span v-if="user && isMobile" class="mobile-user-handle">{{ user.handle }}</span>
    <button class="hamburger" data-testid="hamburger" @click="drawerOpen = true" aria-label="Open menu">
      <span></span><span></span><span></span>
    </button>
  </header>

  <template v-if="drawerOpen">
    <div class="drawer-backdrop" @click="drawerOpen = false"></div>
    <nav class="mobile-drawer" data-testid="mobile-drawer">
      <button class="drawer-lb-toggle" data-testid="drawer-lb-toggle" @click="drawerLbOpen = !drawerLbOpen">
        Leaderboard <span class="drawer-chevron">{{ drawerLbOpen ? '↑' : '↓' }}</span>
      </button>
      <template v-if="drawerLbOpen">
        <a href="/leaderboard" class="drawer-lb-item" @click="drawerOpen = false; drawerLbOpen = false">Current Season</a>
        <a
          v-for="s in seasons"
          :key="s.id"
          :href="`/leaderboard/${s.id}`"
          class="drawer-lb-item"
          @click="drawerOpen = false; drawerLbOpen = false"
        >{{ s.name }}</a>
      </template>
      <a href="/rules" @click="drawerOpen = false">Rules</a>
      <template v-if="user">
        <a :href="`/profile/${user.userID}`" @click="drawerOpen = false">{{ user.handle }}</a>
        <a href="/auth/logout" @click="drawerOpen = false">Sign out</a>
      </template>
      <a v-else-if="!loading" href="/login" @click="drawerOpen = false">Sign In</a>
    </nav>
  </template>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

interface Season {
  id: string
  name: string
  isCurrent: boolean
}

defineProps<{
  user: { userID: string; handle: string } | null
  loading?: boolean
  seasons?: Season[]
}>()

const leaderboardOpen = ref(false)
const drawerOpen = ref(false)
const drawerLbOpen = ref(false)
const isMobile = ref(false)
const navLbRef = ref<HTMLElement | null>(null)

function checkMobile() {
  isMobile.value = window.innerWidth <= 480
}

function onDocClick(e: MouseEvent) {
  if (navLbRef.value && !navLbRef.value.contains(e.target as Node)) {
    leaderboardOpen.value = false
  }
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  document.addEventListener('click', onDocClick)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
  document.removeEventListener('click', onDocClick)
})
</script>
