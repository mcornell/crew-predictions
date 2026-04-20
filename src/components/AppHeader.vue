<template>
  <header class="site-header">
    <a class="brand" href="/">Crew Predictions</a>
    <nav class="site-nav">
      <template v-if="user">
        <a class="btn-ghost" href="/profile">{{ user.handle }}</a>
        <a class="btn-ghost" href="/auth/logout">Sign out</a>
      </template>
      <a v-else class="btn-primary" href="/login">Sign In</a>
    </nav>
  </header>
</template>

<script setup lang="ts">
defineProps<{
  user: { handle: string } | null
}>()
</script>
