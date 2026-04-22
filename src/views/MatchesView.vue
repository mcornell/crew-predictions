<template>
  <div class="page">
    <section v-if="upcomingMatches.length > 0">
      <h1 class="page-title">Upcoming</h1>
      <div class="match-list">
        <div
          v-for="match in upcomingMatches"
          :key="match.id"
          class="match-card"
          :class="{ 'has-prediction': savedPredictions[match.id] }"
          data-testid="match-card"
        >
          <div class="match-info">
            <span v-if="match.state === 'in'" class="live-indicator" data-testid="live-indicator">● LIVE</span>
            <template v-if="savedPredictions[match.id]">
              <div class="matchup matchup--input" data-testid="matchup">
                <span class="team-name team-home">{{ match.homeTeam }}</span>
                <span class="inline-score">{{ savedPredictions[match.id]!.homeGoals }}</span>
                <span class="vs">vs</span>
                <span class="inline-score">{{ savedPredictions[match.id]!.awayGoals }}</span>
                <span class="team-name team-away">{{ match.awayTeam }}</span>
              </div>
              <div class="match-meta">{{ formatKickoff(match.kickoff) }} — <span class="saved-label">Locked in</span></div>
            </template>
            <template v-else>
              <div class="matchup matchup--input" data-testid="matchup">
                <span class="team-name team-home">{{ match.homeTeam }}</span>
                <input class="score-input" name="home_goals" type="number" min="0" max="99" v-model="inputs[match.id].home" placeholder="0" />
                <span class="vs">vs</span>
                <input class="score-input" name="away_goals" type="number" min="0" max="99" v-model="inputs[match.id].away" placeholder="0" />
                <span class="team-name team-away">{{ match.awayTeam }}</span>
              </div>
              <div class="match-meta">{{ formatKickoff(match.kickoff) }}</div>
            </template>
          </div>
          <button v-if="!savedPredictions[match.id]" class="btn-lock" @click="submit(match.id)">Predict</button>
        </div>
      </div>
    </section>

    <section v-if="completedMatches.length > 0" data-testid="results-section">
      <h2 class="page-title" style="margin-top:2.5rem">Results</h2>
      <div class="match-list">
        <div
          v-for="match in completedMatches"
          :key="match.id"
          class="match-card match-card--result"
          data-testid="result-card"
        >
          <div class="match-info">
            <div class="matchup matchup--input" data-testid="matchup">
              <span class="team-name team-home">{{ match.homeTeam }}</span>
              <span class="inline-score">{{ match.homeScore || '–' }}</span>
              <span class="vs">vs</span>
              <span class="inline-score">{{ match.awayScore || '–' }}</span>
              <span class="team-name team-away">{{ match.awayTeam }}</span>
            </div>
            <div class="match-meta">{{ formatKickoff(match.kickoff) }}</div>
            <div v-if="savedPredictions[match.id]" class="match-meta">
              Your pick: {{ savedPredictions[match.id]!.homeGoals }} – {{ savedPredictions[match.id]!.awayGoals }}
            </div>
          </div>
          <div class="btn-spacer"></div>
        </div>
      </div>
    </section>

    <p v-if="upcomingMatches.length === 0 && completedMatches.length === 0" class="empty">No matches found. Check back later.</p>

    <div v-if="showNudge && !currentUser" class="guest-nudge" data-testid="guest-nudge">
      <a href="/signup">Create an account</a> to get on the leaderboard.
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, inject, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import type { Ref } from 'vue'

const currentUser = inject<Ref<{ handle: string; emailVerified: boolean } | null>>('currentUser')
const router = useRouter()

interface Match {
  id: string
  homeTeam: string
  awayTeam: string
  kickoff: string
  status: string
  homeScore: string
  awayScore: string
  state?: string
}

interface Prediction {
  homeGoals: number
  awayGoals: number
}

const GUEST_KEY = 'guestPredictions'

const matches = ref<Match[]>([])
const savedPredictions = reactive<Record<string, Prediction | null>>({})
const inputs = reactive<Record<string, { home: string; away: string }>>({})
const showNudge = ref(false)

const upcomingMatches = computed(() => {
  const cutoff = new Date()
  cutoff.setDate(cutoff.getDate() + 7)
  return matches.value.filter(m => {
    if (m.status !== 'STATUS_SCHEDULED' && m.status !== 'STATUS_IN_PROGRESS') return false
    return new Date(m.kickoff) <= cutoff
  })
})

const completedMatches = computed(() =>
  matches.value
    .filter(m => m.status !== 'STATUS_SCHEDULED' && m.status !== 'STATUS_IN_PROGRESS')
    .reverse()
)

function formatKickoff(iso: string): string {
  const d = new Date(iso)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric', hour: 'numeric', minute: '2-digit', timeZoneName: 'short' })
}

onMounted(async () => {
  document.title = 'Upcoming — Crew Predictions'
  const res = await fetch('/api/matches')
  if (res.ok) {
    const data = await res.json()
    matches.value = data.matches
    const guestPredictions: Record<string, Prediction> = JSON.parse(localStorage.getItem(GUEST_KEY) ?? '{}')
    for (const m of data.matches) {
      inputs[m.id] = { home: '', away: '' }
      savedPredictions[m.id] = data.predictions[m.id] ?? guestPredictions[m.id] ?? null
    }
  }
})

async function submit(matchId: string) {
  const { home, away } = inputs[matchId]
  if (!currentUser?.value) {
    const prediction = { homeGoals: Number(home), awayGoals: Number(away) }
    const stored: Record<string, Prediction> = JSON.parse(localStorage.getItem(GUEST_KEY) ?? '{}')
    stored[matchId] = prediction
    localStorage.setItem(GUEST_KEY, JSON.stringify(stored))
    savedPredictions[matchId] = prediction
    showNudge.value = true
    return
  }
  const body = new URLSearchParams({ match_id: matchId, home_goals: home, away_goals: away })
  const res = await fetch('/api/predictions', { method: 'POST', body })
  if (res.ok) {
    savedPredictions[matchId] = { homeGoals: Number(home), awayGoals: Number(away) }
  }
}
</script>
