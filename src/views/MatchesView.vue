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
            <template v-if="savedPredictions[match.id]">
              <div class="matchup" data-testid="matchup">
                {{ match.homeTeam }}
                <span class="inline-score">{{ savedPredictions[match.id]!.homeGoals }}</span>
                <span class="vs">vs</span>
                <span class="inline-score">{{ savedPredictions[match.id]!.awayGoals }}</span>
                {{ match.awayTeam }}
              </div>
              <div class="match-meta">{{ formatKickoff(match.kickoff) }} — <span class="saved-label">Locked in</span></div>
            </template>
            <template v-else>
              <div class="matchup matchup--input" data-testid="matchup">
                <span class="team-name">{{ match.homeTeam }}</span>
                <input class="score-input" name="home_goals" type="number" min="0" max="99" v-model="inputs[match.id].home" placeholder="0" />
                <span class="vs">vs</span>
                <input class="score-input" name="away_goals" type="number" min="0" max="99" v-model="inputs[match.id].away" placeholder="0" />
                <span class="team-name">{{ match.awayTeam }}</span>
                <button class="btn-lock" @click="submit(match.id)">Predict</button>
              </div>
              <div class="match-meta">{{ formatKickoff(match.kickoff) }}</div>
            </template>
          </div>
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
            <div class="matchup" data-testid="matchup">
              {{ match.homeTeam }}
              <span class="inline-score" v-if="match.homeScore">{{ match.homeScore }}</span>
              <span class="vs">vs</span>
              <span class="inline-score" v-if="match.awayScore">{{ match.awayScore }}</span>
              {{ match.awayTeam }}
            </div>
            <div class="match-meta">{{ formatKickoff(match.kickoff) }}</div>
            <div v-if="savedPredictions[match.id]" class="match-meta">
              Your pick: {{ savedPredictions[match.id]!.homeGoals }} – {{ savedPredictions[match.id]!.awayGoals }}
            </div>
          </div>
        </div>
      </div>
    </section>

    <p v-if="upcomingMatches.length === 0 && completedMatches.length === 0" class="empty">No matches found. Check back later.</p>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'

interface Match {
  id: string
  homeTeam: string
  awayTeam: string
  kickoff: string
  status: string
  homeScore: string
  awayScore: string
}

interface Prediction {
  homeGoals: number
  awayGoals: number
}

const matches = ref<Match[]>([])
const savedPredictions = reactive<Record<string, Prediction | null>>({})
const inputs = reactive<Record<string, { home: string; away: string }>>({})

const upcomingMatches = computed(() => {
  const cutoff = new Date()
  cutoff.setDate(cutoff.getDate() + 7)
  return matches.value.filter(m => {
    if (m.status !== 'STATUS_SCHEDULED' && m.status !== 'STATUS_IN_PROGRESS') return false
    return new Date(m.kickoff) <= cutoff
  })
})

const completedMatches = computed(() =>
  matches.value.filter(m => m.status !== 'STATUS_SCHEDULED' && m.status !== 'STATUS_IN_PROGRESS')
)

function formatKickoff(iso: string): string {
  const d = new Date(iso)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric', hour: 'numeric', minute: '2-digit', timeZoneName: 'short' })
}

onMounted(async () => {
  const res = await fetch('/api/matches')
  if (res.ok) {
    const data = await res.json()
    matches.value = data.matches
    for (const m of data.matches) {
      inputs[m.id] = { home: '', away: '' }
      savedPredictions[m.id] = data.predictions[m.id] ?? null
    }
  }
})

async function submit(matchId: string) {
  const { home, away } = inputs[matchId]
  const body = new URLSearchParams({ match_id: matchId, home_goals: home, away_goals: away })
  const res = await fetch('/api/predictions', { method: 'POST', body })
  if (res.ok) {
    savedPredictions[matchId] = { homeGoals: Number(home), awayGoals: Number(away) }
  }
}
</script>
