<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="login-title">Profile</h1>
      <form class="login-form" data-testid="profile-form" @submit.prevent="handleSubmit">
        <input
          class="form-input"
          data-testid="display-name-input"
          v-model="displayName"
          type="text"
          placeholder="Display name"
          required
        />
        <p v-if="error" class="form-error">{{ error }}</p>
        <button class="btn-submit" type="submit">Save</button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, inject, watch } from 'vue'
import { useRouter } from 'vue-router'
import { updateDisplayName } from '../firebase'
import type { Ref } from 'vue'

const currentUser = inject<Ref<{ handle: string; emailVerified: boolean } | null>>('currentUser')
const router = useRouter()
const displayName = ref('')
const error = ref('')

watch(
  () => currentUser?.value,
  (user) => { if (user && !displayName.value) displayName.value = user.handle },
  { immediate: true }
)

async function handleSubmit() {
  error.value = ''
  try {
    await updateDisplayName(displayName.value)
    router.push('/matches')
  } catch {
    error.value = 'Could not save display name. Please try again.'
  }
}
</script>
