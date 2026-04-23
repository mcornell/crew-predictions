import { describe, it, expect } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import RulesView from '../RulesView.vue'
import { makeRouter } from '../../test-utils/router'

describe('RulesView', () => {
  it('sets document title to Rules — Crew Predictions', async () => {
    mount(RulesView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    expect(document.title).toBe('Rules — Crew Predictions')
  })

  it('renders an Aces Radio section', () => {
    const wrapper = mount(RulesView, { global: { plugins: [makeRouter()] } })
    expect(wrapper.text()).toContain('Aces Radio')
  })

  it('renders an Upper 90 Club section', () => {
    const wrapper = mount(RulesView, { global: { plugins: [makeRouter()] } })
    expect(wrapper.text()).toContain('Upper 90 Club')
  })

  it('shows the +15 exact score rule', () => {
    const wrapper = mount(RulesView, { global: { plugins: [makeRouter()] } })
    expect(wrapper.text()).toContain('+15')
  })

  it('shows the −15 flipped scoreline rule', () => {
    const wrapper = mount(RulesView, { global: { plugins: [makeRouter()] } })
    expect(wrapper.text()).toContain('-15')
  })

  it('shows the +10 correct winner rule', () => {
    const wrapper = mount(RulesView, { global: { plugins: [makeRouter()] } })
    expect(wrapper.text()).toContain('+10')
  })

  it('renders a Grouchy section', () => {
    const wrapper = mount(RulesView, { global: { plugins: [makeRouter()] } })
    expect(wrapper.text()).toContain('Grouchy')
  })

  it('shows the 5 outcome categories for Grouchy', () => {
    const wrapper = mount(RulesView, { global: { plugins: [makeRouter()] } })
    const text = wrapper.text()
    expect(text).toContain('Win by 2')
    expect(text).toContain('Win by 1')
    expect(text).toContain('Draw')
    expect(text).toContain('Lose by 1')
    expect(text).toContain('Lose by 2')
  })

  it('explains Grouchy gives 1 point per correct category', () => {
    const wrapper = mount(RulesView, { global: { plugins: [makeRouter()] } })
    expect(wrapper.text()).toContain('+1')
  })
})
