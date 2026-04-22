import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import ProfileView from '../ProfileView.vue'
import { makeRouter } from '../../test-utils/router'

vi.mock('../../firebase', () => ({
  updateDisplayName: vi.fn(),
}))

describe('ProfileView', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true }))
  })

  it('renders a display name input and submit button', () => {
    const wrapper = mount(ProfileView, { global: { plugins: [makeRouter()] } })
    expect(wrapper.find('form[data-testid="profile-form"]').exists()).toBe(true)
    expect(wrapper.find('input[data-testid="display-name-input"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('calls updateDisplayName, posts to /auth/handle, and navigates to /matches', async () => {
    const { updateDisplayName } = await import('../../firebase')
    vi.mocked(updateDisplayName).mockResolvedValue()

    const r = makeRouter()
    await r.push('/profile')
    const wrapper = mount(ProfileView, { global: { plugins: [r] } })

    await wrapper.find('input[data-testid="display-name-input"]').setValue('Nordecke Regular')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(updateDisplayName).toHaveBeenCalledWith('Nordecke Regular')
    expect(vi.mocked(fetch)).toHaveBeenCalledWith('/auth/handle', expect.objectContaining({ method: 'POST' }))
    expect(r.currentRoute.value.path).toBe('/matches')
  })

  it('shows an error message when updateDisplayName fails', async () => {
    const { updateDisplayName } = await import('../../firebase')
    vi.mocked(updateDisplayName).mockRejectedValue(new Error('network error'))

    const wrapper = mount(ProfileView, { global: { plugins: [makeRouter()] } })
    await wrapper.find('input[data-testid="display-name-input"]').setValue('Nordecke Regular')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(wrapper.find('.form-error').exists()).toBe(true)
  })
})
