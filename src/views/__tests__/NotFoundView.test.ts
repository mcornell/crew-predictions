import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import NotFoundView from '../NotFoundView.vue'

describe('NotFoundView', () => {
  it('renders the 404 heading', () => {
    const wrapper = mount(NotFoundView)
    expect(wrapper.find('[data-testid="not-found"]').exists()).toBe(true)
    expect(wrapper.find('h1').text()).toBe('404')
  })

  it('contains a link back to matches', () => {
    const wrapper = mount(NotFoundView)
    expect(wrapper.find('a[href="/matches"]').exists()).toBe(true)
  })
})
