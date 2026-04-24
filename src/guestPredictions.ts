const GUEST_KEY = 'guestPredictions'

export interface GuestPrediction {
  homeGoals: number
  awayGoals: number
}

export function readGuestPredictions(): Record<string, GuestPrediction> {
  try {
    return JSON.parse(localStorage.getItem(GUEST_KEY) ?? '{}')
  } catch {
    return {}
  }
}

export function writeGuestPredictions(data: Record<string, GuestPrediction>): void {
  try {
    localStorage.setItem(GUEST_KEY, JSON.stringify(data))
  } catch {
    // Safari Private Browsing blocks localStorage writes; degrade silently
  }
}

export function clearGuestPredictions(): void {
  try {
    localStorage.removeItem(GUEST_KEY)
  } catch {}
}

export async function flushGuestPredictions(): Promise<void> {
  const stored = readGuestPredictions()
  const matchIDs = Object.keys(stored)
  if (matchIDs.length === 0) return

  for (const matchID of matchIDs) {
    const { homeGoals, awayGoals } = stored[matchID]
    const body = new URLSearchParams({
      match_id: matchID,
      home_goals: String(homeGoals),
      away_goals: String(awayGoals),
    })
    let networkError = false
    try {
      await fetch('/api/predictions', { method: 'POST', body })
    } catch {
      networkError = true
    }
    // On network error, leave in localStorage so it retries next login.
    // On any server response (accept or reject), remove — if the server rejected
    // (e.g. match kicked off), there's nothing the user can do.
    if (!networkError) {
      delete stored[matchID]
    }
  }

  writeGuestPredictions(stored)
}
