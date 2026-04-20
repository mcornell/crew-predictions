<template>
  <header class="site-header">
    <a class="brand" href="/">Crew Predictions</a>
    <nav class="site-nav">
      <template v-if="user">
        <span class="user-handle">{{ user.handle }}</span>
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
