import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import AppHeader from '../AppHeader.vue'

describe('AppHeader', () => {
  it('renders a banner with "Crew Predictions"', () => {
    const wrapper = mount(AppHeader, { props: { user: null } })
    expect(wrapper.find('header').exists()).toBe(true)
    expect(wrapper.text()).toContain('Crew Predictions')
  })

  it('shows Sign In link when no user', () => {
    const wrapper = mount(AppHeader, { props: { user: null } })
    const link = wrapper.find('a[href="/login"]')
    expect(link.exists()).toBe(true)
    expect(link.text()).toBe('Sign In')
  })

  it('shows username and Sign out link when logged in', () => {
    const wrapper = mount(AppHeader, { props: { user: { userID: 'firebase:abc', handle: 'BlackAndGold@bsky.mock' } } })
    expect(wrapper.text()).toContain('BlackAndGold@bsky.mock')
    expect(wrapper.find('a[href="/auth/logout"]').text()).toBe('Sign out')
  })

  it('hides Sign In link when logged in', () => {
    const wrapper = mount(AppHeader, { props: { user: { userID: 'firebase:abc', handle: 'BlackAndGold@bsky.mock' } } })
    expect(wrapper.find('a[href="/login"]').exists()).toBe(false)
  })

  it('shows a profile link when logged in', () => {
    const wrapper = mount(AppHeader, { props: { user: { userID: 'firebase:abc', handle: 'BlackAndGold@bsky.mock' } } })
    expect(wrapper.find('a[href="/profile/firebase:abc"]').exists()).toBe(true)
  })

  it('hides profile link when not logged in', () => {
    const wrapper = mount(AppHeader, { props: { user: null } })
    expect(wrapper.find('a[href="/profile"]').exists()).toBe(false)
  })

  it('always shows a leaderboard nav trigger', () => {
    const wrapper = mount(AppHeader, { props: { user: null } })
    expect(wrapper.find('[data-testid="season-selector"]').exists()).toBe(true)
  })

  it('clicking leaderboard trigger opens season flyout', async () => {
    const seasons = [{ id: '2026', name: '2026 Season', isCurrent: true }]
    const wrapper = mount(AppHeader, { props: { user: null, seasons } })
    await wrapper.find('[data-testid="season-selector"]').trigger('click')
    expect(wrapper.find('[data-testid="season-flyout"]').exists()).toBe(true)
  })

  it('season flyout contains a current season link and historical season links', async () => {
    const seasons = [{ id: '2026', name: '2026 Season', isCurrent: false }]
    const wrapper = mount(AppHeader, { props: { user: null, seasons } })
    await wrapper.find('[data-testid="season-selector"]').trigger('click')
    const flyout = wrapper.find('[data-testid="season-flyout"]')
    expect(flyout.find('a[href="/leaderboard"]').exists()).toBe(true)
    expect(flyout.text()).toContain('2026 Season')
  })

  it('clicking a flyout season link closes the flyout', async () => {
    const seasons = [{ id: '2026', name: '2026 Season', isCurrent: false }]
    const wrapper = mount(AppHeader, { props: { user: null, seasons } })
    await wrapper.find('[data-testid="season-selector"]').trigger('click')
    await wrapper.find('[data-testid="season-flyout"] a').trigger('click')
    expect(wrapper.find('[data-testid="season-flyout"]').exists()).toBe(false)
  })

  it('mobile drawer shows seasons under leaderboard when expanded', async () => {
    const seasons = [{ id: '2026', name: '2026 Season', isCurrent: false }]
    const wrapper = mount(AppHeader, { props: { user: null, seasons } })
    await wrapper.find('button[data-testid="hamburger"]').trigger('click')
    const drawer = wrapper.find('[data-testid="mobile-drawer"]')
    await drawer.find('[data-testid="drawer-lb-toggle"]').trigger('click')
    expect(drawer.text()).toContain('2026 Season')
  })

  it('hides Sign In when auth is loading', () => {
    const wrapper = mount(AppHeader, { props: { user: null, loading: true } })
    expect(wrapper.find('a[href="/login"]').exists()).toBe(false)
  })

  it('renders a hamburger button', () => {
    const wrapper = mount(AppHeader, { props: { user: null } })
    expect(wrapper.find('button[data-testid="hamburger"]').exists()).toBe(true)
  })

  it('drawer is not visible by default', () => {
    const wrapper = mount(AppHeader, { props: { user: null } })
    expect(wrapper.find('[data-testid="mobile-drawer"]').exists()).toBe(false)
  })

  it('clicking hamburger opens the drawer', async () => {
    const wrapper = mount(AppHeader, { props: { user: null } })
    await wrapper.find('button[data-testid="hamburger"]').trigger('click')
    expect(wrapper.find('[data-testid="mobile-drawer"]').exists()).toBe(true)
  })

  it('clicking a drawer link closes the drawer', async () => {
    const wrapper = mount(AppHeader, { props: { user: null } })
    await wrapper.find('button[data-testid="hamburger"]').trigger('click')
    await wrapper.find('[data-testid="mobile-drawer"] a').trigger('click')
    expect(wrapper.find('[data-testid="mobile-drawer"]').exists()).toBe(false)
  })

  it('drawer shows user profile link and sign out when logged in', async () => {
    const wrapper = mount(AppHeader, { props: { user: { userID: 'firebase:abc', handle: 'CrewFan' } } })
    await wrapper.find('button[data-testid="hamburger"]').trigger('click')
    const drawer = wrapper.find('[data-testid="mobile-drawer"]')
    expect(drawer.find('a[href="/profile/firebase:abc"]').exists()).toBe(true)
    expect(drawer.find('a[href="/auth/logout"]').exists()).toBe(true)
  })

  it('drawer shows Sign In when no user and not loading', async () => {
    const wrapper = mount(AppHeader, { props: { user: null } })
    await wrapper.find('button[data-testid="hamburger"]').trigger('click')
    expect(wrapper.find('[data-testid="mobile-drawer"] a[href="/login"]').exists()).toBe(true)
  })

  it('drawer hides Sign In when auth is loading', async () => {
    const wrapper = mount(AppHeader, { props: { user: null, loading: true } })
    await wrapper.find('button[data-testid="hamburger"]').trigger('click')
    expect(wrapper.find('[data-testid="mobile-drawer"] a[href="/login"]').exists()).toBe(false)
  })

  it('clicking backdrop closes the drawer', async () => {
    const wrapper = mount(AppHeader, { props: { user: null } })
    await wrapper.find('button[data-testid="hamburger"]').trigger('click')
    expect(wrapper.find('[data-testid="mobile-drawer"]').exists()).toBe(true)
    await wrapper.find('.drawer-backdrop').trigger('click')
    expect(wrapper.find('[data-testid="mobile-drawer"]').exists()).toBe(false)
  })

  it('removes resize listener on unmount', () => {
    const removeSpy = vi.spyOn(window, 'removeEventListener')
    const wrapper = mount(AppHeader, { props: { user: null } })
    wrapper.unmount()
    expect(removeSpy).toHaveBeenCalledWith('resize', expect.any(Function))
    removeSpy.mockRestore()
  })
})
