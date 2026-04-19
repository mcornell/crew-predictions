<template>
  <h1>Upcoming</h1>
  <div
    v-for="match in matches"
    :key="match.id"
    data-testid="match-card"
  >
    <span>{{ match.homeTeam }} vs {{ match.awayTeam }}</span>
    <template v-if="savedPredictions[match.id]">
      <span>{{ savedPredictions[match.id]!.homeGoals }} – {{ savedPredictions[match.id]!.awayGoals }}</span>
    </template>
    <template v-else>
      <input name="home_goals" type="number" v-model="inputs[match.id].home" />
      <input name="away_goals" type="number" v-model="inputs[match.id].away" />
      <button @click="submit(match.id)">Lock In</button>
    </template>
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
