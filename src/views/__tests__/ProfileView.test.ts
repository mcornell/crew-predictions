import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import ProfileView from '../ProfileView.vue'

vi.mock('../../firebase', () => ({
  updateDisplayName: vi.fn(),
}))

function makeRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/profile', component: ProfileView },
      { path: '/matches', component: { template: '<div />' } },
    ],
  })
}

describe('ProfileView', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('renders a display name input and submit button', () => {
    const wrapper = mount(ProfileView, { global: { plugins: [makeRouter()] } })
    expect(wrapper.find('form[data-testid="profile-form"]').exists()).toBe(true)
    expect(wrapper.find('input[data-testid="display-name-input"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('calls updateDisplayName and navigates to /matches on success', async () => {
    const { updateDisplayName } = await import('../../firebase')
    vi.mocked(updateDisplayName).mockResolvedValue()

    const r = makeRouter()
    await r.push('/profile')
    const wrapper = mount(ProfileView, { global: { plugins: [r] } })

    await wrapper.find('input[data-testid="display-name-input"]').setValue('Nordecke Regular')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(updateDisplayName).toHaveBeenCalledWith('Nordecke Regular')
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
