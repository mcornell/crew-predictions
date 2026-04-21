import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import App from '../App.vue'

vi.mock('../firebase', () => ({}))

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

  it('hides email verification banner when user is verified', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ handle: 'fan@crew.mock', emailVerified: true }),
    }))
    const wrapper = mount(App, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="email-verification-banner"]').exists()).toBe(false)
  })

  it('shows email verification banner when user is not verified', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ handle: 'fan@crew.mock', emailVerified: false }),
    }))
    const wrapper = mount(App, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="email-verification-banner"]').exists()).toBe(true)
  })

  it('re-fetches /api/me after route change to update auth state', async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce({ ok: false, status: 401 })
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ handle: 'testfan@crew.mock' }) })
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
})
