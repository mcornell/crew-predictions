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
})
