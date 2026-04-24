import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { readGuestPredictions, writeGuestPredictions, clearGuestPredictions, flushGuestPredictions } from '../guestPredictions'

beforeEach(() => {
  localStorage.clear()
})

afterEach(() => {
  vi.restoreAllMocks()
  localStorage.clear()
})

describe('readGuestPredictions', () => {
  it('returns empty object when nothing stored', () => {
    expect(readGuestPredictions()).toEqual({})
  })

  it('returns stored predictions', () => {
    localStorage.setItem('guestPredictions', JSON.stringify({ 'm1': { homeGoals: 2, awayGoals: 1 } }))
    expect(readGuestPredictions()).toEqual({ 'm1': { homeGoals: 2, awayGoals: 1 } })
  })

  it('returns empty object on malformed JSON', () => {
    localStorage.setItem('guestPredictions', 'not-json')
    expect(readGuestPredictions()).toEqual({})
  })
})

describe('writeGuestPredictions', () => {
  it('writes predictions to localStorage', () => {
    writeGuestPredictions({ 'm1': { homeGoals: 1, awayGoals: 0 } })
    expect(JSON.parse(localStorage.getItem('guestPredictions')!)).toEqual({ 'm1': { homeGoals: 1, awayGoals: 0 } })
  })
})

describe('clearGuestPredictions', () => {
  it('removes the key from localStorage', () => {
    localStorage.setItem('guestPredictions', '{}')
    clearGuestPredictions()
    expect(localStorage.getItem('guestPredictions')).toBeNull()
  })
})

describe('flushGuestPredictions', () => {
  it('does nothing when no guest predictions stored', async () => {
    const fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)
    await flushGuestPredictions()
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('POSTs each prediction to /api/predictions', async () => {
    localStorage.setItem('guestPredictions', JSON.stringify({
      'm1': { homeGoals: 2, awayGoals: 1 },
      'm2': { homeGoals: 0, awayGoals: 0 },
    }))
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true }))
    await flushGuestPredictions()
    const calls = vi.mocked(fetch).mock.calls
    expect(calls).toHaveLength(2)
    const urls = calls.map(c => c[0])
    expect(urls.every(u => u === '/api/predictions')).toBe(true)
  })

  it('clears localStorage after flushing', async () => {
    localStorage.setItem('guestPredictions', JSON.stringify({ 'm1': { homeGoals: 1, awayGoals: 0 } }))
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true }))
    await flushGuestPredictions()
    expect(readGuestPredictions()).toEqual({})
  })

  it('clears prediction even when server rejects (match kicked off)', async () => {
    localStorage.setItem('guestPredictions', JSON.stringify({ 'm1': { homeGoals: 1, awayGoals: 0 } }))
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false, status: 403 }))
    await flushGuestPredictions()
    expect(readGuestPredictions()).toEqual({})
  })

  it('leaves prediction in localStorage on network error', async () => {
    localStorage.setItem('guestPredictions', JSON.stringify({ 'm1': { homeGoals: 1, awayGoals: 0 } }))
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('network')))
    await flushGuestPredictions()
    expect(readGuestPredictions()).toEqual({ 'm1': { homeGoals: 1, awayGoals: 0 } })
  })
})
