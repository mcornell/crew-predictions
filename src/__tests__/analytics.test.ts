import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

const mockGetAnalytics = vi.fn()
vi.mock('firebase/analytics', () => ({ getAnalytics: mockGetAnalytics }))
vi.mock('firebase/app', () => ({
  getApps: () => [],
  initializeApp: vi.fn(() => ({})),
}))

describe('initAnalytics', () => {
  beforeEach(() => {
    mockGetAnalytics.mockClear()
  })

  afterEach(() => {
    vi.resetModules()
    delete (window as any).__firebaseConfig
  })

  it('initializes analytics when measurementId is present', async () => {
    ;(window as any).__firebaseConfig = { measurementId: 'G-TEST123', appId: '1:123:web:abc' }
    const { initAnalytics } = await import('../firebase')
    initAnalytics()
    expect(mockGetAnalytics).toHaveBeenCalledOnce()
  })

  it('does not initialize analytics when measurementId is absent', async () => {
    ;(window as any).__firebaseConfig = { projectId: 'crew-predictions' }
    const { initAnalytics } = await import('../firebase')
    initAnalytics()
    expect(mockGetAnalytics).not.toHaveBeenCalled()
  })
})
