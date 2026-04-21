<template>
  <header class="site-header">
    <a class="brand" href="/">Crew Predictions</a>
    <nav class="site-nav">
      <a class="btn-ghost" href="/leaderboard">Leaderboard</a>
      <a class="btn-ghost" href="/rules">Rules</a>
      <template v-if="user">
        <a class="btn-ghost" href="/profile">{{ user.handle }}</a>
        <a class="btn-ghost" href="/auth/logout">Sign out</a>
      </template>
      <a v-else-if="!loading" class="btn-primary" href="/login">Sign In</a>
    </nav>
  </header>
</template>

<script setup lang="ts">
defineProps<{
  user: { handle: string } | null
  loading?: boolean
}>()
</script>
