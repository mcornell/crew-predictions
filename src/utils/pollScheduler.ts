export const POLL_INTERVAL_MS = 60_000
export const PRE_KICKOFF_MS = 30 * 60_000
export const POST_KICKOFF_WINDOW_MS = 2 * 60 * 60_000

interface Match {
  kickoff: string
  state?: string
  status: string
}

export function isInActiveWindow(matches: Match[], now: number): boolean {
  return matches.some(m => {
    if ((m.state ?? '') === 'in' || m.status === 'STATUS_DELAYED') return true
    if ((m.state ?? '') === 'post') return false
    const kickoff = new Date(m.kickoff).getTime()
    return kickoff - now <= PRE_KICKOFF_MS && now - kickoff <= POST_KICKOFF_WINDOW_MS
  })
}

export function msUntilActiveWindow(matches: Match[], now: number): number | null {
  const times = matches
    .filter(m => (m.state ?? '') !== 'post')
    .map(m => new Date(m.kickoff).getTime() - PRE_KICKOFF_MS)
    .filter(t => t > now)
    .sort((a, b) => a - b)
  return times.length > 0 ? times[0] - now : null
}
