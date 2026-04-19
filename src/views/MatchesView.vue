<template>
  <div class="page">
    <h1 class="page-title">Upcoming</h1>
    <div class="match-list">
      <div
        v-for="match in matches"
        :key="match.id"
        class="match-card"
        :class="{ 'has-prediction': savedPredictions[match.id] }"
        data-testid="match-card"
      >
        <div class="match-info">
          <div class="matchup">{{ match.homeTeam }} <span style="color:var(--muted)">vs</span> {{ match.awayTeam }}</div>
          <div class="match-meta">{{ formatKickoff(match.kickoff) }}</div>
        </div>
        <div class="prediction">
          <template v-if="savedPredictions[match.id]">
            <div class="saved-score">
              <span class="score-display">{{ savedPredictions[match.id]!.homeGoals }} – {{ savedPredictions[match.id]!.awayGoals }}</span>
              <span class="saved-label">Your Pick</span>
            </div>
          </template>
          <template v-else>
            <div class="score-inputs">
              <input class="score-input" name="home_goals" type="number" min="0" max="99" v-model="inputs[match.id].home" placeholder="0" />
              <span class="score-sep">–</span>
              <input class="score-input" name="away_goals" type="number" min="0" max="99" v-model="inputs[match.id].away" placeholder="0" />
            </div>
            <button class="btn-lock" @click="submit(match.id)">Lock In</button>
          </template>
        </div>
      </div>
    </div>
    <p v-if="matches.length === 0" class="empty">No upcoming matches. Check back later.</p>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'

interface Match {
  id: string
  homeTeam: string
  awayTeam: string
  kickoff: string
  status: string
}

interface Prediction {
  homeGoals: number
  awayGoals: number
}

const matches = ref<Match[]>([])
const savedPredictions = reactive<Record<string, Prediction | null>>({})
const inputs = reactive<Record<string, { home: string; away: string }>>({})

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
