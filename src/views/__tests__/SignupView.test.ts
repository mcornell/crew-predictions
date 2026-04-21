import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import SignupView from '../SignupView.vue'

vi.mock('../../firebase', () => ({
  signUp: vi.fn(),
  signInWithGoogle: vi.fn(),
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

  it('renders a Sign in with Google button', () => {
    const wrapper = mount(SignupView, { global: { plugins: [makeRouter()] } })
    expect(wrapper.find('button[data-testid="google-signin"]').exists()).toBe(true)
  })

  it('renders a link to the login page', () => {
    const wrapper = mount(SignupView, { global: { plugins: [makeRouter()] } })
    const link = wrapper.find('a[href="/login"]')
    expect(link.exists()).toBe(true)
    expect(link.text()).toBe('Sign in')
  })

  it('calls signInWithGoogle on button click', async () => {
    const { signInWithGoogle } = await import('../../firebase')
    vi.mocked(signInWithGoogle).mockResolvedValue(undefined)

    const router = makeRouter()
    await router.push('/signup')
    const wrapper = mount(SignupView, { global: { plugins: [router] } })

    await wrapper.find('button[data-testid="google-signin"]').trigger('click')
    await flushPromises()

    expect(signInWithGoogle).toHaveBeenCalled()
  })

  it('shows a specific message when email is already registered', async () => {
    const { signUp } = await import('../../firebase')
    const err: any = new Error('email in use')
    err.code = 'auth/email-already-in-use'
    vi.mocked(signUp).mockRejectedValue(err)

    const wrapper = mount(SignupView, { global: { plugins: [makeRouter()] } })
    await wrapper.find('input[type="email"]').setValue('taken@crew.mock')
    await wrapper.find('input[type="password"]').setValue('pass123')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(wrapper.find('.form-error').text()).toBe('That email is already registered. Sign in instead.')
  })

  it('shows a specific message when password is too weak', async () => {
    const { signUp } = await import('../../firebase')
    const err: any = new Error('weak password')
    err.code = 'auth/weak-password'
    vi.mocked(signUp).mockRejectedValue(err)

    const wrapper = mount(SignupView, { global: { plugins: [makeRouter()] } })
    await wrapper.find('input[type="email"]').setValue('new@crew.mock')
    await wrapper.find('input[type="password"]').setValue('abc')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(wrapper.find('.form-error').text()).toBe('Password must be at least 6 characters.')
  })

  it('shows a generic message for unknown sign-up errors', async () => {
    const { signUp } = await import('../../firebase')
    const err: any = new Error('something weird')
    err.code = 'auth/something-weird'
    vi.mocked(signUp).mockRejectedValue(err)

    const wrapper = mount(SignupView, { global: { plugins: [makeRouter()] } })
    await wrapper.find('input[type="email"]').setValue('new@crew.mock')
    await wrapper.find('input[type="password"]').setValue('pass123')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(wrapper.find('.form-error').text()).toBe('Could not create account.')
  })
})
