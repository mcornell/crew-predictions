import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import MatchDetailView from '../MatchDetailView.vue'
import { POLL_INTERVAL_MS } from '../../utils/pollScheduler'

const mockMatch = {
  id: 'm-test',
  homeTeam: 'Columbus Crew',
  awayTeam: 'FC Dallas',
  kickoff: '2026-04-20T19:00:00Z',
  homeScore: '2',
  awayScore: '1',
}

const mockPredictions = [
  { userID: 'google:u1', handle: 'fan1@bsky.mock', homeGoals: 2, awayGoals: 1, acesRadioPoints: 15, upper90ClubPoints: 3, grouchyPoints: 1 },
  { userID: 'google:u2', handle: 'fan2@bsky.mock', homeGoals: 0, awayGoals: 0, acesRadioPoints: 0, upper90ClubPoints: 0, grouchyPoints: 0 },
]

const mockScoringFormats = [
  { key: 'acesRadio', label: 'Aces Radio' },
  { key: 'upper90Club', label: 'Upper 90 Club' },
  { key: 'grouchy', label: 'Grouchy™' },
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

  it('renders mobile sort buttons for both scoring formats', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="mobile-sort-aces"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="mobile-sort-upper90"]').exists()).toBe(true)
  })

  it('mobile sort upper90 button triggers sort change', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    await wrapper.find('[data-testid="mobile-sort-upper90"]').trigger('click')
    const rows = wrapper.findAll('[data-testid="prediction-row"]')
    expect(rows[0].find('[data-testid="prediction-upper90-points"]').text()).toBe('3')
  })

  it('shows Grouchy points in each row', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    const rows = wrapper.findAll('[data-testid="prediction-row"]')
    expect(rows[0].find('[data-testid="prediction-grouchy-points"]').text()).toBe('1')
    expect(rows[1].find('[data-testid="prediction-grouchy-points"]').text()).toBe('0')
  })

  it('sorts by Grouchy when Grouchy column header is clicked', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: mockMatch,
        predictions: [
          { userID: 'u1', handle: 'fan1', homeGoals: 1, awayGoals: 0, acesRadioPoints: 0, upper90ClubPoints: 0, grouchyPoints: 1 },
          { userID: 'u2', handle: 'fan2', homeGoals: 0, awayGoals: 1, acesRadioPoints: 0, upper90ClubPoints: 0, grouchyPoints: 0 },
        ],
        scoringFormats: mockScoringFormats,
      }),
    }))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    await wrapper.find('[data-testid="sort-grouchy"]').trigger('click')
    const rows = wrapper.findAll('[data-testid="prediction-row"]')
    expect(rows[0].find('[data-testid="prediction-grouchy-points"]').text()).toBe('1')
  })

  it('renders mobile sort button for Grouchy', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="mobile-sort-grouchy"]').exists()).toBe(true)
  })

  it('mobile sort Grouchy button triggers sort change', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    await wrapper.find('[data-testid="mobile-sort-grouchy"]').trigger('click')
    const rows = wrapper.findAll('[data-testid="prediction-row"]')
    expect(rows[0].find('[data-testid="prediction-grouchy-points"]').text()).toBe('1')
  })

  it('mobile sort Aces button switches back to Aces Radio', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    await wrapper.find('[data-testid="mobile-sort-upper90"]').trigger('click')
    await wrapper.find('[data-testid="mobile-sort-aces"]').trigger('click')
    const rows = wrapper.findAll('[data-testid="prediction-row"]')
    expect(rows[0].find('[data-testid="prediction-aces-points"]').text()).toBe('15')
  })

  it('desktop sort Aces button switches back to Aces Radio from another sort', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    await wrapper.find('[data-testid="sort-upper90"]').trigger('click')
    await wrapper.find('[data-testid="sort-aces"]').trigger('click')
    const rows = wrapper.findAll('[data-testid="prediction-row"]')
    expect(rows[0].find('[data-testid="prediction-aces-points"]').text()).toBe('15')
  })

  it('shows live indicator when match state is "in"', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: { ...mockMatch, state: 'in', displayClock: '67\'', status: 'STATUS_IN_PROGRESS' },
        predictions: [],
        scoringFormats: mockScoringFormats,
      }),
    }))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="live-indicator-detail"]').exists()).toBe(true)
    expect(wrapper.text()).toContain("67'")
  })

  it('shows halftime indicator when match status is STATUS_HALFTIME', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: { ...mockMatch, state: 'in', status: 'STATUS_HALFTIME' },
        predictions: [],
        scoringFormats: mockScoringFormats,
      }),
    }))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.text()).toContain('HT')
  })

  it('shows projected label when isProjected is true', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: mockMatch,
        predictions: mockPredictions,
        isProjected: true,
        scoringFormats: mockScoringFormats,
      }),
    }))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="projected-label"]').exists()).toBe(true)
  })

  it('formatKickoff returns empty string for invalid date', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: { ...mockMatch, kickoff: 'not-a-date' },
        predictions: [],
        scoringFormats: mockScoringFormats,
      }),
    }))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.find('.match-meta').text()).toBe('')
  })

  it('poll timer callback re-fetches detail when it fires', async () => {
    vi.useFakeTimers()
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: { ...mockMatch, state: 'in', status: 'STATUS_IN_PROGRESS' },
        predictions: [],
        scoringFormats: mockScoringFormats,
      }),
    })
    vi.stubGlobal('fetch', fetchMock)
    const router = makeRouter()
    await router.isReady()
    mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    const callsBefore = fetchMock.mock.calls.length
    await vi.advanceTimersByTimeAsync(POLL_INTERVAL_MS)
    await flushPromises()
    expect(fetchMock.mock.calls.length).toBeGreaterThan(callsBefore)
    vi.useRealTimers()
  })

  it('schedules a poll when kickoff is in the future but outside the active window', async () => {
    vi.useFakeTimers()
    const futureKickoff = new Date(Date.now() + 3 * 60 * 60 * 1000).toISOString()
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: { ...mockMatch, kickoff: futureKickoff, state: '', status: 'STATUS_SCHEDULED' },
        predictions: [],
        scoringFormats: mockScoringFormats,
      }),
    }))
    const setTimeoutSpy = vi.spyOn(global, 'setTimeout')
    const router = makeRouter()
    await router.isReady()
    mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(setTimeoutSpy).toHaveBeenCalled()
    vi.useRealTimers()
  })

  it('future-window poll timer fires and re-fetches', async () => {
    vi.useFakeTimers()
    // Kickoff 31 minutes from now → msUntilActiveWindow ≈ 60 000 ms (1 minute before active window)
    const kickoff = new Date(Date.now() + 31 * 60 * 1000).toISOString()
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: { ...mockMatch, kickoff, state: 'pre', status: 'STATUS_SCHEDULED' },
        predictions: [],
        scoringFormats: mockScoringFormats,
      }),
    })
    vi.stubGlobal('fetch', fetchMock)
    const router = makeRouter()
    await router.isReady()
    mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    const callsBefore = fetchMock.mock.calls.length
    await vi.advanceTimersByTimeAsync(61_000)
    await flushPromises()
    expect(fetchMock.mock.calls.length).toBeGreaterThan(callsBefore)
    vi.useRealTimers()
  })

  it('shows venue on match detail page', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: { ...mockMatch, venue: 'ScottsMiracle-Gro Field' },
        predictions: mockPredictions,
        scoringFormats: mockScoringFormats,
      }),
    }))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="match-detail-venue"]').text()).toBe('ScottsMiracle-Gro Field')
  })

  it('shows ESPN link on match detail page', async () => {
    const router = makeRouter('m-test')
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    const link = wrapper.find('[data-testid="espn-link"]')
    expect(link.exists()).toBe(true)
    expect(link.attributes('href')).toContain('gameId/m-test')
  })

  it('shows record and form on match detail page', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: { ...mockMatch, homeRecord: '5-3-2', awayRecord: '4-4-2', homeForm: 'WWWLL', awayForm: 'LWDWL' },
        predictions: [],
        scoringFormats: mockScoringFormats,
      }),
    }))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="home-record"]').text()).toBe('5-3-2')
    expect(wrapper.find('[data-testid="home-form"]').text()).toBe('WWWLL')
    vi.restoreAllMocks()
  })

  it('shows formatted attendance when match has attendance', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: { ...mockMatch, attendance: 19903 },
        predictions: [],
        scoringFormats: mockScoringFormats,
      }),
    }))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="match-detail-attendance"]').text()).toBe('19,903')
  })

  it('does not show attendance element when attendance is 0', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: { ...mockMatch, attendance: 0 },
        predictions: [],
        scoringFormats: mockScoringFormats,
      }),
    }))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="match-detail-attendance"]').exists()).toBe(false)
  })

  it('shows event timeline when match has displayable events', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: {
          ...mockMatch,
          events: [
            { clock: "4'", typeID: 'goal', team: 'Columbus Crew', players: ['Max Arfsten'] },
            { clock: "90'+4'", typeID: 'red-card', team: 'Philadelphia Union', players: ['Japhet Sery'] },
          ],
        },
        predictions: [],
        scoringFormats: mockScoringFormats,
      }),
    }))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="match-events"]').exists()).toBe(true)
    expect(wrapper.findAll('[data-testid="match-event"]')).toHaveLength(2)
  })

  it('does not show event timeline when match has no events', async () => {
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="match-events"]').exists()).toBe(false)
  })

  it('filters out non-displayable event types (kickoff, halftime, start-2nd-half, end-regular-time)', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: {
          ...mockMatch,
          events: [
            { clock: '', typeID: 'kickoff', team: '', players: [] },
            { clock: "4'", typeID: 'goal', team: 'Columbus Crew', players: ['Max Arfsten'] },
            { clock: '', typeID: 'halftime', team: '', players: [] },
            { clock: '', typeID: 'start-2nd-half', team: '', players: [] },
            { clock: '', typeID: 'end-regular-time', team: '', players: [] },
          ],
        },
        predictions: [],
        scoringFormats: mockScoringFormats,
      }),
    }))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    expect(wrapper.findAll('[data-testid="match-event"]')).toHaveLength(1)
  })

  it('renders substitution with on and off players', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: {
          ...mockMatch,
          events: [
            { clock: "63'", typeID: 'substitution', team: 'Columbus Crew', players: ['Steven Moreira', 'Hugo Picard'] },
          ],
        },
        predictions: [],
        scoringFormats: mockScoringFormats,
      }),
    }))
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    const eventEl = wrapper.find('[data-testid="match-event"]')
    expect(eventEl.text()).toContain('Steven Moreira')
    expect(eventEl.text()).toContain('Hugo Picard')
  })

  it('clears poll timer on unmount when a poll was scheduled', async () => {
    vi.useFakeTimers()
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        match: { ...mockMatch, state: 'in', status: 'STATUS_IN_PROGRESS' },
        predictions: [],
        scoringFormats: mockScoringFormats,
      }),
    }))
    const clearTimeoutSpy = vi.spyOn(global, 'clearTimeout')
    const router = makeRouter()
    await router.isReady()
    const wrapper = mount(MatchDetailView, { global: { plugins: [router] } })
    await flushPromises()
    wrapper.unmount()
    expect(clearTimeoutSpy).toHaveBeenCalled()
    vi.useRealTimers()
  })
})
