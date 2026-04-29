<template>
  <AppHeader :user="currentUser" :loading="authLoading" :seasons="seasons" />
<RouterView />
</template>

<script setup lang="ts">
import { ref, provide, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppHeader from './components/AppHeader.vue'
import { getGoogleRedirectResult } from './firebase'
import { flushGuestPredictions } from './guestPredictions'

const route = useRoute()
const router = useRouter()
const currentUser = ref<{ userID: string; handle: string; emailVerified: boolean } | null>(null)
const authLoading = ref(true)
const seasons = ref<{ id: string; name: string; isCurrent: boolean }[]>([])

provide('currentUser', currentUser)

async function fetchUser() {
  const res = await fetch('/api/me')
  currentUser.value = res.ok ? await res.json() : null
  authLoading.value = false
}

async function fetchSeasons() {
  try {
    const res = await fetch('/api/seasons')
    if (res.ok) {
      const data = await res.json()
      const all: { id: string; name: string; isCurrent: boolean }[] = data.seasons ?? []
      const currentIdx = all.findIndex(s => s.isCurrent)
      const past = currentIdx >= 0 ? all.slice(0, currentIdx) : []
      seasons.value = [...past].reverse()
    }
  } catch { /* non-critical */ }
}

onMounted(async () => {
  try {
    const token = await getGoogleRedirectResult()
    if (token) {
      const res = await fetch('/auth/session', {
        method: 'POST',
        body: new URLSearchParams({ idToken: token }),
      })
      if (res.ok) {
        await flushGuestPredictions()
        await fetchUser()
        router.push('/matches')
        return
      }
    }
  } catch (err) {
    console.error('Google redirect result failed:', err)
  }
  await fetchUser()
  await fetchSeasons()
})

watch(() => route.path, fetchUser)
</script>
