<template>
  <AppHeader :user="currentUser" />
  <RouterView />
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppHeader from './components/AppHeader.vue'

const route = useRoute()
const currentUser = ref<{ handle: string } | null>(null)

async function fetchUser() {
  const res = await fetch('/api/me')
  currentUser.value = res.ok ? await res.json() : null
}

onMounted(fetchUser)
watch(() => route.path, fetchUser)
</script>
