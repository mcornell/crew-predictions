import { describe, it, expect } from 'vitest'
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

  it('always shows a leaderboard link', () => {
    const wrapper = mount(AppHeader, { props: { user: null } })
    expect(wrapper.find('a[href="/leaderboard"]').exists()).toBe(true)
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
})
