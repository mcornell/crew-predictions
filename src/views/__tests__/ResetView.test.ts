import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import ResetView from '../ResetView.vue'

vi.mock('../../firebase', () => ({
  sendPasswordReset: vi.fn(),
}))

function makeRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/reset', component: ResetView },
      { path: '/login', component: { template: '<div />' } },
    ],
  })
}

describe('ResetView', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('renders an email input and submit button', () => {
    const wrapper = mount(ResetView, { global: { plugins: [makeRouter()] } })
    expect(wrapper.find('form[data-testid="reset-form"]').exists()).toBe(true)
    expect(wrapper.find('input[type="email"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('shows confirmation after successful submission', async () => {
    const { sendPasswordReset } = await import('../../firebase')
    vi.mocked(sendPasswordReset).mockResolvedValue()

    const wrapper = mount(ResetView, { global: { plugins: [makeRouter()] } })
    await wrapper.find('input[type="email"]').setValue('fan@example.com')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(sendPasswordReset).toHaveBeenCalledWith('fan@example.com')
    expect(wrapper.find('[data-testid="reset-confirmation"]').exists()).toBe(true)
    expect(wrapper.find('form[data-testid="reset-form"]').exists()).toBe(false)
  })

  it('shows error message on failure', async () => {
    const { sendPasswordReset } = await import('../../firebase')
    vi.mocked(sendPasswordReset).mockRejectedValue(new Error('network error'))

    const wrapper = mount(ResetView, { global: { plugins: [makeRouter()] } })
    await wrapper.find('input[type="email"]').setValue('fan@example.com')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(wrapper.find('.form-error').exists()).toBe(true)
    expect(wrapper.find('[data-testid="reset-confirmation"]').exists()).toBe(false)
  })

  it('renders a link back to the login page', () => {
    const wrapper = mount(ResetView, { global: { plugins: [makeRouter()] } })
    const link = wrapper.find('a[href="/login"]')
    expect(link.exists()).toBe(true)
  })
})
