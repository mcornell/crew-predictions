<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="login-title">Sign In</h1>
      <p class="login-sub">Pick your scores. Be wrong in public. It's tradition.</p>
      <form class="login-form" data-testid="login-form" @submit.prevent="handleSubmit">
        <input class="form-input" v-model="email" type="email" placeholder="Email" required />
        <input class="form-input" v-model="password" type="password" placeholder="Password" required />
        <p v-if="error" class="form-error">{{ error }}</p>
        <button class="btn-submit" type="submit">Sign In</button>
      </form>
      <button class="btn-google" data-testid="google-signin" @click="handleGoogle">
        Sign in with Google
      </button>
      <p class="auth-alt">
        Don't have an account?
        <router-link to="/signup">Sign up</router-link>
      </p>
      <p class="auth-alt">
        <router-link to="/reset">Forgot password?</router-link>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { signIn, signInWithGoogle } from '../firebase'

const router = useRouter()
const email = ref('')
const password = ref('')
const error = ref('')

async function postSession(token: string) {
  await fetch('/auth/session', {
    method: 'POST',
    body: new URLSearchParams({ idToken: token }),
  })
}

async function handleSubmit() {
  error.value = ''
  try {
    const token = await signIn(email.value, password.value)
    await postSession(token)
    router.push('/matches')
  } catch {
    error.value = 'Invalid email or password'
  }
}

async function handleGoogle() {
  error.value = ''
  try {
    const token = await signInWithGoogle()
    await postSession(token)
    router.push('/matches')
  } catch {
    error.value = 'Google sign-in failed'
  }
}
</script>
