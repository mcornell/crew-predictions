import { describe, it, expect } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import RulesView from '../RulesView.vue'

function makeRouter() {
  return createRouter({
    history: createMemoryHistory(),
    routes: [{ path: '/rules', component: RulesView }],
  })
}

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
})
