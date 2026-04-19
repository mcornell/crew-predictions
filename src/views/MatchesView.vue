<template>
  <h1>Upcoming</h1>
  <div
    v-for="match in matches"
    :key="match.id"
    data-testid="match-card"
  >
    {{ match.homeTeam }} vs {{ match.awayTeam }}
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface Match {
  id: string
  homeTeam: string
  awayTeam: string
  kickoff: string
  status: string
}

const matches = ref<Match[]>([])

onMounted(async () => {
  const res = await fetch('/api/matches')
  if (res.ok) {
    const data = await res.json()
    matches.value = data.matches
  }
})
</script>
