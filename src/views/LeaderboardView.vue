<template>
  <div class="page">
    <h1 class="page-title">Leaderboard</h1>

    <div class="leaderboard-section">
      <h2 class="section-title">Aces Radio</h2>
      <div
        v-for="(entry, i) in leaderboard.acesRadio"
        :key="entry.handle"
        class="leaderboard-row"
        data-testid="leaderboard-row"
      >
        <span class="leaderboard-rank">{{ i + 1 }}</span>
        <span class="leaderboard-handle">{{ entry.handle }}</span>
        <span class="leaderboard-points" data-testid="leaderboard-points">{{ entry.points }}</span>
      </div>
      <p v-if="leaderboard.acesRadio.length === 0" class="empty">No predictions scored yet.</p>
    </div>

    <div class="leaderboard-section">
      <h2 class="section-title">Upper 90 Club</h2>
      <div
        v-for="(entry, i) in leaderboard.upper90Club"
        :key="entry.handle"
        class="leaderboard-row"
        data-testid="leaderboard-row"
      >
        <span class="leaderboard-rank">{{ i + 1 }}</span>
        <span class="leaderboard-handle">{{ entry.handle }}</span>
        <span class="leaderboard-points" data-testid="leaderboard-points">{{ entry.points }}</span>
      </div>
      <p v-if="leaderboard.upper90Club.length === 0" class="empty">No predictions scored yet.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue'

interface Entry {
  handle: string
  points: number
}

const leaderboard = reactive<{ acesRadio: Entry[]; upper90Club: Entry[] }>({
  acesRadio: [],
  upper90Club: [],
})

onMounted(async () => {
  const res = await fetch('/api/leaderboard')
  if (res.ok) {
    const data = await res.json()
    leaderboard.acesRadio = data.acesRadio
    leaderboard.upper90Club = data.upper90Club
  }
})
</script>
