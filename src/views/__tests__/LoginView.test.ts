import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import LoginView from '../LoginView.vue'

vi.mock('../../firebase', () => ({
  signIn: vi.fn(),
  signInWithGoogle: vi.fn(),
}))

const router = createRouter({
  history: createMemoryHistory(),
  routes: [
    { path: '/login', component: LoginView },
    { path: '/matches', component: { template: '<div />' } },
  ],
})

describe('LoginView', () => {
  it('renders an email/password form', () => {
    const wrapper = mount(LoginView, { global: { plugins: [router] } })
    expect(wrapper.find('form[data-testid="login-form"]').exists()).toBe(true)
    expect(wrapper.find('input[type="email"]').exists()).toBe(true)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('calls signIn and navigates to /matches on success', async () => {
    const { signIn } = await import('../../firebase')
    vi.mocked(signIn).mockResolvedValue('fake-token')
    global.fetch = vi.fn().mockResolvedValue({ ok: true })

    await router.push('/login')
    const wrapper = mount(LoginView, { global: { plugins: [router] } })

    await wrapper.find('input[type="email"]').setValue('test@crew.mock')
    await wrapper.find('input[type="password"]').setValue('pass123')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(signIn).toHaveBeenCalledWith('test@crew.mock', 'pass123')
    expect(router.currentRoute.value.path).toBe('/matches')
  })

  it('email input has autocomplete="email"', () => {
    const wrapper = mount(LoginView, { global: { plugins: [router] } })
    expect(wrapper.find('input[type="email"]').attributes('autocomplete')).toBe('email')
  })

  it('password input has autocomplete="current-password"', () => {
    const wrapper = mount(LoginView, { global: { plugins: [router] } })
    expect(wrapper.find('input[type="password"]').attributes('autocomplete')).toBe('current-password')
  })

  it('renders a Sign in with Google button', () => {
    const wrapper = mount(LoginView, { global: { plugins: [router] } })
    expect(wrapper.find('button[data-testid="google-signin"]').exists()).toBe(true)
  })

  it('renders a link to the sign-up page', () => {
    const wrapper = mount(LoginView, { global: { plugins: [router] } })
    const link = wrapper.find('a[href="/signup"]')
    expect(link.exists()).toBe(true)
    expect(link.text()).toBe('Sign up')
  })

  it('renders a "Forgot password?" link to /reset', () => {
    const wrapper = mount(LoginView, { global: { plugins: [router] } })
    const link = wrapper.find('a[href="/reset"]')
    expect(link.exists()).toBe(true)
    expect(link.text()).toBe('Forgot password?')
  })

  it('calls signInWithGoogle on button click', async () => {
    const { signInWithGoogle } = await import('../../firebase')
    vi.mocked(signInWithGoogle).mockResolvedValue(undefined)

    await router.push('/login')
    const wrapper = mount(LoginView, { global: { plugins: [router] } })

    await wrapper.find('button[data-testid="google-signin"]').trigger('click')
    await flushPromises()

    expect(signInWithGoogle).toHaveBeenCalled()
  })
})
