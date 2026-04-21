<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="login-title">Reset Password</h1>
      <template v-if="sent">
        <p data-testid="reset-confirmation" class="login-sub">
          Check your email — a reset link is on its way.
        </p>
        <p class="auth-alt">
          <router-link to="/login">Back to sign in</router-link>
        </p>
      </template>
      <template v-else>
        <p class="login-sub">Enter your email and we'll send a reset link.</p>
        <form class="login-form" data-testid="reset-form" @submit.prevent="handleSubmit">
          <input class="form-input" v-model="email" type="email" placeholder="Email" autocomplete="email" required />
          <p v-if="error" class="form-error">{{ error }}</p>
          <button class="btn-submit" type="submit">Send Reset Link</button>
        </form>
        <p class="auth-alt">
          <router-link to="/login">Back to sign in</router-link>
        </p>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { sendPasswordReset } from '../firebase'

const email = ref('')
const error = ref('')
const sent = ref(false)

async function handleSubmit() {
  error.value = ''
  try {
    await sendPasswordReset(email.value)
    sent.value = true
  } catch {
    error.value = 'Could not send reset email. Please try again.'
  }
}
</script>
