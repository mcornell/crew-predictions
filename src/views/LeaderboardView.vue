<template>
  <h1>Leaderboard</h1>
  <section>
    <h2>Aces Radio</h2>
    <div
      v-for="entry in leaderboard.acesRadio"
      :key="entry.handle"
      data-testid="leaderboard-row"
    >
      <span>{{ entry.handle }}</span>
      <span data-testid="leaderboard-points">{{ entry.points }}</span>
    </div>
  </section>
  <section>
    <h2>Upper 90 Club</h2>
    <div
      v-for="entry in leaderboard.upper90Club"
      :key="entry.handle"
      data-testid="leaderboard-row"
    >
      <span>{{ entry.handle }}</span>
      <span data-testid="leaderboard-points">{{ entry.points }}</span>
    </div>
  </section>
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
