import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { ref, type Ref } from 'vue'
import MatchesView from '../MatchesView.vue'
import { makeRouter } from '../../test-utils/router'

type User = { handle: string; emailVerified: boolean } | null
type Provide = { currentUser: Ref<User> }

const loggedOutProvide: Provide = { currentUser: ref(null) }
const loggedInProvide: Provide = {
  currentUser: ref({ handle: 'testfan@crew.mock', emailVerified: true }),
}

function mountMatches(provide: Provide = loggedOutProvide) {
  return mount(MatchesView, { global: { plugins: [makeRouter()], provide } })
}

function futureKickoff(hoursFromNow: number): string {
  return new Date(Date.now() + hoursFromNow * 60 * 60 * 1000).toISOString()
}
function pastKickoff(hoursAgo: number): string {
  return new Date(Date.now() - hoursAgo * 60 * 60 * 1000).toISOString()
}

const mockMatches = [
  { id: 'match-past', homeTeam: 'New England Revolution', awayTeam: 'Columbus Crew', kickoff: pastKickoff(96), status: 'STATUS_FULL_TIME', state: 'post', homeScore: '2', awayScore: '1' },
  { id: 'match-1', homeTeam: 'Columbus Crew', awayTeam: 'LA Galaxy', kickoff: futureKickoff(24), status: 'STATUS_SCHEDULED', state: 'pre', homeScore: '', awayScore: '' },
  { id: 'match-2', homeTeam: 'Columbus Crew', awayTeam: 'Philadelphia Union', kickoff: futureKickoff(72), status: 'STATUS_SCHEDULED', state: 'pre', homeScore: '', awayScore: '' },
]

const liveMatch = { id: 'match-live', homeTeam: 'Columbus Crew', awayTeam: 'Philadelphia Union', kickoff: pastKickoff(1), status: 'STATUS_FIRST_HALF', state: 'in', homeScore: '', awayScore: '' }

beforeEach(() => {
  localStorage.clear()
  vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
    ok: true,
    json: () => Promise.resolve({ matches: mockMatches, predictions: {} }),
  }))
})

describe('MatchesView', () => {
  it('sets document title to Upcoming — Crew Predictions', async () => {
    mountMatches()
    await flushPromises()
    expect(document.title).toBe('Upcoming — Crew Predictions')
  })

  it('renders an Upcoming heading', async () => {
    const wrapper = mountMatches()
    await flushPromises()
    expect(wrapper.find('h1').text()).toBe('Upcoming')
  })

  it('renders a card for each match', async () => {
    const wrapper = mountMatches()
    await flushPromises()
    expect(wrapper.findAll('[data-testid="match-card"]')).toHaveLength(2)
  })

  it('shows Columbus Crew in at least one card', async () => {
    const wrapper = mountMatches()
    await flushPromises()
    expect(wrapper.text()).toContain('Columbus Crew')
  })

  it('each card has home_goals and away_goals inputs and a Lock In button', async () => {
    const wrapper = mountMatches(loggedInProvide)
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

    const wrapper = mountMatches(loggedInProvide)
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
    const wrapper = mountMatches()
    await flushPromises()
    const card = wrapper.findAll('[data-testid="match-card"]')[0]
    // New layout: Columbus Crew [2] vs [0] LA Galaxy
    expect(card.text()).toMatch(/Columbus Crew\s*2\s*vs\s*0\s*LA Galaxy/)
  })

  it('shows a Results section for completed matches', async () => {
    const wrapper = mountMatches()
    await flushPromises()
    expect(wrapper.text()).toContain('Results')
  })

  it('results section shows most recent match first', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        matches: [
          { id: 'older', homeTeam: 'CF Montréal', awayTeam: 'Columbus Crew', kickoff: '2026-04-10T23:30:00Z', status: 'STATUS_FULL_TIME', state: 'post', homeScore: '0', awayScore: '2' },
          { id: 'newer', homeTeam: 'Columbus Crew', awayTeam: 'Atlanta United', kickoff: '2026-04-17T23:30:00Z', status: 'STATUS_FULL_TIME', state: 'post', homeScore: '3', awayScore: '1' },
        ],
        predictions: {},
      }),
    }))
    const wrapper = mountMatches()
    await flushPromises()
    const cards = wrapper.findAll('[data-testid="result-card"]')
    expect(cards[0].text()).toContain('Atlanta United')
    expect(cards[1].text()).toContain('CF Montréal')
  })

  it('completed match appears in results section, not upcoming', async () => {
    const wrapper = mountMatches()
    await flushPromises()
    const resultsSection = wrapper.find('[data-testid="results-section"]')
    expect(resultsSection.exists()).toBe(true)
    expect(resultsSection.text()).toContain('New England Revolution')
  })

  it('upcoming match does not appear in results section', async () => {
    const wrapper = mountMatches()
    await flushPromises()
    const resultsSection = wrapper.find('[data-testid="results-section"]')
    expect(resultsSection.text()).not.toContain('LA Galaxy')
  })

  it('shows final score between team names on result cards', async () => {
    const wrapper = mountMatches()
    await flushPromises()
    const card = wrapper.find('[data-testid="result-card"]')
    const text = card.text().replace(/\s+/g, ' ')
    const neIdx = text.indexOf('New England Revolution')
    const clbIdx = text.indexOf('Columbus Crew')
    const scoreIdx = text.indexOf('2')
    expect(neIdx).toBeGreaterThanOrEqual(0)
    expect(scoreIdx).toBeGreaterThan(neIdx)
    expect(clbIdx).toBeGreaterThan(scoreIdx)
  })

  it('result card matchup line contains score inline', async () => {
    const wrapper = mountMatches()
    await flushPromises()
    const matchup = wrapper.find('[data-testid="result-card"] [data-testid="matchup"]')
    expect(matchup.text()).toMatch(/New England Revolution\s*2\s*vs\s*1\s*Columbus Crew/i)
  })

  it('logged-out user sees an enabled Predict button with score inputs', async () => {
    const wrapper = mountMatches()
    await flushPromises()
    const card = wrapper.findAll('[data-testid="match-card"]')[0]
    expect(card.find('input[name="home_goals"]').exists()).toBe(true)
    expect(card.find('input[name="away_goals"]').exists()).toBe(true)
    const btn = card.find('button')
    expect(btn.text()).toBe('Predict')
    expect(btn.attributes('disabled')).toBeUndefined()
  })

  it('guest prediction is saved to localStorage after submit', async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ matches: mockMatches, predictions: {} }) })
    vi.stubGlobal('fetch', fetchMock)
    const storageSpy = vi.spyOn(Storage.prototype, 'setItem')

    const wrapper = mountMatches()
    await flushPromises()

    const card = wrapper.findAll('[data-testid="match-card"]')[0]
    await card.find('input[name="home_goals"]').setValue('2')
    await card.find('input[name="away_goals"]').setValue('0')
    await card.find('button').trigger('click')
    await flushPromises()

    expect(storageSpy).toHaveBeenCalledWith(
      'guestPredictions',
      expect.stringContaining('"match-1"'),
    )
    storageSpy.mockRestore()
  })

  it('shows sign-in nudge after guest submits a prediction', async () => {
    const wrapper = mountMatches()
    await flushPromises()

    const card = wrapper.findAll('[data-testid="match-card"]')[0]
    await card.find('input[name="home_goals"]').setValue('2')
    await card.find('input[name="away_goals"]').setValue('0')
    await card.find('button').trigger('click')
    await flushPromises()

    expect(wrapper.find('[data-testid="guest-nudge"]').exists()).toBe(true)
  })

  it('guest prediction loaded from localStorage on mount', async () => {
    localStorage.setItem('guestPredictions', JSON.stringify({ 'match-1': { homeGoals: 3, awayGoals: 2 } }))

    const wrapper = mountMatches()
    await flushPromises()

    const card = wrapper.findAll('[data-testid="match-card"]')[0]
    expect(card.text()).toMatch(/Columbus Crew\s*3\s*vs\s*2\s*LA Galaxy/)

    localStorage.removeItem('guestPredictions')
  })

  it('shows LIVE indicator on match card when state is "in"', async () => {
    const now = new Date()
    const liveMatch = { id: 'live-1', homeTeam: 'Columbus Crew', awayTeam: 'FC Dallas', kickoff: new Date(now.getTime() - 60 * 60 * 1000).toISOString(), status: 'STATUS_SCHEDULED', homeScore: '', awayScore: '', state: 'in' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ matches: [liveMatch], predictions: {} }),
    }))
    const wrapper = mountMatches()
    await flushPromises()
    expect(wrapper.find('[data-testid="live-indicator"]').exists()).toBe(true)
  })

  it('does not show LIVE indicator when state is not "in"', async () => {
    const wrapper = mountMatches()
    await flushPromises()
    expect(wrapper.find('[data-testid="live-indicator"]').exists()).toBe(false)
  })

  it('match more than 8 days away is not shown in upcoming', async () => {
    const farFuture = new Date()
    farFuture.setDate(farFuture.getDate() + 10)
    const farMatch = { id: 'far', homeTeam: 'Columbus Crew', awayTeam: 'Inter Miami', kickoff: farFuture.toISOString(), status: 'STATUS_SCHEDULED', state: 'pre' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ matches: [farMatch], predictions: {} }),
    }))
    const wrapper = mountMatches()
    await flushPromises()
    expect(wrapper.findAll('[data-testid="match-card"]')).toHaveLength(0)
  })

  it('match exactly 8 days away is shown in upcoming', async () => {
    const eightDays = new Date()
    eightDays.setDate(eightDays.getDate() + 8)
    const match = { id: 'm8', homeTeam: 'Columbus Crew', awayTeam: 'Inter Miami', kickoff: eightDays.toISOString(), status: 'STATUS_SCHEDULED', state: 'pre' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ matches: [match], predictions: {} }),
    }))
    const wrapper = mountMatches()
    await flushPromises()
    expect(wrapper.findAll('[data-testid="match-card"]')).toHaveLength(1)
  })

  it('live match shows 0-0 score, not dashes', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ matches: [liveMatch], predictions: {} }),
    }))
    const wrapper = mountMatches()
    await flushPromises()
    const matchup = wrapper.find('[data-testid="now-playing-card"] [data-testid="matchup"]')
    expect(matchup.text()).toMatch(/Columbus Crew\s*0\s*vs\s*0\s*Philadelphia Union/i)
  })

  it('live match does not appear in results', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ matches: [liveMatch], predictions: {} }),
    }))
    const wrapper = mountMatches()
    await flushPromises()
    expect(wrapper.findAll('[data-testid="result-card"]')).toHaveLength(0)
  })

  it('shows countdown on upcoming match card', async () => {
    const wrapper = mountMatches()
    await flushPromises()
    const card = wrapper.findAll('[data-testid="match-card"]')[0]
    expect(card.find('[data-testid="match-countdown"]').exists()).toBe(true)
  })

  it('shows Unlock button after submitting a prediction', async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ matches: mockMatches, predictions: {} }) })
      .mockResolvedValueOnce({ ok: true })
    vi.stubGlobal('fetch', fetchMock)

    const wrapper = mountMatches(loggedInProvide)
    await flushPromises()

    const card = wrapper.findAll('[data-testid="match-card"]')[0]
    await card.find('input[name="home_goals"]').setValue('2')
    await card.find('input[name="away_goals"]').setValue('1')
    await card.find('button').trigger('click')
    await flushPromises()

    expect(card.find('button').text()).toBe('Unlock')
  })

  it('clears the countdown interval on unmount', async () => {
    vi.useFakeTimers()
    const clearSpy = vi.spyOn(globalThis, 'clearInterval')

    const wrapper = mountMatches()
    await flushPromises()

    wrapper.unmount()

    expect(clearSpy).toHaveBeenCalled()
    clearSpy.mockRestore()
    vi.useRealTimers()
  })

  it('clicking Unlock restores score inputs pre-populated with previous pick', async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({ matches: mockMatches, predictions: {} }) })
      .mockResolvedValueOnce({ ok: true })
    vi.stubGlobal('fetch', fetchMock)

    const wrapper = mountMatches(loggedInProvide)
    await flushPromises()

    const card = wrapper.findAll('[data-testid="match-card"]')[0]
    await card.find('input[name="home_goals"]').setValue('2')
    await card.find('input[name="away_goals"]').setValue('1')
    await card.find('button').trigger('click')
    await flushPromises()

    await card.find('button').trigger('click') // Unlock
    await flushPromises()

    expect((card.find('input[name="home_goals"]').element as HTMLInputElement).value).toBe('2')
    expect((card.find('input[name="away_goals"]').element as HTMLInputElement).value).toBe('1')
    expect(card.find('button').text()).toBe('Predict')
  })

  it('delayed match appears in now-playing section, not upcoming or results', async () => {
    const now = new Date()
    const delayedMatch = { id: 'del-1', homeTeam: 'Columbus Crew', awayTeam: 'LA Galaxy', kickoff: new Date(now.getTime() - 60 * 60 * 1000).toISOString(), status: 'STATUS_DELAYED', homeScore: '0', awayScore: '0', state: '' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ matches: [delayedMatch], predictions: {} }) }))
    const wrapper = mountMatches()
    await flushPromises()
    expect(wrapper.find('[data-testid="now-playing-section"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="delayed-indicator"]').exists()).toBe(true)
    expect(wrapper.findAll('[data-testid="match-card"]')).toHaveLength(0)
    expect(wrapper.findAll('[data-testid="result-card"]')).toHaveLength(0)
  })

  it('delayed match card has no Predict or Unlock button', async () => {
    const now = new Date()
    const delayedMatch = { id: 'del-2', homeTeam: 'Columbus Crew', awayTeam: 'LA Galaxy', kickoff: new Date(now.getTime() - 60 * 60 * 1000).toISOString(), status: 'STATUS_DELAYED', homeScore: '0', awayScore: '0', state: '' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ matches: [delayedMatch], predictions: {} }) }))
    const wrapper = mountMatches(loggedInProvide)
    await flushPromises()
    const section = wrapper.find('[data-testid="now-playing-section"]')
    expect(section.find('button').exists()).toBe(false)
  })

  it('in-progress match appears in now-playing section, not upcoming', async () => {
    const now = new Date()
    const liveMatch = { id: 'live-2', homeTeam: 'Columbus Crew', awayTeam: 'FC Dallas', kickoff: new Date(now.getTime() - 60 * 60 * 1000).toISOString(), status: 'STATUS_SCHEDULED', homeScore: '', awayScore: '', state: 'in' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ matches: [liveMatch], predictions: {} }) }))
    const wrapper = mountMatches()
    await flushPromises()
    expect(wrapper.find('[data-testid="now-playing-section"]').exists()).toBe(true)
    expect(wrapper.findAll('[data-testid="match-card"]')).toHaveLength(0)
  })

  it('polls every 60 seconds when a live match is present', async () => {
    vi.useFakeTimers()
    const liveMatches = [...mockMatches, liveMatch]
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ matches: liveMatches, predictions: {} }),
    })
    vi.stubGlobal('fetch', fetchMock)

    mountMatches()
    await flushPromises()
    const callsAfterMount = fetchMock.mock.calls.length

    vi.advanceTimersByTime(60_000)
    await flushPromises()

    expect(fetchMock.mock.calls.length).toBeGreaterThan(callsAfterMount)
    vi.useRealTimers()
  })

  it('does not poll when no match is in active window', async () => {
    vi.useFakeTimers()
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ matches: mockMatches, predictions: {} }),
    })
    vi.stubGlobal('fetch', fetchMock)

    mountMatches()
    await flushPromises()
    const callsAfterMount = fetchMock.mock.calls.length

    vi.advanceTimersByTime(60_000)
    await flushPromises()

    expect(fetchMock.mock.calls.length).toBe(callsAfterMount)
    vi.useRealTimers()
  })

  it('clears the poll interval on unmount', async () => {
    vi.useFakeTimers()
    const clearSpy = vi.spyOn(globalThis, 'clearInterval')

    const wrapper = mountMatches()
    await flushPromises()
    const callsBefore = clearSpy.mock.calls.length

    wrapper.unmount()

    expect(clearSpy.mock.calls.length).toBeGreaterThan(callsBefore)
    clearSpy.mockRestore()
    vi.useRealTimers()
  })

  it('mounts without error when currentUser is not provided', async () => {
    const wrapper = mount(MatchesView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    expect(wrapper.exists()).toBe(true)
  })

  it('still renders matches when localStorage.getItem throws (Safari Private Browsing)', async () => {
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => { throw new Error('SecurityError') })
    const wrapper = mountMatches()
    await flushPromises()
    expect(wrapper.findAll('[data-testid="match-card"]').length).toBeGreaterThan(0)
    vi.restoreAllMocks()
  })

  it('handles localStorage.setItem throwing when saving guest prediction', async () => {
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => { throw new Error('SecurityError') })
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ matches: mockMatches, predictions: {} }),
    }))
    const wrapper = mountMatches({ currentUser: ref(null) })
    await flushPromises()
    const card = wrapper.find('[data-testid="match-card"]')
    await card.find('input[name="home_goals"]').setValue('2')
    await card.find('input[name="away_goals"]').setValue('1')
    await card.find('button').trigger('click')
    await flushPromises()
    expect(wrapper.exists()).toBe(true)
    vi.restoreAllMocks()
  })

  it('shows loading state before fetch resolves', async () => {
    vi.stubGlobal('fetch', vi.fn().mockReturnValue(new Promise(() => {})))
    const wrapper = mountMatches()
    expect(wrapper.find('[data-testid="loading"]').exists()).toBe(true)
  })

  it('shows error state when fetch fails', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false }))
    const wrapper = mountMatches()
    await flushPromises()
    expect(wrapper.find('[data-testid="error"]').exists()).toBe(true)
  })

  it('shows venue on upcoming match card', async () => {
    const matchWithVenue = { id: 'match-ven-1', homeTeam: 'Columbus Crew', awayTeam: 'FC Dallas', kickoff: futureKickoff(24), status: 'STATUS_SCHEDULED', state: 'pre', homeScore: '', awayScore: '', venue: 'ScottsMiracle-Gro Field' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ matches: [matchWithVenue], predictions: {} }) }))
    const wrapper = mountMatches()
    await flushPromises()
    const card = wrapper.find('[data-testid="match-card"][data-match-id="match-ven-1"]')
    expect(card.find('[data-testid="match-venue"]').text()).toBe('ScottsMiracle-Gro Field')
    vi.restoreAllMocks()
  })

  it('shows venue on live match card', async () => {
    const liveWithVenue = { id: 'match-ven-live', homeTeam: 'Columbus Crew', awayTeam: 'FC Dallas', kickoff: pastKickoff(1), status: 'STATUS_FIRST_HALF', state: 'in', homeScore: '1', awayScore: '0', venue: 'ScottsMiracle-Gro Field' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ matches: [liveWithVenue], predictions: {} }) }))
    const wrapper = mountMatches()
    await flushPromises()
    const card = wrapper.find('[data-testid="now-playing-card"][data-match-id="match-ven-live"]')
    expect(card.find('[data-testid="match-venue"]').text()).toBe('ScottsMiracle-Gro Field')
    vi.restoreAllMocks()
  })

  it('shows venue on result card', async () => {
    const resultWithVenue = { id: 'match-ven-result', homeTeam: 'Columbus Crew', awayTeam: 'FC Dallas', kickoff: pastKickoff(96), status: 'STATUS_FULL_TIME', state: 'post', homeScore: '2', awayScore: '1', venue: 'ScottsMiracle-Gro Field' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ matches: [resultWithVenue], predictions: {} }) }))
    const wrapper = mountMatches()
    await flushPromises()
    const card = wrapper.find('[data-testid="result-card"][data-match-id="match-ven-result"]')
    expect(card.find('[data-testid="match-venue"]').text()).toBe('ScottsMiracle-Gro Field')
    vi.restoreAllMocks()
  })

  it('past-kickoff scheduled match has no Predict button', async () => {
    vi.useFakeTimers()
    const pastKickoff = new Date(Date.now() - 5 * 60 * 1000).toISOString()
    const pastMatch = { id: 'past-scheduled', homeTeam: 'Columbus Crew', awayTeam: 'FC Dallas', kickoff: pastKickoff, status: 'STATUS_SCHEDULED', homeScore: '', awayScore: '', state: 'pre' }
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({ matches: [pastMatch], predictions: {} }) }))
    const wrapper = mountMatches()
    await flushPromises()
    expect(wrapper.find('button[class*="btn-lock"]').exists()).toBe(false)
    vi.useRealTimers()
  })
})
