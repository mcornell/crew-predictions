<template>
  <AppHeader :user="currentUser" :loading="authLoading" />
  <div v-if="currentUser && !currentUser.emailVerified" data-testid="email-verification-banner" class="verification-banner">
    Please verify your email — check your inbox for the verification link.
  </div>
  <RouterView />
</template>

<script setup lang="ts">
import { ref, provide, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppHeader from './components/AppHeader.vue'

const route = useRoute()
const currentUser = ref<{ handle: string; emailVerified: boolean } | null>(null)
const authLoading = ref(true)

provide('currentUser', currentUser)

async function fetchUser() {
  const res = await fetch('/api/me')
  currentUser.value = res.ok ? await res.json() : null
  authLoading.value = false
}

onMounted(fetchUser)

watch(() => route.path, fetchUser)
</script>
