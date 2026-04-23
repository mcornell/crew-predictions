<template>
  <div class="page">
    <RouterLink to="/matches" class="back-link" data-testid="back-link">← All Matches</RouterLink>

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

    <div class="sort-controls" v-if="scoringFormats.length > 0">
      <button
        v-for="fmt in scoringFormats"
        :key="fmt.key"
        class="sort-btn"
        :class="{ 'sort-btn--active': activeFormat === fmt.key }"
        @click="activeFormat = fmt.key"
      >
        {{ fmt.label }}
      </button>
    </div>

    <div v-if="sortedPredictions.length > 0" class="predictions-table">
      <div
        v-for="(entry, i) in sortedPredictions"
        :key="entry.userID"
        class="prediction-row"
        data-testid="prediction-row"
      >
        <span class="prediction-rank">{{ rankFor(i, sortedPredictions) }}</span>
        <span class="prediction-handle">{{ entry.handle }}</span>
        <span class="prediction-score">{{ entry.homeGoals }} – {{ entry.awayGoals }}</span>
        <span class="prediction-points">{{ pointsFor(entry) }}</span>
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

interface ScoringFormat {
  key: string
  label: string
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
const scoringFormats = ref<ScoringFormat[]>([])
const activeFormat = ref('acesRadio')
const loaded = ref(false)

const sortedPredictions = computed(() => {
  const key = activeFormat.value === 'acesRadio' ? 'acesRadioPoints' : 'upper90ClubPoints'
  return [...predictions.value].sort((a, b) => b[key] - a[key])
})

function pointsFor(entry: PredictionEntry): number {
  return activeFormat.value === 'acesRadio' ? entry.acesRadioPoints : entry.upper90ClubPoints
}

function rankFor(i: number, sorted: PredictionEntry[]): number {
  if (i === 0) return 1
  const key = activeFormat.value === 'acesRadio' ? 'acesRadioPoints' : 'upper90ClubPoints'
  if (sorted[i][key] === sorted[i - 1][key]) return rankFor(i - 1, sorted)
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
    scoringFormats.value = data.scoringFormats ?? []
    if (data.scoringFormats?.length > 0) {
      activeFormat.value = data.scoringFormats[0].key
    }
    if (match.value) {
      document.title = `${match.value.homeTeam} vs ${match.value.awayTeam} — Crew Predictions`
    }
  }
  loaded.value = true
})
</script>
