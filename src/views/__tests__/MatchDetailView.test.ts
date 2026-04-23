import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import MatchDetailView from '../MatchDetailView.vue'

const mockMatch = {
  id: 'm-test',
  homeTeam: 'Columbus Crew',
  awayTeam: 'FC Dallas',
  kickoff: '2026-04-20T19:00:00Z',
  homeScore: '2',
  awayScore: '1',
}

const mockPredictions = [
  { userID: 'google:u1', handle: 'fan1@bsky.mock', homeGoals: 2, awayGoals: 1, acesRadioPoints: 15, upper90ClubPoints: 10 },
  { userID: 'google:u2', handle: 'fan2@bsky.mock', homeGoals: 0, awayGoals: 0, acesRadioPoints: 0, upper90ClubPoints: 0 },
]

const mockScoringFormats = [
  { key: 'acesRadio', label: 'Aces Radio' },
  { key: 'upper90Club', label: 'Upper 90 Club' },
]

function makeRouter(matchId = 'm-test') {
  const stub = { template: '<div />' }
  const router = createRouter({
    history: createMemoryHistory(),
    routes: [
      { path: '/matches', component: stub },
      { path: '/matches/:matchId', component: MatchDetailView },
    ],
  })
  router.push(`/matches/${matchId}`)
  return router
}

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
    ok: true,
    json: () => Promise.resolve({ match: mockMatch, predictions: mockPredictions, scoringFormats: mockScoringFormats }),
  }))
})

afterEach(() => {
  vi.restoreAllMocks()
})

describe('MatchDetailView', () => {
  it('shows match header with team names and score', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, {
      global: { plugins: [router] },
    })
    await flushPromises()
    expect(wrapper.text()).toContain('Columbus Crew')
    expect(wrapper.text()).toContain('FC Dallas')
    expect(wrapper.find('[data-testid="match-score"]').text()).toContain('2')
  })

  it('renders a row for each predictor', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, {
      global: { plugins: [router] },
    })
    await flushPromises()
    const rows = wrapper.findAll('[data-testid="prediction-row"]')
    expect(rows.length).toBe(2)
  })

  it('shows sort buttons for each scoring format', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, {
      global: { plugins: [router] },
    })
    await flushPromises()
    expect(wrapper.text()).toContain('Aces Radio')
    expect(wrapper.text()).toContain('Upper 90 Club')
  })

  it('sorts by Aces Radio descending by default — highest points first', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, {
      global: { plugins: [router] },
    })
    await flushPromises()
    const rows = wrapper.findAll('[data-testid="prediction-row"]')
    expect(rows[0].text()).toContain('fan1@bsky.mock')
  })

  it('shows empty state when no predictions were made', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ match: mockMatch, predictions: [], scoringFormats: mockScoringFormats }),
    }))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, {
      global: { plugins: [router] },
    })
    await flushPromises()
    expect(wrapper.text()).toContain('No predictions were made for this match')
  })

  it('shows a back link to /matches', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, {
      global: { plugins: [router] },
    })
    await flushPromises()
    const backLink = wrapper.find('[data-testid="back-link"]')
    expect(backLink.exists()).toBe(true)
    expect(backLink.attributes('href')).toContain('/matches')
  })
})
