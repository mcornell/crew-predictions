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
          <button
            class="lb-sort-btn"
            :class="{ 'lb-sort-btn--active': activeSort === 'grouchy' }"
            data-testid="mobile-sort-grouchy"
            @click="activeSort = 'grouchy'"
          >Grouchy™</button>
        </div>

        <div class="lb-table lb-5col">
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
          <button
            class="lb-cell lb-pts lb-sort-btn"
            :class="{ 'lb-sort-btn--active': activeSort === 'grouchy' }"
            data-testid="sort-grouchy"
            @click="activeSort = 'grouchy'"
          >GROUCHY™</button>
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
          >{{ upper90For(entry) }}</span>
          <span
            class="lb-cell lb-pts"
            :class="{ 'lb-pts--active': activeSort === 'grouchy' }"
            data-testid="leaderboard-grouchy-points"
            data-label="Grouchy™"
          >{{ entry.grouchyPoints }}</span>
        </div>
        </div>
      </template>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'

interface Entry {
  userID?: string
  handle: string
  acesRadioPoints: number
  upper90ClubPoints?: number
  upper90Points?: number
  grouchyPoints: number
  hasProfile?: boolean
}

const route = useRoute()

const entries = ref<Entry[]>([])
const activeSort = ref<'aces' | 'upper90' | 'grouchy'>('aces')
const loading = ref(true)
const error = ref<string | null>(null)

function upper90For(e: Entry): number {
  return e.upper90ClubPoints ?? e.upper90Points ?? 0
}

const sortedEntries = computed(() => {
  const key = activeSort.value
  return [...entries.value].sort((a, b) => {
    if (key === 'aces') return b.acesRadioPoints - a.acesRadioPoints
    if (key === 'upper90') return upper90For(b) - upper90For(a)
    return b.grouchyPoints - a.grouchyPoints
  })
})

function rankFor(i: number): number {
  if (i === 0) return 1
  const a = sortedEntries.value[i]
  const b = sortedEntries.value[i - 1]
  const key = activeSort.value
  const av = key === 'aces' ? a.acesRadioPoints : key === 'upper90' ? upper90For(a) : a.grouchyPoints
  const bv = key === 'aces' ? b.acesRadioPoints : key === 'upper90' ? upper90For(b) : b.grouchyPoints
  if (av === bv) return rankFor(i - 1)
  return i + 1
}

onMounted(async () => {
  document.title = 'Leaderboard — Crew Predictions'
  const seasonID = route.params.season as string | undefined
  const url = seasonID ? `/api/leaderboard/${seasonID}` : '/api/leaderboard'
  const res = await fetch(url)
  if (res.ok) {
    const data = await res.json()
    entries.value = data.entries ?? []
  } else {
    error.value = 'Could not load leaderboard. Try again later.'
  }
  loading.value = false
})
</script>
