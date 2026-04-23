<template>
  <div class="page">
    <div v-if="profile" class="profile-page">

      <div class="profile-header">
        <h1 class="profile-handle">{{ profile.handle }}</h1>
        <p v-if="profile.location" class="profile-location">
          <span class="profile-location-icon">▸</span>{{ profile.location }}
        </p>
      </div>

      <div class="profile-stats">
        <div class="profile-stat">
          <span class="profile-stat-value" data-testid="prediction-count">{{ profile.predictionCount }}</span>
          <span class="profile-stat-label">Predictions</span>
        </div>
        <div class="profile-stat">
          <span class="profile-stat-value" data-testid="aces-radio-points">{{ profile.acesRadio.rank > 0 ? profile.acesRadio.points : '—' }}</span>
          <span class="profile-stat-label">Aces Radio</span>
          <span v-if="profile.acesRadio.rank > 0" class="profile-stat-rank">#{{ profile.acesRadio.rank }}</span>
        </div>
        <div class="profile-stat">
          <span class="profile-stat-value">{{ profile.upper90Club.rank > 0 ? profile.upper90Club.points : '—' }}</span>
          <span class="profile-stat-label">Upper 90 Club</span>
          <span v-if="profile.upper90Club.rank > 0" class="profile-stat-rank">#{{ profile.upper90Club.rank }}</span>
        </div>
      </div>

      <div v-if="isOwnProfile" class="profile-edit">
        <h2 class="profile-edit-title">Edit Profile</h2>
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

const currentUser = inject<Ref<{ userID: string; handle: string; emailVerified: boolean } | null>>('currentUser', ref(null))
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
