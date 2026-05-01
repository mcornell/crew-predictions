<template>
  <div class="page">
    <RouterLink to="/matches" class="back-link" data-testid="back-link">← All Matches</RouterLink>

    <p v-if="loading" data-testid="loading" class="status-msg">Loading…</p>
    <p v-else-if="error" data-testid="error" class="status-msg status-msg--error">{{ error }}</p>

    <div v-if="match" class="match-detail-header">
      <div v-if="isLive" class="match-detail-live-bar" data-testid="live-indicator-detail">
        <span v-if="match.status === 'STATUS_HALFTIME'" class="live-indicator">● HT</span>
        <span v-else class="live-indicator">● {{ match.displayClock || 'LIVE' }}</span>
      </div>
      <div class="matchup matchup--input" data-testid="match-score" :class="{ 'matchup--has-form': match.homeRecord || match.awayRecord || match.homeForm || match.awayForm }">
        <span class="team-name team-home">
          <img v-if="match.homeLogo" :src="match.homeLogo" :alt="`${match.homeTeam} logo`" class="team-logo" data-testid="home-logo" loading="lazy" />
          <span class="team-label">{{ match.homeTeam }}</span>
        </span>
        <span class="inline-score">{{ match.homeScore || '–' }}</span>
        <span class="vs">vs</span>
        <span class="inline-score">{{ match.awayScore || '–' }}</span>
        <span class="team-name team-away">
          <img v-if="match.awayLogo" :src="match.awayLogo" :alt="`${match.awayTeam} logo`" class="team-logo" data-testid="away-logo" loading="lazy" />
          <span class="team-label">{{ match.awayTeam }}</span>
        </span>
        <div v-if="match.homeRecord || match.homeForm" class="matchup-team-form matchup-team-form--home">
          <span v-if="match.homeRecord" class="match-record">{{ match.homeRecord }}</span>
          <span v-if="match.homeForm" class="match-form"><span v-for="(c, i) in match.homeForm.split('')" :key="i" :class="`form-letter form-letter--${c.toLowerCase()}`">{{ c }}</span></span>
        </div>
        <div v-if="match.awayRecord || match.awayForm" class="matchup-team-form matchup-team-form--away">
          <span v-if="match.awayRecord" class="match-record">{{ match.awayRecord }}</span>
          <span v-if="match.awayForm" class="match-form"><span v-for="(c, i) in match.awayForm.split('')" :key="i" :class="`form-letter form-letter--${c.toLowerCase()}`">{{ c }}</span></span>
        </div>
      </div>
      <div v-if="match.homeRecord || match.homeForm" class="match-form-row">
        <div class="match-form-team">
          <span v-if="match.homeRecord" class="match-record" data-testid="home-record">{{ match.homeRecord }}</span>
          <span v-if="match.homeForm" class="match-form" data-testid="home-form"><span v-for="(c, i) in match.homeForm.split('')" :key="i" :class="`form-letter form-letter--${c.toLowerCase()}`">{{ c }}</span></span>
        </div>
        <div class="match-form-team match-form-team--away">
          <span v-if="match.awayForm" class="match-form" data-testid="away-form"><span v-for="(c, i) in match.awayForm.split('')" :key="i" :class="`form-letter form-letter--${c.toLowerCase()}`">{{ c }}</span></span>
          <span v-if="match.awayRecord" class="match-record" data-testid="away-record">{{ match.awayRecord }}</span>
        </div>
      </div>
      <div class="match-meta">{{ formatKickoff(match.kickoff) }}</div>
      <div v-if="match.venue" class="match-meta match-venue" data-testid="match-detail-venue">{{ match.venue }}</div>
      <div v-if="match.attendance" class="match-meta match-attendance" data-testid="match-detail-attendance">{{ formatAttendance(match.attendance) }}</div>
      <div v-if="match.referee" class="match-meta match-referee" data-testid="match-referee">Referee: {{ match.referee }}</div>
      <div v-if="displayableEvents.length > 0" class="match-events" data-testid="match-events">
        <div
          v-for="(event, i) in displayableEvents"
          :key="i"
          class="match-event"
          :class="[`match-event--${event.typeID}`, `match-event--${eventSide(event)}`]"
          data-testid="match-event"
        >
          <div class="event-detail">
            <span class="event-icon" :aria-label="event.typeID">{{ eventIcon(event.typeID) }}</span>
            <span v-if="event.typeID === 'substitution'" class="event-players">
              <span class="sub-on">
                <span class="full-name">{{ event.players[0] }}</span>
                <span class="short-name">{{ surname(event.players[0] || '') }}</span>
                <span class="sub-arrow">↑</span>
              </span>
              <span v-if="event.players[1]" class="sub-off">
                <span class="full-name">{{ event.players[1] }}</span>
                <span class="short-name">{{ surname(event.players[1]) }}</span>
                <span class="sub-arrow">↓</span>
              </span>
            </span>
            <span v-else class="event-players">
              <span class="full-name">{{ event.players.join(', ') }}</span>
              <span class="short-name">{{ event.players.map(surname).join(', ') }}</span>
            </span>
          </div>
          <span class="event-clock">{{ event.clock }}</span>
        </div>
      </div>
      <a
        :href="`https://www.espn.com/soccer/match/_/gameId/${match.id}`"
        target="_blank"
        rel="noopener noreferrer"
        class="espn-link"
        data-testid="espn-link"
      >View on ESPN ↗</a>
    </div>

    <template v-if="sortedPredictions.length > 0">
      <p v-if="isProjected" class="projected-label" data-testid="projected-label">Projected scores based on current live result</p>
      <div class="lb-mobile-sort">
        <button
          class="lb-sort-btn"
          :class="{ 'lb-sort-btn--active': activeFormat === 'acesRadio' }"
          data-testid="mobile-sort-aces"
          @click="activeFormat = 'acesRadio'"
        >Aces Radio</button>
        <button
          class="lb-sort-btn"
          :class="{ 'lb-sort-btn--active': activeFormat === 'upper90Club' }"
          data-testid="mobile-sort-upper90"
          @click="activeFormat = 'upper90Club'"
        >Upper 90 Club</button>
        <button
          class="lb-sort-btn"
          :class="{ 'lb-sort-btn--active': activeFormat === 'grouchy' }"
          data-testid="mobile-sort-grouchy"
          @click="activeFormat = 'grouchy'"
        >Grouchy™</button>
      </div>

    <div class="lb-table lb-5col">
      <div class="lb-header">
        <span class="lb-cell lb-rank">RANK</span>
        <span class="lb-cell lb-handle">PREDICTOR</span>
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
        <button
          class="lb-cell lb-pts lb-sort-btn"
          :class="{ 'lb-sort-btn--active': activeFormat === 'grouchy' }"
          data-testid="sort-grouchy"
          @click="activeFormat = 'grouchy'"
        >Grouchy™</button>
      </div>

      <div
        v-for="(entry, i) in sortedPredictions"
        :key="entry.userID"
        class="lb-row"
        data-testid="prediction-row"
      >
        <span class="lb-cell lb-rank" data-testid="prediction-rank">{{ rankFor(i) }}</span>
        <span class="lb-cell lb-handle">
          <span class="lb-handle-name">{{ entry.handle }}</span>
          <span class="lb-pick-sub" data-testid="prediction-score">{{ entry.homeGoals }} – {{ entry.awayGoals }}</span>
        </span>
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
        <span
          class="lb-cell lb-pts"
          :class="{ 'lb-pts--active': activeFormat === 'grouchy' }"
          data-testid="prediction-grouchy-points"
          data-label="Grouchy™"
        >{{ entry.grouchyPoints }}</span>
      </div>
    </div>
    </template>
    <p v-else-if="loaded" class="empty">No predictions were made for this match</p>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { isInActiveWindow, msUntilActiveWindow, POLL_INTERVAL_MS } from '../utils/pollScheduler'

interface MatchEvent {
  clock: string
  typeID: string
  team: string
  players: string[]
}

interface MatchInfo {
  id: string
  homeTeam: string
  awayTeam: string
  kickoff: string
  homeScore: string
  awayScore: string
  state?: string
  displayClock?: string
  status?: string
  venue?: string
  homeRecord?: string
  awayRecord?: string
  homeForm?: string
  awayForm?: string
  homeLogo?: string
  awayLogo?: string
  attendance?: number
  referee?: string
  events?: MatchEvent[]
}

const NON_DISPLAYABLE_EVENTS = new Set(['kickoff', 'halftime', 'start-2nd-half', 'end-regular-time'])

const GOAL_TYPES = new Set(['goal', 'goal---header', 'goal---volley', 'own-goal', 'penalty---scored'])

function eventIcon(typeID: string): string {
  if (GOAL_TYPES.has(typeID)) return '⚽'
  if (typeID === 'yellow-card') return '🟨'
  if (typeID === 'red-card') return '🟥'
  if (typeID === 'penalty---saved') return '🧤'
  if (typeID === 'substitution') return '🔄'
  return ''
}

function surname(fullName: string): string {
  const parts = fullName.trim().split(/\s+/)
  return parts[parts.length - 1] || fullName
}

interface PredictionEntry {
  userID: string
  handle: string
  homeGoals: number
  awayGoals: number
  acesRadioPoints: number
  upper90ClubPoints: number
  grouchyPoints: number
}

const route = useRoute()
const match = ref<MatchInfo | null>(null)
const predictions = ref<PredictionEntry[]>([])
const isProjected = ref(false)
const activeFormat = ref<'acesRadio' | 'upper90Club' | 'grouchy'>('acesRadio')
const loaded = ref(false)
const loading = ref(true)
const error = ref<string | null>(null)

const isLive = computed(() => match.value?.state === 'in')

const displayableEvents = computed(() =>
  (match.value?.events ?? []).filter(e => !NON_DISPLAYABLE_EVENTS.has(e.typeID))
)

function eventSide(event: MatchEvent): 'home' | 'away' {
  return match.value && event.team === match.value.homeTeam ? 'home' : 'away'
}

const sortedPredictions = computed(() => {
  const key = activeFormat.value === 'acesRadio' ? 'acesRadioPoints' : activeFormat.value === 'upper90Club' ? 'upper90ClubPoints' : 'grouchyPoints'
  return [...predictions.value].sort((a, b) => b[key] - a[key])
})

function rankFor(i: number): number {
  if (i === 0) return 1
  const key = activeFormat.value === 'acesRadio' ? 'acesRadioPoints' : activeFormat.value === 'upper90Club' ? 'upper90ClubPoints' : 'grouchyPoints'
  if (sortedPredictions.value[i][key] === sortedPredictions.value[i - 1][key]) return rankFor(i - 1)
  return i + 1
}

function formatAttendance(n: number): string {
  return n.toLocaleString('en-US')
}

function formatKickoff(iso: string): string {
  const d = new Date(iso)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric', hour: 'numeric', minute: '2-digit', timeZoneName: 'short' })
}

const matchId = route.params.matchId as string
let pollTimer: ReturnType<typeof setTimeout> | null = null

async function fetchDetail() {
  const res = await fetch(`/api/matches/${matchId}`)
  if (!res.ok) return
  const data = await res.json()
  match.value = data.match
  predictions.value = data.predictions ?? []
  isProjected.value = data.isProjected ?? false
}

function schedulePoll() {
  if (pollTimer !== null) clearTimeout(pollTimer)
  const m = match.value
  if (!m) return
  const asSchedulerMatch = [{ kickoff: m.kickoff, state: m.state, status: '' }]
  if (isInActiveWindow(asSchedulerMatch, Date.now())) {
    pollTimer = setTimeout(async () => {
      await fetchDetail()
      schedulePoll()
    }, POLL_INTERVAL_MS)
  } else {
    const ms = msUntilActiveWindow(asSchedulerMatch, Date.now())
    if (ms !== null) {
      pollTimer = setTimeout(async () => {
        await fetchDetail()
        schedulePoll()
      }, ms)
    }
  }
}

onMounted(async () => {
  const res = await fetch(`/api/matches/${matchId}`)
  if (res.ok) {
    const data = await res.json()
    match.value = data.match
    predictions.value = data.predictions ?? []
    isProjected.value = data.isProjected ?? false
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
  schedulePoll()
})

onUnmounted(() => {
  if (pollTimer !== null) clearTimeout(pollTimer)
})
</script>
