import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
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
    const wrapper = mount(AppHeader, { props: { user: { handle: 'BlackAndGold@bsky.mock' } } })
    expect(wrapper.text()).toContain('BlackAndGold@bsky.mock')
    expect(wrapper.find('a[href="/auth/logout"]').text()).toBe('Sign out')
  })

  it('hides Sign In link when logged in', () => {
    const wrapper = mount(AppHeader, { props: { user: { handle: 'BlackAndGold@bsky.mock' } } })
    expect(wrapper.find('a[href="/login"]').exists()).toBe(false)
  })
})
