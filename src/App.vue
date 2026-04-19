<template>
  <AppHeader :user="currentUser" />
  <RouterView />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AppHeader from './components/AppHeader.vue'

const currentUser = ref<{ handle: string } | null>(null)
</script>
