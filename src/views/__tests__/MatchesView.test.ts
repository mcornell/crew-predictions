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

  it('each card has home_goals and away_goals inputs and a Lock In button', async () => {
    const wrapper = mount(MatchesView)
    await flushPromises()
    const card = wrapper.findAll('[data-testid="match-card"]')[0]
    expect(card.find('input[name="home_goals"]').exists()).toBe(true)
    expect(card.find('input[name="away_goals"]').exists()).toBe(true)
    expect(card.find('button').text()).toBe('Lock In')
  })

  it('shows saved prediction score after submitting', async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ matches: mockMatches, predictions: {} }) })
      .mockResolvedValueOnce({ ok: true })
    vi.stubGlobal('fetch', fetchMock)

    const wrapper = mount(MatchesView)
    await flushPromises()

    const card = wrapper.findAll('[data-testid="match-card"]')[0]
    await card.find('input[name="home_goals"]').setValue('3')
    await card.find('input[name="away_goals"]').setValue('1')
    await card.find('button').trigger('click')
    await flushPromises()

    expect(card.text()).toContain('3 – 1')
  })

  it('shows existing prediction from initial load', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        matches: mockMatches,
        predictions: { 'match-1': { homeGoals: 2, awayGoals: 0 } },
      }),
    }))
    const wrapper = mount(MatchesView)
    await flushPromises()
    const card = wrapper.findAll('[data-testid="match-card"]')[0]
    expect(card.text()).toContain('2 – 0')
  })
})
