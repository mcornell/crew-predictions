<template>
  <div class="page">
    <div class="lb-page-header">
      <h1 class="page-title">Leaderboard</h1>
      <div v-if="showSwitcher" class="lb-season-switcher" ref="switcherRef">
        <button
          class="lb-season-trigger"
          data-testid="season-selector"
          @click="seasonOpen = !seasonOpen"
        >{{ activeSeasonLabel }} <span class="lb-season-chevron">▾</span></button>
        <div v-if="seasonOpen" class="season-flyout" data-testid="season-flyout">
          <a
            href="/leaderboard"
            class="season-flyout-item"
            :class="{ 'season-flyout-item--current': !routeSeasonID }"
            @click="seasonOpen = false"
          >Current Season</a>
          <a
            v-for="s in pastSeasons"
            :key="s.id"
            :href="`/leaderboard/${s.id}`"
            class="season-flyout-item"
            :class="{ 'season-flyout-item--current': routeSeasonID === s.id }"
            @click="seasonOpen = false"
          >{{ s.name }}</a>
        </div>
      </div>
    </div>

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
import { ref, computed, onMounted, onUnmounted } from 'vue'
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

interface Season {
  id: string
  name: string
  isCurrent: boolean
}

const route = useRoute()

const entries = ref<Entry[]>([])
const activeSort = ref<'aces' | 'upper90' | 'grouchy'>('aces')
const loading = ref(true)
const error = ref<string | null>(null)

// All seasons known to the backend (current + past + future). The selector
// only surfaces past seasons; future seasons are filtered out so an admin
// pre-creating next year's season doesn't leak into the dropdown.
const allSeasons = ref<Season[]>([])
const seasonOpen = ref(false)
const switcherRef = ref<HTMLElement | null>(null)

const routeSeasonID = computed(() => route.params.season as string | undefined)

const pastSeasons = computed<Season[]>(() => {
  const currentIdx = allSeasons.value.findIndex(s => s.isCurrent)
  if (currentIdx < 0) return []
  return [...allSeasons.value.slice(0, currentIdx)].reverse()
})

const currentSeason = computed<Season | undefined>(() =>
  allSeasons.value.find(s => s.isCurrent)
)

// Show the switcher if there are past seasons to navigate to, OR if the user
// has landed directly on a historical leaderboard URL — in that case they
// need a way back to current even when no other past seasons exist.
const showSwitcher = computed(() => pastSeasons.value.length > 0 || !!routeSeasonID.value)

const activeSeasonLabel = computed(() => {
  if (routeSeasonID.value) {
    const found = allSeasons.value.find(s => s.id === routeSeasonID.value)
    return found?.name ?? routeSeasonID.value
  }
  return currentSeason.value?.name ?? 'Current Season'
})

function onDocClick(e: MouseEvent) {
  if (switcherRef.value && !switcherRef.value.contains(e.target as Node)) {
    seasonOpen.value = false
  }
}

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

  const [lbRes, seasonsRes] = await Promise.all([fetch(url), fetch('/api/seasons')])
  if (lbRes.ok) {
    const data = await lbRes.json()
    entries.value = data.entries ?? []
  } else {
    error.value = 'Could not load leaderboard. Try again later.'
  }
  if (seasonsRes.ok) {
    const sData = await seasonsRes.json()
    allSeasons.value = sData.seasons ?? []
  }
  loading.value = false
  document.addEventListener('click', onDocClick)
})

onUnmounted(() => {
  document.removeEventListener('click', onDocClick)
})
</script>
