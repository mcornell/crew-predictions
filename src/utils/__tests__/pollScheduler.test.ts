import { describe, it, expect } from 'vitest'
import { isInActiveWindow, msUntilActiveWindow, POLL_INTERVAL_MS, PRE_KICKOFF_MS } from '../pollScheduler'

const now = new Date('2026-05-01T19:00:00Z').getTime()

function match(overrides: Partial<{ kickoff: string; state: string; status: string }>) {
  return { kickoff: new Date(now + 2 * 3600_000).toISOString(), state: 'pre', status: 'STATUS_SCHEDULED', ...overrides }
}

describe('isInActiveWindow', () => {
  it('returns true when a match is live', () => {
    expect(isInActiveWindow([match({ state: 'in' })], now)).toBe(true)
  })

  it('returns true when a match is delayed', () => {
    expect(isInActiveWindow([match({ status: 'STATUS_DELAYED' })], now)).toBe(true)
  })

  it('returns true when kickoff is within 30 minutes', () => {
    const kickoff = new Date(now + 29 * 60_000).toISOString()
    expect(isInActiveWindow([match({ kickoff })], now)).toBe(true)
  })

  it('returns true when kickoff was up to 2 hours ago and match not post', () => {
    const kickoff = new Date(now - 90 * 60_000).toISOString()
    expect(isInActiveWindow([match({ kickoff, state: 'pre' })], now)).toBe(true)
  })

  it('returns false when kickoff was over 2 hours ago and state is still pre', () => {
    const kickoff = new Date(now - 121 * 60_000).toISOString()
    expect(isInActiveWindow([match({ kickoff, state: 'pre' })], now)).toBe(false)
  })

  it('returns false when match is post', () => {
    const kickoff = new Date(now - 90 * 60_000).toISOString()
    expect(isInActiveWindow([match({ kickoff, state: 'post' })], now)).toBe(false)
  })

  it('returns false when kickoff is more than 30 minutes away', () => {
    const kickoff = new Date(now + 31 * 60_000).toISOString()
    expect(isInActiveWindow([match({ kickoff })], now)).toBe(false)
  })

  it('returns false for empty match list', () => {
    expect(isInActiveWindow([], now)).toBe(false)
  })

  it('returns true if any match is in active window', () => {
    const live = match({ state: 'in' })
    const distant = match({ kickoff: new Date(now + 5 * 3600_000).toISOString() })
    expect(isInActiveWindow([distant, live], now)).toBe(true)
  })
})

describe('msUntilActiveWindow', () => {
  it('returns ms until 30 minutes before next pre kickoff', () => {
    const kickoff = new Date(now + 2 * 3600_000).toISOString()
    const result = msUntilActiveWindow([match({ kickoff })], now)
    expect(result).toBe(2 * 3600_000 - PRE_KICKOFF_MS)
  })

  it('returns null when no upcoming pre matches', () => {
    expect(msUntilActiveWindow([], now)).toBeNull()
    expect(msUntilActiveWindow([match({ state: 'post' })], now)).toBeNull()
  })

  it('picks the nearest kickoff when multiple upcoming', () => {
    const near = match({ kickoff: new Date(now + 2 * 3600_000).toISOString() })
    const far = match({ kickoff: new Date(now + 5 * 3600_000).toISOString() })
    const result = msUntilActiveWindow([far, near], now)
    expect(result).toBe(2 * 3600_000 - PRE_KICKOFF_MS)
  })
})

describe('POLL_INTERVAL_MS', () => {
  it('is 60 seconds', () => {
    expect(POLL_INTERVAL_MS).toBe(60_000)
  })
})
