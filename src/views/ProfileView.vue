<template>
  <div class="page">
    <div v-if="profile" class="profile-page">
      <div class="profile-header">
        <h1 class="page-title">{{ profile.handle }}</h1>
        <p v-if="profile.location" class="profile-location">{{ profile.location }}</p>
      </div>

      <div class="profile-stats" v-if="profile.predictionCount > 0 || profile.acesRadio.rank > 0">
        <div class="stat-item">
          <span class="stat-label">Predictions</span>
          <span class="stat-value" data-testid="prediction-count">{{ profile.predictionCount }}</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">Aces Radio</span>
          <span class="stat-value" data-testid="aces-radio-points">
            {{ profile.acesRadio.rank > 0 ? `${profile.acesRadio.points} pts · #${profile.acesRadio.rank}` : '—' }}
          </span>
        </div>
        <div class="stat-item">
          <span class="stat-label">Upper 90 Club</span>
          <span class="stat-value">
            {{ profile.upper90Club.rank > 0 ? `${profile.upper90Club.points} pts · #${profile.upper90Club.rank}` : '—' }}
          </span>
        </div>
      </div>

      <div v-if="isOwnProfile" class="login-card" style="margin-top: 2rem">
        <h2 class="login-title" style="font-size: 1.2rem">Edit Profile</h2>
        <form class="login-form" data-testid="profile-form" @submit.prevent="handleSubmit">
          <input
            class="form-input"
            data-testid="display-name-input"
            v-model="displayName"
            type="text"
            placeholder="Display name"
            required
          />
          <input
            class="form-input"
            data-testid="location-input"
            v-model="location"
            type="text"
            placeholder="Location (e.g. Columbus, OH)"
          />
          <p v-if="error" class="form-error">{{ error }}</p>
          <button class="btn-submit" type="submit">Save</button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, inject, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { updateDisplayName } from '../firebase'
import type { Ref } from 'vue'

interface Standing { points: number; rank: number }
interface Profile {
  userID: string
  handle: string
  location: string
  predictionCount: number
  acesRadio: Standing
  upper90Club: Standing
}

const currentUser = inject<Ref<{ userID: string; handle: string; emailVerified: boolean } | null>>('currentUser')
const route = useRoute()
const router = useRouter()

const profile = ref<Profile | null>(null)
const displayName = ref('')
const location = ref('')
const error = ref('')

const isOwnProfile = computed(() =>
  !!currentUser?.value && !!profile.value && currentUser.value.userID === profile.value.userID
)

onMounted(async () => {
  const userID = route.params.userID as string
  if (!userID) {
    router.replace('/login')
    return
  }
  const res = await fetch(`/api/profile/${userID}`)
  if (!res.ok) {
    router.replace('/login')
    return
  }
  profile.value = await res.json()
  displayName.value = profile.value!.handle
  location.value = profile.value!.location
})

async function handleSubmit() {
  error.value = ''
  try {
    await updateDisplayName(displayName.value)
    const body = new URLSearchParams({ handle: displayName.value, location: location.value })
    await fetch('/auth/handle', { method: 'POST', body })
    router.push('/matches')
  } catch {
    error.value = 'Could not save profile. Please try again.'
  }
}
</script>
