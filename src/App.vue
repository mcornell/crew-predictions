<template>
  <AppHeader :user="currentUser" :loading="authLoading" />
<RouterView />
</template>

<script setup lang="ts">
import { ref, provide, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppHeader from './components/AppHeader.vue'
import { getGoogleRedirectResult } from './firebase'

const route = useRoute()
const router = useRouter()
const currentUser = ref<{ userID: string; handle: string; emailVerified: boolean } | null>(null)
const authLoading = ref(true)

provide('currentUser', currentUser)

async function fetchUser() {
  const res = await fetch('/api/me')
  currentUser.value = res.ok ? await res.json() : null
  authLoading.value = false
}

onMounted(async () => {
  try {
    const token = await getGoogleRedirectResult()
    if (token) {
      await fetch('/auth/session', {
        method: 'POST',
        body: new URLSearchParams({ idToken: token }),
      })
      await fetchUser()
      router.push('/matches')
      return
    }
  } catch {
    // redirect result failed — fall through and fetch current session state
  }
  await fetchUser()
})

watch(() => route.path, fetchUser)
</script>
