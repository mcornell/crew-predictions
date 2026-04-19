<template>
  <AppHeader :user="currentUser" />
  <RouterView />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AppHeader from './components/AppHeader.vue'

const currentUser = ref<{ handle: string } | null>(null)

onMounted(async () => {
  const res = await fetch('/api/me')
  if (res.ok) {
    const data = await res.json()
    currentUser.value = { handle: data.handle }
  }
})
</script>
