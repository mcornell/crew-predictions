import { describe, it, expect } from 'vitest'
import { formatCountdown } from '../countdown'

describe('formatCountdown', () => {
  it('returns "kicks off now" for zero or negative ms', () => {
    expect(formatCountdown(0)).toBe('kicks off now')
    expect(formatCountdown(-1000)).toBe('kicks off now')
  })

  it('formats minutes only when under an hour', () => {
    expect(formatCountdown(45 * 60 * 1000)).toBe('locks in 45m')
    expect(formatCountdown(1 * 60 * 1000)).toBe('locks in 1m')
  })

  it('formats hours and minutes when under a day', () => {
    expect(formatCountdown(2 * 3600 * 1000 + 34 * 60 * 1000)).toBe('locks in 2h 34m')
    expect(formatCountdown(23 * 3600 * 1000)).toBe('locks in 23h 0m')
  })

  it('formats days and hours when over a day', () => {
    expect(formatCountdown(3 * 86400 * 1000 + 12 * 3600 * 1000)).toBe('locks in 3d 12h')
    expect(formatCountdown(1 * 86400 * 1000)).toBe('locks in 1d 0h')
  })
})
