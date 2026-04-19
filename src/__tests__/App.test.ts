import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import App from '../App.vue'

const router = createRouter({
  history: createMemoryHistory(),
  routes: [{ path: '/', component: { template: '<div />' } }],
})

beforeEach(() => {
  vi.restoreAllMocks()
})

describe('App', () => {
  it('shows Sign In link when /api/me returns 401', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false, status: 401 }))
    const wrapper = mount(App, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.find('a[href="/login"]').exists()).toBe(true)
  })

  it('shows handle in header when /api/me returns user', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ handle: 'BlackAndGold@bsky.mock' }),
    }))
    const wrapper = mount(App, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.text()).toContain('BlackAndGold@bsky.mock')
    expect(wrapper.find('a[href="/login"]').exists()).toBe(false)
  })
})
