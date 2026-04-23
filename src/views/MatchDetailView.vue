<template>
  <div class="page">
    <RouterLink to="/matches" class="back-link" data-testid="back-link">← All Matches</RouterLink>

    <p v-if="loading" data-testid="loading" class="status-msg">Loading…</p>
    <p v-else-if="error" data-testid="error" class="status-msg status-msg--error">{{ error }}</p>

    <div v-if="match" class="match-detail-header">
      <div class="matchup matchup--input" data-testid="match-score">
        <span class="team-name team-home">{{ match.homeTeam }}</span>
        <span class="inline-score">{{ match.homeScore }}</span>
        <span class="vs">vs</span>
        <span class="inline-score">{{ match.awayScore }}</span>
        <span class="team-name team-away">{{ match.awayTeam }}</span>
      </div>
      <div class="match-meta">{{ formatKickoff(match.kickoff) }}</div>
    </div>

    <div v-if="sortedPredictions.length > 0" class="lb-table lb-5col">
      <div class="lb-header">
        <span class="lb-cell lb-rank">RANK</span>
        <span class="lb-cell lb-handle">PREDICTOR</span>
        <span class="lb-cell lb-pick">PICK</span>
        <button
          class="lb-cell lb-pts lb-sort-btn"
          :class="{ 'lb-sort-btn--active': activeFormat === 'acesRadio' }"
          data-testid="sort-aces"
          @click="activeFormat = 'acesRadio'"
        >Aces Radio</button>
        <button
          class="lb-cell lb-pts lb-sort-btn"
          :class="{ 'lb-sort-btn--active': activeFormat === 'upper90Club' }"
          data-testid="sort-upper90"
          @click="activeFormat = 'upper90Club'"
        >Upper 90 Club</button>
      </div>

      <div
        v-for="(entry, i) in sortedPredictions"
        :key="entry.userID"
        class="lb-row"
        data-testid="prediction-row"
      >
        <span class="lb-cell lb-rank" data-testid="prediction-rank">{{ rankFor(i) }}</span>
        <span class="lb-cell lb-handle">{{ entry.handle }}</span>
        <span class="lb-cell lb-pick" data-testid="prediction-score">{{ entry.homeGoals }} – {{ entry.awayGoals }}</span>
        <span
          class="lb-cell lb-pts"
          :class="{ 'lb-pts--active': activeFormat === 'acesRadio' }"
          data-testid="prediction-aces-points"
          data-label="Aces Radio"
        >{{ entry.acesRadioPoints }}</span>
        <span
          class="lb-cell lb-pts"
          :class="{ 'lb-pts--active': activeFormat === 'upper90Club' }"
          data-testid="prediction-upper90-points"
          data-label="Upper 90 Club"
        >{{ entry.upper90ClubPoints }}</span>
      </div>
    </div>
    <p v-else-if="loaded" class="empty">No predictions were made for this match</p>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'

interface MatchInfo {
  id: string
  homeTeam: string
  awayTeam: string
  kickoff: string
  homeScore: string
  awayScore: string
}

interface PredictionEntry {
  userID: string
  handle: string
  homeGoals: number
  awayGoals: number
  acesRadioPoints: number
  upper90ClubPoints: number
}

const route = useRoute()
const match = ref<MatchInfo | null>(null)
const predictions = ref<PredictionEntry[]>([])
const activeFormat = ref<'acesRadio' | 'upper90Club'>('acesRadio')
const loaded = ref(false)
const loading = ref(true)
const error = ref<string | null>(null)

const sortedPredictions = computed(() => {
  const key = activeFormat.value === 'acesRadio' ? 'acesRadioPoints' : 'upper90ClubPoints'
  return [...predictions.value].sort((a, b) => b[key] - a[key])
})

function rankFor(i: number): number {
  if (i === 0) return 1
  const key = activeFormat.value === 'acesRadio' ? 'acesRadioPoints' : 'upper90ClubPoints'
  if (sortedPredictions.value[i][key] === sortedPredictions.value[i - 1][key]) return rankFor(i - 1)
  return i + 1
}

function formatKickoff(iso: string): string {
  const d = new Date(iso)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric', hour: 'numeric', minute: '2-digit', timeZoneName: 'short' })
}

onMounted(async () => {
  const matchId = route.params.matchId as string
  const res = await fetch(`/api/matches/${matchId}`)
  if (res.ok) {
    const data = await res.json()
    match.value = data.match
    predictions.value = data.predictions ?? []
    if (data.scoringFormats?.length > 0) {
      activeFormat.value = data.scoringFormats[0].key
    }
    if (match.value) {
      document.title = `${match.value.homeTeam} vs ${match.value.awayTeam} — Crew Predictions`
    }
  } else {
    error.value = 'Could not load match. Try again later.'
  }
  loaded.value = true
  loading.value = false
})
</script>
