const STAGING_API_KEY = process.env.STAGING_API_KEY!
const SMOKE_PASSWORD = process.env.SMOKE_TEST_PASSWORD!
const SMOKE_NEW_EMAIL = 'smoke-new@crew-predictions-staging.web.app'

export default async function globalTeardown() {
  const signIn = await fetch(
    `https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=${STAGING_API_KEY}`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: SMOKE_NEW_EMAIL, password: SMOKE_PASSWORD, returnSecureToken: true }),
    }
  )
  if (!signIn.ok) return
  const { idToken } = await signIn.json()
  await fetch(
    `https://identitytoolkit.googleapis.com/v1/accounts:delete?key=${STAGING_API_KEY}`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ idToken }),
    }
  )
}
