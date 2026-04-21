import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import ResetView from '../ResetView.vue'
import { makeRouter } from '../../test-utils/router'

vi.mock('../../firebase', () => ({
  sendPasswordReset: vi.fn(),
}))

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

  it('email input has autocomplete="email"', () => {
    const wrapper = mount(ResetView, { global: { plugins: [makeRouter()] } })
    expect(wrapper.find('input[type="email"]').attributes('autocomplete')).toBe('email')
  })

  it('renders a link back to the login page', () => {
    const wrapper = mount(ResetView, { global: { plugins: [makeRouter()] } })
    const link = wrapper.find('a[href="/login"]')
    expect(link.exists()).toBe(true)
  })
})
