<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="login-title">Sign Up</h1>
      <p class="login-sub">Pick your scores. Be wrong in public. It's tradition.</p>
      <form class="login-form" data-testid="signup-form" @submit.prevent="handleSubmit">
        <input class="form-input" v-model="email" type="email" placeholder="Email" required />
        <input class="form-input" v-model="password" type="password" placeholder="Password" required />
        <p v-if="error" class="form-error">{{ error }}</p>
        <button class="btn-submit" type="submit">Sign Up</button>
      </form>
      <button class="btn-google" data-testid="google-signin" @click="handleGoogle">
        Sign in with Google
      </button>
      <p class="auth-alt">
        Already have an account?
        <router-link to="/login">Sign in</router-link>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { signUp, signInWithGoogle } from '../firebase'

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

const signUpErrorMessages: Record<string, string> = {
  'auth/email-already-in-use': 'That email is already registered. Sign in instead.',
  'auth/weak-password': 'Password must be at least 6 characters.',
}

async function handleSubmit() {
  error.value = ''
  try {
    const token = await signUp(email.value, password.value)
    await postSession(token)
    router.push('/matches')
  } catch (e: any) {
    error.value = signUpErrorMessages[e?.code] ?? 'Could not create account.'
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
