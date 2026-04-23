<template>
  <div class="page">
    <h1 class="page-title">Leaderboard</h1>

    <p v-if="loading" data-testid="loading" class="status-msg">Loading…</p>
    <p v-else-if="error" data-testid="error" class="status-msg status-msg--error">{{ error }}</p>

    <template v-else>
      <div class="leaderboard-section">
        <h2 class="section-title">Aces Radio</h2>
        <div
          v-for="(entry, i) in leaderboard.acesRadio"
          :key="entry.userID"
          class="leaderboard-row"
          data-testid="leaderboard-row"
        >
          <span class="leaderboard-rank">{{ i + 1 }}</span>
          <RouterLink v-if="entry.hasProfile" :to="`/profile/${entry.userID}`" class="leaderboard-handle" data-testid="leaderboard-handle">{{ entry.handle }}</RouterLink>
          <span v-else class="leaderboard-handle" data-testid="leaderboard-handle">{{ entry.handle }}</span>
          <span class="leaderboard-points" data-testid="leaderboard-points">{{ entry.points }}</span>
        </div>
        <p v-if="leaderboard.acesRadio.length === 0" class="empty">No predictions scored yet.</p>
      </div>

      <div class="leaderboard-section">
        <h2 class="section-title">Upper 90 Club</h2>
        <div
          v-for="(entry, i) in leaderboard.upper90Club"
          :key="entry.userID"
          class="leaderboard-row"
          data-testid="leaderboard-row"
        >
          <span class="leaderboard-rank">{{ i + 1 }}</span>
          <RouterLink v-if="entry.hasProfile" :to="`/profile/${entry.userID}`" class="leaderboard-handle" data-testid="leaderboard-handle">{{ entry.handle }}</RouterLink>
          <span v-else class="leaderboard-handle" data-testid="leaderboard-handle">{{ entry.handle }}</span>
          <span class="leaderboard-points" data-testid="leaderboard-points">{{ entry.points }}</span>
        </div>
        <p v-if="leaderboard.upper90Club.length === 0" class="empty">No predictions scored yet.</p>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'

interface Entry {
  userID: string
  handle: string
  points: number
  hasProfile: boolean
}

const leaderboard = reactive<{ acesRadio: Entry[]; upper90Club: Entry[] }>({
  acesRadio: [],
  upper90Club: [],
})
const loading = ref(true)
const error = ref<string | null>(null)

onMounted(async () => {
  document.title = 'Leaderboard — Crew Predictions'
  const res = await fetch('/api/leaderboard')
  if (res.ok) {
    const data = await res.json()
    leaderboard.acesRadio = data.acesRadio
    leaderboard.upper90Club = data.upper90Club
  } else {
    error.value = 'Could not load leaderboard. Try again later.'
  }
  loading.value = false
})
</script>
