import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import MatchesView from '../MatchesView.vue'

const mockMatches = [
  { id: 'match-past', homeTeam: 'New England Revolution', awayTeam: 'Columbus Crew', kickoff: '2026-04-18T23:30:00Z', status: 'STATUS_FULL_TIME' },
  { id: 'match-1', homeTeam: 'Columbus Crew', awayTeam: 'LA Galaxy', kickoff: '2026-04-22T23:30:00Z', status: 'STATUS_SCHEDULED' },
  { id: 'match-2', homeTeam: 'Columbus Crew', awayTeam: 'Philadelphia Union', kickoff: '2026-04-25T23:30:00Z', status: 'STATUS_SCHEDULED' },
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

  it('shows a Results section for completed matches', async () => {
    const wrapper = mount(MatchesView)
    await flushPromises()
    expect(wrapper.text()).toContain('Results')
  })

  it('completed match appears in results section, not upcoming', async () => {
    const wrapper = mount(MatchesView)
    await flushPromises()
    const resultsSection = wrapper.find('[data-testid="results-section"]')
    expect(resultsSection.exists()).toBe(true)
    expect(resultsSection.text()).toContain('New England Revolution')
  })

  it('upcoming match does not appear in results section', async () => {
    const wrapper = mount(MatchesView)
    await flushPromises()
    const resultsSection = wrapper.find('[data-testid="results-section"]')
    expect(resultsSection.text()).not.toContain('LA Galaxy')
  })

  it('match more than 7 days away is not shown in upcoming', async () => {
    const farFuture = new Date()
    farFuture.setDate(farFuture.getDate() + 10)
    const farMatch = { id: 'far', homeTeam: 'Columbus Crew', awayTeam: 'Inter Miami', kickoff: farFuture.toISOString(), status: 'STATUS_SCHEDULED' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ matches: [farMatch], predictions: {} }),
    }))
    const wrapper = mount(MatchesView)
    await flushPromises()
    expect(wrapper.findAll('[data-testid="match-card"]')).toHaveLength(0)
  })
})
