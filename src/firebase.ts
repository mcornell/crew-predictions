import { initializeApp, getApps } from 'firebase/app'
import { getAuth, connectAuthEmulator, signInWithEmailAndPassword, createUserWithEmailAndPassword, signInWithRedirect, getRedirectResult, GoogleAuthProvider, sendPasswordResetEmail, updateProfile, onAuthStateChanged, type User } from 'firebase/auth'

declare global {
  interface Window {
    __firebaseConfig?: Record<string, string>
  }
}

let emulatorConnected = false

function getFirebaseAuth() {
  const app = getApps().length ? getApps()[0] : initializeApp(window.__firebaseConfig ?? {})
  const auth = getAuth(app)
  if (!emulatorConnected && window.__firebaseConfig?.authEmulatorHost) {
    connectAuthEmulator(auth, `http://${window.__firebaseConfig.authEmulatorHost}`, { disableWarnings: true })
    emulatorConnected = true
  }
  return auth
}

export async function signIn(email: string, password: string): Promise<string> {
  const auth = getFirebaseAuth()
  const result = await signInWithEmailAndPassword(auth, email, password)
  return result.user.getIdToken()
}

export async function signUp(email: string, password: string): Promise<string> {
  const auth = getFirebaseAuth()
  const result = await createUserWithEmailAndPassword(auth, email, password)
  return result.user.getIdToken()
}

export async function signInWithGoogle(): Promise<void> {
  const auth = getFirebaseAuth()
  await signInWithRedirect(auth, new GoogleAuthProvider())
}

export async function getGoogleRedirectResult(): Promise<string | null> {
  const auth = getFirebaseAuth()
  const result = await getRedirectResult(auth)
  if (!result) return null
  return result.user.getIdToken()
}

export async function sendPasswordReset(email: string): Promise<void> {
  const auth = getFirebaseAuth()
  await sendPasswordResetEmail(auth, email)
}

function waitForCurrentUser(): Promise<User> {
  return new Promise((resolve, reject) => {
    const auth = getFirebaseAuth()
    const unsubscribe = onAuthStateChanged(auth, (user) => {
      unsubscribe()
      if (user) resolve(user)
      else reject(new Error('not signed in'))
    })
  })
}

export async function updateDisplayName(name: string): Promise<void> {
  const user = await waitForCurrentUser()
  await updateProfile(user, { displayName: name })
  // Force-refresh token so the server reads the updated name claim
  await fetch('/auth/session', {
    method: 'POST',
    body: new URLSearchParams({ idToken: await user.getIdToken(true) }),
  })
}
