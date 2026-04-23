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
  { userID: 'google:u1', handle: 'fan1@bsky.mock', homeGoals: 2, awayGoals: 1, acesRadioPoints: 15, upper90ClubPoints: 3 },
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
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.text()).toContain('Columbus Crew')
    expect(wrapper.text()).toContain('FC Dallas')
    expect(wrapper.find('[data-testid="match-score"]').text()).toContain('2')
  })

  it('renders a row for each predictor', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    const rows = wrapper.findAll('[data-testid="prediction-row"]')
    expect(rows).toHaveLength(2)
  })

  it('shows column headers for both scoring formats', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="sort-aces"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="sort-upper90"]').exists()).toBe(true)
  })

  it('shows both Aces Radio and Upper 90 Club points in each row', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    const rows = wrapper.findAll('[data-testid="prediction-row"]')
    expect(rows[0].find('[data-testid="prediction-aces-points"]').text()).toBe('15')
    expect(rows[0].find('[data-testid="prediction-upper90-points"]').text()).toBe('3')
  })

  it('shows the predicted score in each row', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    const rows = wrapper.findAll('[data-testid="prediction-row"]')
    expect(rows[0].find('[data-testid="prediction-score"]').text()).toContain('2')
  })

  it('sorts by Aces Radio descending by default — highest points first', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    const rows = wrapper.findAll('[data-testid="prediction-row"]')
    expect(rows[0].text()).toContain('fan1@bsky.mock')
  })

  it('sorts by Upper 90 Club when that column header is clicked', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    await wrapper.find('[data-testid="sort-upper90"]').trigger('click')
    const rows = wrapper.findAll('[data-testid="prediction-row"]')
    expect(rows[0].find('[data-testid="prediction-upper90-points"]').text()).toBe('3')
  })

  it('rank is dynamic — ties share rank', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: mockMatch,
        predictions: [
          { userID: 'u1', handle: 'Fan1', homeGoals: 2, awayGoals: 1, acesRadioPoints: 10, upper90ClubPoints: 2 },
          { userID: 'u2', handle: 'Fan2', homeGoals: 2, awayGoals: 0, acesRadioPoints: 10, upper90ClubPoints: 1 },
          { userID: 'u3', handle: 'Fan3', homeGoals: 1, awayGoals: 0, acesRadioPoints: 0, upper90ClubPoints: 0 },
        ],
        scoringFormats: mockScoringFormats,
      }),
    }))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    const rows = wrapper.findAll('[data-testid="prediction-row"]')
    expect(rows[0].find('[data-testid="prediction-rank"]').text()).toBe('1')
    expect(rows[1].find('[data-testid="prediction-rank"]').text()).toBe('1')
    expect(rows[2].find('[data-testid="prediction-rank"]').text()).toBe('3')
  })

  it('shows empty state when no predictions were made', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ match: mockMatch, predictions: [], scoringFormats: mockScoringFormats }),
    }))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.text()).toContain('No predictions were made for this match')
  })

  it('shows loading state before fetch resolves', async () => {
    vi.stubGlobal('fetch', vi.fn().mockReturnValue(new Promise(() => {})))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    expect(wrapper.find('[data-testid="loading"]').exists()).toBe(true)
  })

  it('shows error state when fetch fails', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false }))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="error"]').exists()).toBe(true)
  })

  it('shows a back link to /matches', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    const backLink = wrapper.find('[data-testid="back-link"]')
    expect(backLink.exists()).toBe(true)
    expect(backLink.attributes('href')).toContain('/matches')
  })
})
