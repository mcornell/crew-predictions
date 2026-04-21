<template>
  <AppHeader :user="currentUser" />
  <div v-if="currentUser && !currentUser.emailVerified" data-testid="email-verification-banner" class="verification-banner">
    Please verify your email — check your inbox for the verification link.
  </div>
  <RouterView />
</template>

<script setup lang="ts">
import { ref, provide, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppHeader from './components/AppHeader.vue'
import { getGoogleRedirectResult } from './firebase'

const route = useRoute()
const router = useRouter()
const currentUser = ref<{ handle: string; emailVerified: boolean } | null>(null)

provide('currentUser', currentUser)

async function fetchUser() {
  const res = await fetch('/api/me')
  currentUser.value = res.ok ? await res.json() : null
}

onMounted(async () => {
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
  await fetchUser()
})

watch(() => route.path, fetchUser)
</script>
