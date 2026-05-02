import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import App from '../App.vue'
import AppHeader from '../components/AppHeader.vue'

vi.mock('../firebase', () => ({
  getGoogleRedirectResult: vi.fn().mockResolvedValue(null),
}))

function makeRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/', component: { template: '<div />' } },
      { path: '/matches', component: { template: '<div />' } },
    ],
  })
}

beforeEach(() => {
  vi.restoreAllMocks()
})

describe('App', () => {
  it('shows Sign In link when /api/me returns 401', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false, status: 401 }))
    const wrapper = mount(App, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    expect(wrapper.find('a[href="/login"]').exists()).toBe(true)
  })

  it('shows handle in header when /api/me returns user', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ handle: 'BlackAndGold@bsky.mock' }),
    }))
    const wrapper = mount(App, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    expect(wrapper.text()).toContain('BlackAndGold@bsky.mock')
    expect(wrapper.find('a[href="/login"]').exists()).toBe(false)
  })

  it('does not show an email verification banner', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ handle: 'fan@crew.mock', emailVerified: false }),
    }))
    const wrapper = mount(App, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="email-verification-banner"]').exists()).toBe(false)
  })

  it('still fetches /api/me when getGoogleRedirectResult throws', async () => {
    const { getGoogleRedirectResult } = await import('../firebase')
    vi.mocked(getGoogleRedirectResult).mockRejectedValueOnce(new Error('auth/popup-closed-by-user'))
    const fetchMock = vi.fn().mockResolvedValue({ ok: false, status: 401 })
    vi.stubGlobal('fetch', fetchMock)

    mount(App, { global: { plugins: [makeRouter()] } })
    await flushPromises()

    expect(fetchMock).toHaveBeenCalledWith('/api/me')
  })

  it('completes google redirect sign-in and navigates to /matches', async () => {
    const { getGoogleRedirectResult } = await import('../firebase')
    vi.mocked(getGoogleRedirectResult).mockResolvedValueOnce('google-id-token')
    const fetchMock = vi.fn()
      .mockResolvedValueOnce({ ok: true })                                                                   // /auth/session
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ handle: 'NewUser@bsky.mock' }) })    // /api/me (fetchUser)
      .mockResolvedValueOnce({ ok: false })                                                                  // /api/seasons
    vi.stubGlobal('fetch', fetchMock)

    const router = makeRouter()
    mount(App, { global: { plugins: [router] } })
    await flushPromises()

    expect(fetchMock).toHaveBeenCalledWith('/auth/session', expect.objectContaining({ method: 'POST' }))
    expect(router.currentRoute.value.path).toBe('/matches')
  })

  it('falls through to fetchUser when google redirect session create fails', async () => {
    const { getGoogleRedirectResult } = await import('../firebase')
    vi.mocked(getGoogleRedirectResult).mockResolvedValueOnce('google-id-token')
    const fetchMock = vi.fn()
      .mockResolvedValueOnce({ ok: false })                                                                  // /auth/session (fails)
      .mockResolvedValueOnce({ ok: false, status: 401 })                                                    // /api/me
      .mockResolvedValueOnce({ ok: false })                                                                  // /api/seasons
    vi.stubGlobal('fetch', fetchMock)

    mount(App, { global: { plugins: [makeRouter()] } })
    await flushPromises()

    expect(fetchMock).toHaveBeenCalledWith('/api/me')
  })

  it('re-fetches /api/me after route change to update auth state', async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce({ ok: false, status: 401 })                                                   // /api/me (initial)
      .mockResolvedValueOnce({ ok: false })                                                                 // /api/seasons
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ handle: 'testfan@crew.mock' }) })   // /api/me (route change)
    vi.stubGlobal('fetch', fetchMock)

    const router = makeRouter()
    await router.push('/')
    const wrapper = mount(App, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.find('a[href="/login"]').exists()).toBe(true)

    await router.push('/matches')
    await flushPromises()
    expect(wrapper.text()).toContain('testfan@crew.mock')
  })

  it.each([
    {
      desc: 'returns past seasons before the current one in descending order',
      input: [
        { id: '2026', name: '2026 Season', isCurrent: false },
        { id: '2027-sprint', name: '2027 Sprint Season', isCurrent: true },
        { id: '2027-28', name: '2027-28 Season', isCurrent: false },
      ],
      expected: [
        { id: '2026', name: '2026 Season', isCurrent: false },
      ],
    },
    {
      desc: 'returns no past seasons when the current season is first',
      input: [
        { id: '2026', name: '2026 Season', isCurrent: true },
        { id: '2027-sprint', name: '2027 Sprint Season', isCurrent: false },
      ],
      expected: [],
    },
    {
      desc: 'returns multiple past seasons in descending order',
      input: [
        { id: '2026', name: '2026 Season', isCurrent: false },
        { id: '2027-sprint', name: '2027 Sprint Season', isCurrent: false },
        { id: '2027-28', name: '2027-28 Season', isCurrent: true },
        { id: '2028-29', name: '2028-29 Season', isCurrent: false },
      ],
      expected: [
        { id: '2027-sprint', name: '2027 Sprint Season', isCurrent: false },
        { id: '2026', name: '2026 Season', isCurrent: false },
      ],
    },
  ])('passes seasons to AppHeader: $desc', async ({ input, expected }) => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce({ ok: false, status: 401 })
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ seasons: input }) })
    vi.stubGlobal('fetch', fetchMock)

    const wrapper = mount(App, { global: { plugins: [makeRouter()] } })
    await flushPromises()

    const header = wrapper.findComponent(AppHeader)
    expect(header.props('seasons')).toEqual(expected)
  })
})
