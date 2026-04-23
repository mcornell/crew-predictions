<template>
  <div class="page">
    <h1 class="page-title">Leaderboard</h1>

    <p v-if="loading" data-testid="loading" class="status-msg">Loading…</p>
    <p v-else-if="error" data-testid="error" class="status-msg status-msg--error">{{ error }}</p>

    <template v-else>
      <p v-if="sortedEntries.length === 0" class="empty">No predictions scored yet.</p>

      <template v-else>
        <div class="lb-mobile-sort">
          <button
            class="lb-sort-btn"
            :class="{ 'lb-sort-btn--active': activeSort === 'aces' }"
            data-testid="mobile-sort-aces"
            @click="activeSort = 'aces'"
          >Aces Radio</button>
          <button
            class="lb-sort-btn"
            :class="{ 'lb-sort-btn--active': activeSort === 'upper90' }"
            data-testid="mobile-sort-upper90"
            @click="activeSort = 'upper90'"
          >Upper 90 Club</button>
        </div>

        <div class="lb-table lb-4col">
        <div class="lb-header">
          <span class="lb-cell lb-rank">RANK</span>
          <span class="lb-cell lb-handle">PREDICTOR</span>
          <button
            class="lb-cell lb-pts lb-sort-btn"
            :class="{ 'lb-sort-btn--active': activeSort === 'aces' }"
            data-testid="sort-aces"
            @click="activeSort = 'aces'"
          >ACES RADIO</button>
          <button
            class="lb-cell lb-pts lb-sort-btn"
            :class="{ 'lb-sort-btn--active': activeSort === 'upper90' }"
            data-testid="sort-upper90"
            @click="activeSort = 'upper90'"
          >UPPER 90 CLUB</button>
        </div>

        <div
          v-for="(entry, i) in sortedEntries"
          :key="entry.userID"
          class="lb-row"
          data-testid="leaderboard-row"
        >
          <span class="lb-cell lb-rank" data-testid="leaderboard-rank">{{ rankFor(i) }}</span>
          <span class="lb-cell lb-handle">
            <RouterLink
              v-if="entry.hasProfile"
              :to="`/profile/${entry.userID}`"
              class="lb-handle-link"
              data-testid="leaderboard-handle"
            >{{ entry.handle }}</RouterLink>
            <span v-else data-testid="leaderboard-handle">{{ entry.handle }}</span>
          </span>
          <span
            class="lb-cell lb-pts"
            :class="{ 'lb-pts--active': activeSort === 'aces' }"
            data-testid="leaderboard-aces-points"
            data-label="Aces Radio"
          >{{ entry.acesRadioPoints }}</span>
          <span
            class="lb-cell lb-pts"
            :class="{ 'lb-pts--active': activeSort === 'upper90' }"
            data-testid="leaderboard-upper90-points"
            data-label="Upper 90 Club"
          >{{ entry.upper90ClubPoints }}</span>
        </div>
        </div>
      </template>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

interface Entry {
  userID: string
  handle: string
  acesRadioPoints: number
  upper90ClubPoints: number
  hasProfile: boolean
}

const entries = ref<Entry[]>([])
const activeSort = ref<'aces' | 'upper90'>('aces')
const loading = ref(true)
const error = ref<string | null>(null)

const sortedEntries = computed(() => {
  const key = activeSort.value === 'aces' ? 'acesRadioPoints' : 'upper90ClubPoints'
  return [...entries.value].sort((a, b) => b[key] - a[key])
})

function rankFor(i: number): number {
  if (i === 0) return 1
  const key = activeSort.value === 'aces' ? 'acesRadioPoints' : 'upper90ClubPoints'
  if (sortedEntries.value[i][key] === sortedEntries.value[i - 1][key]) return rankFor(i - 1)
  return i + 1
}

onMounted(async () => {
  document.title = 'Leaderboard — Crew Predictions'
  const res = await fetch('/api/leaderboard')
  if (res.ok) {
    const data = await res.json()
    entries.value = data.entries ?? []
  } else {
    error.value = 'Could not load leaderboard. Try again later.'
  }
  loading.value = false
})
</script>
