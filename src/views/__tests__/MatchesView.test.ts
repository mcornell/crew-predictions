import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import MatchesView from '../MatchesView.vue'

const mockMatches = [
  { id: 'match-1', homeTeam: 'Portland Timbers', awayTeam: 'Columbus Crew', kickoff: '2026-05-01T20:00:00Z', status: 'STATUS_SCHEDULED' },
  { id: 'match-2', homeTeam: 'Sporting Kansas City', awayTeam: 'Columbus Crew', kickoff: '2026-05-08T23:00:00Z', status: 'STATUS_SCHEDULED' },
]

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
    ok: true,
    json: () => Promise.resolve({ matches: mockMatches, predictions: {} }),
  }))
})

describe('MatchesView', () => {
  it('renders an Upcoming heading', async () => {
    const wrapper = mount(MatchesView)
    await flushPromises()
    expect(wrapper.find('h1').text()).toBe('Upcoming')
  })

  it('renders a card for each match', async () => {
    const wrapper = mount(MatchesView)
    await flushPromises()
    expect(wrapper.findAll('[data-testid="match-card"]')).toHaveLength(2)
  })

  it('shows Columbus Crew in at least one card', async () => {
    const wrapper = mount(MatchesView)
    await flushPromises()
    expect(wrapper.text()).toContain('Columbus Crew')
  })
})
