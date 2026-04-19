<template>
  <form data-testid="login-form" @submit.prevent="handleSubmit">
    <input v-model="email" type="email" placeholder="Email" required />
    <input v-model="password" type="password" placeholder="Password" required />
    <p v-if="error">{{ error }}</p>
    <button type="submit">Sign In</button>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { signIn } from '../firebase'

const router = useRouter()
const email = ref('')
const password = ref('')
const error = ref('')

async function handleSubmit() {
  error.value = ''
  try {
    const token = await signIn(email.value, password.value)
    await fetch('/auth/session', {
      method: 'POST',
      body: new URLSearchParams({ idToken: token }),
    })
    router.push('/matches')
  } catch (e: any) {
    error.value = 'Invalid email or password'
  }
}
</script>
