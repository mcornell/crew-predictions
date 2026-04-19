<template>
  <header>
    <span>Crew Predictions</span>
    <nav>
      <template v-if="user">
        <span>{{ user.handle }}</span>
        <a href="/logout">Sign out</a>
      </template>
      <a v-else href="/login">Sign In</a>
    </nav>
  </header>
</template>

<script setup lang="ts">
defineProps<{
  user: { handle: string } | null
}>()
</script>
