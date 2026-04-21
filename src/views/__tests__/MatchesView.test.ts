import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { ref } from 'vue'
import MatchesView from '../MatchesView.vue'

const loggedInProvide = {
  currentUser: ref({ handle: 'testfan@crew.mock', emailVerified: true }),
}

const mockMatches = [
  { id: 'match-past', homeTeam: 'New England Revolution', awayTeam: 'Columbus Crew', kickoff: '2026-04-18T23:30:00Z', status: 'STATUS_FULL_TIME', homeScore: '2', awayScore: '1' },
  { id: 'match-1', homeTeam: 'Columbus Crew', awayTeam: 'LA Galaxy', kickoff: '2026-04-22T23:30:00Z', status: 'STATUS_SCHEDULED', homeScore: '', awayScore: '' },
  { id: 'match-2', homeTeam: 'Columbus Crew', awayTeam: 'Philadelphia Union', kickoff: '2026-04-25T23:30:00Z', status: 'STATUS_SCHEDULED', homeScore: '', awayScore: '' },
]

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
    ok: true,
    json: () => Promise.resolve({ matches: mockMatches, predictions: {} }),
  }))
})

describe('MatchesView', () => {
  it('sets document title to Upcoming — Crew Predictions', async () => {
    mount(MatchesView)
    await flushPromises()
    expect(document.title).toBe('Upcoming — Crew Predictions')
  })

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
    const wrapper = mount(MatchesView, { global: { provide: loggedInProvide } })
    await flushPromises()
    const card = wrapper.findAll('[data-testid="match-card"]')[0]
    expect(card.find('input[name="home_goals"]').exists()).toBe(true)
    expect(card.find('input[name="away_goals"]').exists()).toBe(true)
    expect(card.find('button').text()).toBe('Predict')
  })

  it('shows saved prediction score after submitting', async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ matches: mockMatches, predictions: {} }) })
      .mockResolvedValueOnce({ ok: true })
    vi.stubGlobal('fetch', fetchMock)

    const wrapper = mount(MatchesView, { global: { provide: loggedInProvide } })
    await flushPromises()

    const card = wrapper.findAll('[data-testid="match-card"]')[0]
    await card.find('input[name="home_goals"]').setValue('3')
    await card.find('input[name="away_goals"]').setValue('1')
    await card.find('button').trigger('click')
    await flushPromises()

    // New layout: Columbus Crew [3] vs [1] LA Galaxy
    expect(card.text()).toMatch(/Columbus Crew\s*3\s*vs\s*1\s*LA Galaxy/)
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
    // New layout: Columbus Crew [2] vs [0] LA Galaxy
    expect(card.text()).toMatch(/Columbus Crew\s*2\s*vs\s*0\s*LA Galaxy/)
  })

  it('shows a Results section for completed matches', async () => {
    const wrapper = mount(MatchesView)
    await flushPromises()
    expect(wrapper.text()).toContain('Results')
  })

  it('results section shows most recent match first', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        matches: [
          { id: 'older', homeTeam: 'CF Montréal', awayTeam: 'Columbus Crew', kickoff: '2026-04-10T23:30:00Z', status: 'STATUS_FULL_TIME', homeScore: '0', awayScore: '2' },
          { id: 'newer', homeTeam: 'Columbus Crew', awayTeam: 'Atlanta United', kickoff: '2026-04-17T23:30:00Z', status: 'STATUS_FULL_TIME', homeScore: '3', awayScore: '1' },
        ],
        predictions: {},
      }),
    }))
    const wrapper = mount(MatchesView)
    await flushPromises()
    const cards = wrapper.findAll('[data-testid="result-card"]')
    expect(cards[0].text()).toContain('Atlanta United')
    expect(cards[1].text()).toContain('CF Montréal')
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

  it('shows final score between team names on result cards', async () => {
    const wrapper = mount(MatchesView)
    await flushPromises()
    const card = wrapper.find('[data-testid="result-card"]')
    const text = card.text().replace(/\s+/g, ' ')
    // Score must appear between the two team names
    const neIdx = text.indexOf('New England Revolution')
    const clbIdx = text.indexOf('Columbus Crew')
    const scoreIdx = text.indexOf('2')
    expect(neIdx).toBeGreaterThanOrEqual(0)
    expect(scoreIdx).toBeGreaterThan(neIdx)
    expect(clbIdx).toBeGreaterThan(scoreIdx)
  })

  it('result card matchup line contains score inline', async () => {
    const wrapper = mount(MatchesView)
    await flushPromises()
    const matchup = wrapper.find('[data-testid="result-card"] [data-testid="matchup"]')
    expect(matchup.text()).toMatch(/New England Revolution\s*2\s*vs\s*1\s*Columbus Crew/i)
  })

  it('logged-out user sees a disabled Predict button, not a Sign in link', async () => {
    const wrapper = mount(MatchesView)
    await flushPromises()
    const card = wrapper.findAll('[data-testid="match-card"]')[0]
    const btn = card.find('button')
    expect(btn.text()).toBe('Predict')
    expect(btn.attributes('disabled')).toBeDefined()
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
