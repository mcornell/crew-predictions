export const SMOKE_NEW_EMAILS = [
  'smoke-new-desktop@crew-predictions-staging.web.app',
  'smoke-new-ios@crew-predictions-staging.web.app',
]

export async function deleteSmokeNewAccounts() {
  const apiKey = process.env.STAGING_API_KEY!
  const password = process.env.SMOKE_TEST_PASSWORD!

  await Promise.all(SMOKE_NEW_EMAILS.map(async (email) => {
    const signIn = await fetch(
      `https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=${apiKey}`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password, returnSecureToken: true }),
      }
    )
    if (!signIn.ok) return
    const { idToken } = await signIn.json()
    await fetch(
      `https://identitytoolkit.googleapis.com/v1/accounts:delete?key=${apiKey}`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ idToken }),
      }
    )
  }))
}
