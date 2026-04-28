import { initializeApp, getApps } from 'firebase/app'
import { getAuth, connectAuthEmulator, signInWithEmailAndPassword, createUserWithEmailAndPassword, signInWithRedirect, getRedirectResult, GoogleAuthProvider, sendPasswordResetEmail } from 'firebase/auth'
import { getAnalytics } from 'firebase/analytics'

declare global {
  interface Window {
    __firebaseConfig?: Record<string, string | undefined>
  }
}

export function initAnalytics() {
  const cfg = window.__firebaseConfig
  if (!cfg?.measurementId || !cfg?.appId || !cfg?.projectId) return
  const app = getApps().length ? getApps()[0] : initializeApp(cfg)
  getAnalytics(app)
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

