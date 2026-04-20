import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import SignupView from '../SignupView.vue'

vi.mock('../../firebase', () => ({
  signUp: vi.fn(),
}))

function makeRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/signup', component: SignupView },
      { path: '/matches', component: { template: '<div />' } },
    ],
  })
}

beforeEach(() => {
  vi.restoreAllMocks()
})

describe('SignupView', () => {
  it('renders an email/password form', () => {
    const wrapper = mount(SignupView, { global: { plugins: [makeRouter()] } })
    expect(wrapper.find('form[data-testid="signup-form"]').exists()).toBe(true)
    expect(wrapper.find('input[type="email"]').exists()).toBe(true)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('calls signUp and navigates to /matches on success', async () => {
    const { signUp } = await import('../../firebase')
    vi.mocked(signUp).mockResolvedValue('fake-token')
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true }))

    const router = makeRouter()
    await router.push('/signup')
    const wrapper = mount(SignupView, { global: { plugins: [router] } })

    await wrapper.find('input[type="email"]').setValue('new@crew.mock')
    await wrapper.find('input[type="password"]').setValue('pass123')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(signUp).toHaveBeenCalledWith('new@crew.mock', 'pass123')
    expect(router.currentRoute.value.path).toBe('/matches')
  })
})
