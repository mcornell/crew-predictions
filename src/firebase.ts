import { initializeApp, getApps } from 'firebase/app'
import { getAuth, connectAuthEmulator, signInWithEmailAndPassword } from 'firebase/auth'

declare global {
  interface Window {
    __firebaseConfig?: Record<string, string>
  }
}

let emulatorConnected = false

export function getFirebaseAuth() {
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
