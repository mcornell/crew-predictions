<template>
  <div class="page">
    <p v-if="loading" data-testid="loading" class="status-msg">Loading…</p>
    <p v-else-if="error" data-testid="error" class="status-msg status-msg--error">{{ error }}</p>
    <template v-else>
    <section v-if="nowPlayingMatches.length > 0" data-testid="now-playing-section" class="now-playing-section">
      <h1 class="page-title">Now Playing</h1>
      <div class="match-list">
        <RouterLink
          v-for="match in nowPlayingMatches"
          :key="match.id"
          :to="`/matches/${match.id}`"
          class="match-card match-card--live match-card--link"
          :data-match-id="match.id"
          data-testid="now-playing-card"
        >
          <div class="match-info">
            <span v-if="match.status === 'STATUS_DELAYED'" class="delayed-indicator" data-testid="delayed-indicator">▊ DELAYED</span>
            <span v-else-if="match.status === 'STATUS_HALFTIME'" class="live-indicator" data-testid="live-indicator">● HT</span>
            <span v-else class="live-indicator" data-testid="live-indicator">● {{ match.displayClock || 'LIVE' }}</span>
            <div class="matchup matchup--input" data-testid="matchup">
              <span class="team-name team-home">{{ match.homeTeam }}</span>
              <span class="inline-score">{{ match.homeScore || '0' }}</span>
              <span class="vs">vs</span>
              <span class="inline-score">{{ match.awayScore || '0' }}</span>
              <span class="team-name team-away">{{ match.awayTeam }}</span>
            </div>
            <div class="match-meta">{{ formatKickoff(match.kickoff) }}</div>
          </div>
          <div class="btn-spacer"></div>
        </RouterLink>
      </div>
    </section>

    <section v-if="upcomingMatches.length > 0">
      <h1 class="page-title" :style="nowPlayingMatches.length > 0 ? 'margin-top:2.5rem' : ''">Upcoming</h1>
      <div class="match-list">
        <div
          v-for="match in upcomingMatches"
          :key="match.id"
          class="match-card"
          :class="{ 'has-prediction': savedPredictions[match.id] }"
          :data-match-id="match.id"
          data-testid="match-card"
        >
          <div class="match-info">
            <span class="match-countdown" data-testid="match-countdown">{{ countdowns[match.id] }}</span>
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
          <template v-if="!isLocked(match)">
            <button v-if="savedPredictions[match.id]" class="btn-lock btn-unlock" @click="unlock(match.id)">Unlock</button>
            <button v-else class="btn-lock" @click="submit(match.id)">Predict</button>
          </template>
          <div v-else class="btn-spacer"></div>
        </div>
      </div>
    </section>

    <section v-if="completedMatches.length > 0" data-testid="results-section">
      <h2 class="page-title" style="margin-top:2.5rem">Results</h2>
      <div class="match-list">
        <RouterLink
          v-for="match in completedMatches"
          :key="match.id"
          :to="`/matches/${match.id}`"
          class="match-card match-card--result match-card--link"
          :data-match-id="match.id"
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
        </RouterLink>
      </div>
    </section>

    <p v-if="nowPlayingMatches.length === 0 && upcomingMatches.length === 0 && completedMatches.length === 0" class="empty">No matches found. Check back later.</p>

    <div v-if="showNudge && !currentUser" class="guest-nudge" data-testid="guest-nudge">
      <a href="/signup">Create an account</a> to get on the leaderboard.
    </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, inject, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { formatCountdown } from '../utils/countdown'
import { isInActiveWindow, msUntilActiveWindow, POLL_INTERVAL_MS } from '../utils/pollScheduler'
import { readGuestPredictions, writeGuestPredictions } from '../guestPredictions'
import type { Ref } from 'vue'

const currentUser = inject<Ref<{ handle: string; emailVerified: boolean } | null>>('currentUser', ref(null))
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
  displayClock?: string
}

interface Prediction {
  homeGoals: number
  awayGoals: number
}

const matches = ref<Match[]>([])
const savedPredictions = reactive<Record<string, Prediction | null>>({})
const inputs = reactive<Record<string, { home: string; away: string }>>({})
const countdowns = reactive<Record<string, string>>({})
const showNudge = ref(false)
const nowMs = ref(Date.now())
const loading = ref(true)
const error = ref<string | null>(null)

const nowPlayingMatches = computed(() =>
  matches.value.filter(m => m.state === 'in' || m.status === 'STATUS_DELAYED')
)

const upcomingMatches = computed(() => {
  const cutoff = new Date()
  cutoff.setDate(cutoff.getDate() + 8)
  return matches.value.filter(m => m.state === 'pre' && new Date(m.kickoff) <= cutoff)
})

const completedMatches = computed(() =>
  matches.value
    .filter(m => m.state === 'post')
    .sort((a, b) => new Date(b.kickoff).getTime() - new Date(a.kickoff).getTime())
)

function isLocked(match: Match): boolean {
  return nowMs.value >= new Date(match.kickoff).getTime() || match.state === 'in' || match.status === 'STATUS_DELAYED'
}

function updateCountdowns() {
  nowMs.value = Date.now()
  for (const m of upcomingMatches.value) {
    countdowns[m.id] = formatCountdown(new Date(m.kickoff).getTime() - nowMs.value)
  }
}

let countdownTimer: ReturnType<typeof setInterval> | null = null
let pollTimer: ReturnType<typeof setTimeout> | null = null

function schedulePoll() {
  if (pollTimer !== null) clearTimeout(pollTimer)
  if (isInActiveWindow(matches.value, Date.now())) {
    pollTimer = setTimeout(async () => {
      await fetchMatches()
      schedulePoll()
    }, POLL_INTERVAL_MS)
  } else {
    const ms = msUntilActiveWindow(matches.value, Date.now())
    if (ms !== null) {
      pollTimer = setTimeout(() => {
        fetchMatches().then(() => schedulePoll())
      }, ms)
    }
  }
}

function formatKickoff(iso: string): string {
  const d = new Date(iso)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric', hour: 'numeric', minute: '2-digit', timeZoneName: 'short' })
}

async function fetchMatches() {
  const res = await fetch('/api/matches')
  if (!res.ok) return
  const data = await res.json()
  matches.value = data.matches
  const guestPredictions = readGuestPredictions()
  for (const m of data.matches) {
    if (!inputs[m.id]) inputs[m.id] = { home: '', away: '' }
    if (savedPredictions[m.id] === undefined) {
      savedPredictions[m.id] = data.predictions[m.id] ?? guestPredictions[m.id] ?? null
    }
  }
}

onMounted(async () => {
  document.title = 'Upcoming — Crew Predictions'
  const res = await fetch('/api/matches')
  if (res.ok) {
    const data = await res.json()
    matches.value = data.matches
    const guestPredictions = readGuestPredictions()
    for (const m of data.matches) {
      inputs[m.id] = { home: '', away: '' }
      savedPredictions[m.id] = data.predictions[m.id] ?? guestPredictions[m.id] ?? null
    }
    updateCountdowns()
    countdownTimer = setInterval(updateCountdowns, 1000)
  } else {
    error.value = 'Could not load matches. Try again later.'
  }
  loading.value = false
  schedulePoll()
})

onUnmounted(() => {
  if (countdownTimer !== null) clearInterval(countdownTimer)
  if (pollTimer !== null) clearTimeout(pollTimer)
})

function unlock(matchId: string) {
  const prev = savedPredictions[matchId]
  if (prev) {
    inputs[matchId] = { home: String(prev.homeGoals), away: String(prev.awayGoals) }
  }
  savedPredictions[matchId] = null
}

async function submit(matchId: string) {
  const { home, away } = inputs[matchId]
  if (!currentUser?.value) {
    const prediction = { homeGoals: Number(home), awayGoals: Number(away) }
    const stored = readGuestPredictions()
    stored[matchId] = prediction
    writeGuestPredictions(stored)
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
