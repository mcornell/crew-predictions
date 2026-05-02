import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import LeaderboardView from '../LeaderboardView.vue'
import { makeRouter } from '../../test-utils/router'

const mockEntries = [
  { userID: 'firebase:abc', handle: 'BlackAndGold@bsky.mock', acesRadioPoints: 15, upper90ClubPoints: 1, grouchyPoints: 3, hasProfile: true },
  { userID: 'firebase:def', handle: 'ColumbusNordecke@bsky.mock', acesRadioPoints: 10, upper90ClubPoints: 2, grouchyPoints: 1, hasProfile: true },
]

const mockData = { entries: mockEntries }

// Default seasons response: a current season and a past one. Tests that need
// different season state stub fetch themselves before mounting.
const defaultSeasons = [
  { id: '2025', name: '2025 Season', isCurrent: false },
  { id: '2026', name: '2026 Season', isCurrent: true },
]

function urlRouter(overrides: Record<string, unknown> = {}) {
  return vi.fn(async (url: string) => {
    if (url in overrides) return overrides[url] as Response
    if (url.startsWith('/api/leaderboard')) return { ok: true, json: () => Promise.resolve(mockData) } as unknown as Response
    if (url === '/api/seasons') return { ok: true, json: () => Promise.resolve({ seasons: defaultSeasons }) } as unknown as Response
    return { ok: false, status: 404 } as Response
  })
}

beforeEach(() => {
  vi.stubGlobal('fetch', urlRouter())
})
// Note: individual tests that need distinct fetch responses stub fetch themselves

afterEach(() => {
  vi.restoreAllMocks()
})

function makeLeaderboardRouter(path = '/leaderboard') {
  const r = makeRouter()
  r.addRoute({ path: '/leaderboard/:season', component: LeaderboardView })
  return r
}

describe('LeaderboardView', () => {
  it('fetches from /api/leaderboard/:season when season route param is set', async () => {
    const fetchMock = urlRouter({
      '/api/leaderboard/2026': { ok: true, json: () => Promise.resolve({ entries: [
        { handle: 'HistoryFan', acesRadioPoints: 15, upper90ClubPoints: 3, grouchyPoints: 1 }
      ]}) } as unknown as Response,
    })
    vi.stubGlobal('fetch', fetchMock)
    const r = makeLeaderboardRouter()
    await r.push('/leaderboard/2026')
    const wrapper = mount(LeaderboardView, { global: { plugins: [r] } })
    await flushPromises()
    expect(fetchMock).toHaveBeenCalledWith('/api/leaderboard/2026')
    const rows = wrapper.findAll('[data-testid="leaderboard-row"]')
    expect(rows[0].find('[data-testid="leaderboard-aces-points"]').text()).toBe('15')
  })

  it('sets document title to Leaderboard — Crew Predictions', async () => {
    mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    expect(document.title).toBe('Leaderboard — Crew Predictions')
  })

  it('shows a helpful message when no predictions scored yet', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ entries: [] }),
    }))
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    expect(wrapper.text()).toContain('No predictions scored yet')
  })

  it('renders one row per predictor in the unified table', async () => {
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    const rows = wrapper.findAll('[data-testid="leaderboard-row"]')
    expect(rows).toHaveLength(2)
  })

  it('shows Aces Radio points in each row', async () => {
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    const rows = wrapper.findAll('[data-testid="leaderboard-row"]')
    expect(rows[0].find('[data-testid="leaderboard-aces-points"]').text()).toBe('15')
    expect(rows[1].find('[data-testid="leaderboard-aces-points"]').text()).toBe('10')
  })

  it('shows Upper 90 Club points in each row', async () => {
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    const rows = wrapper.findAll('[data-testid="leaderboard-row"]')
    expect(rows[0].find('[data-testid="leaderboard-upper90-points"]').text()).toBe('1')
    expect(rows[1].find('[data-testid="leaderboard-upper90-points"]').text()).toBe('2')
  })

  it.each([
    ['sort-upper90',  'leaderboard-upper90-points',  '2'],
    ['sort-grouchy',  'leaderboard-grouchy-points',  '3'],
  ])('sorts by %s when that column header is clicked', async (sortTestid, pointsTestid, expectedTopValue) => {
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    await wrapper.find(`[data-testid="${sortTestid}"]`).trigger('click')
    const rows = wrapper.findAll('[data-testid="leaderboard-row"]')
    expect(rows[0].find(`[data-testid="${pointsTestid}"]`).text()).toBe(expectedTopValue)
  })

  it('sorts by Aces Radio when that column header is clicked', async () => {
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    await wrapper.find('[data-testid="sort-upper90"]').trigger('click')
    await wrapper.find('[data-testid="sort-aces"]').trigger('click')
    const rows = wrapper.findAll('[data-testid="leaderboard-row"]')
    expect(rows[0].find('[data-testid="leaderboard-aces-points"]').text()).toBe('15')
  })

  it('rank is dynamic — ties share rank and next rank skips', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        entries: [
          { userID: 'u1', handle: 'Fan1', acesRadioPoints: 10, upper90ClubPoints: 2, hasProfile: true },
          { userID: 'u2', handle: 'Fan2', acesRadioPoints: 10, upper90ClubPoints: 1, hasProfile: true },
          { userID: 'u3', handle: 'Fan3', acesRadioPoints: 5, upper90ClubPoints: 1, hasProfile: true },
        ],
      }),
    }))
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    const rows = wrapper.findAll('[data-testid="leaderboard-row"]')
    expect(rows[0].find('[data-testid="leaderboard-rank"]').text()).toBe('1')
    expect(rows[1].find('[data-testid="leaderboard-rank"]').text()).toBe('1')
    expect(rows[2].find('[data-testid="leaderboard-rank"]').text()).toBe('3')
  })

  it('handle links to profile page when hasProfile is true', async () => {
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    const link = wrapper.find('[data-testid="leaderboard-row"] a')
    expect(link.attributes('href')).toContain('firebase:abc')
  })

  it('renders plain span instead of link when hasProfile is false', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        entries: [{ userID: 'legacyfan', handle: 'legacyfan', acesRadioPoints: 5, upper90ClubPoints: 1, hasProfile: false }],
      }),
    }))
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    const row = wrapper.find('[data-testid="leaderboard-row"]')
    expect(row.find('a').exists()).toBe(false)
    expect(row.find('[data-testid="leaderboard-handle"]').text()).toBe('legacyfan')
  })

  it('shows loading state before fetch resolves', async () => {
    vi.stubGlobal('fetch', vi.fn().mockReturnValue(new Promise(() => {})))
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    expect(wrapper.find('[data-testid="loading"]').exists()).toBe(true)
  })

  it('shows error state when fetch fails', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false }))
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="error"]').exists()).toBe(true)
  })

  it('renders mobile sort buttons for both scoring formats', async () => {
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="mobile-sort-aces"]').exists()).toBe(true)
    expect(wrapper.find('[data-testid="mobile-sort-upper90"]').exists()).toBe(true)
  })

  it.each([
    ['mobile-sort-upper90',  'leaderboard-upper90-points',  '2'],
    ['mobile-sort-grouchy',  'leaderboard-grouchy-points',  '3'],
  ])('mobile %s button triggers sort change', async (sortTestid, pointsTestid, expectedTopValue) => {
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    await wrapper.find(`[data-testid="${sortTestid}"]`).trigger('click')
    const rows = wrapper.findAll('[data-testid="leaderboard-row"]')
    expect(rows[0].find(`[data-testid="${pointsTestid}"]`).text()).toBe(expectedTopValue)
  })

  it('shows Grouchy points in each row', async () => {
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    const rows = wrapper.findAll('[data-testid="leaderboard-row"]')
    expect(rows[0].find('[data-testid="leaderboard-grouchy-points"]').text()).toBe('3')
    expect(rows[1].find('[data-testid="leaderboard-grouchy-points"]').text()).toBe('1')
  })

  it('renders mobile sort button for Grouchy', async () => {
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    expect(wrapper.find('[data-testid="mobile-sort-grouchy"]').exists()).toBe(true)
  })

  it('mobile sort Aces button switches back from another sort', async () => {
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    await wrapper.find('[data-testid="mobile-sort-upper90"]').trigger('click')
    await wrapper.find('[data-testid="mobile-sort-aces"]').trigger('click')
    const rows = wrapper.findAll('[data-testid="leaderboard-row"]')
    expect(rows[0].find('[data-testid="leaderboard-aces-points"]').text()).toBe('15')
  })

  describe('season selector', () => {
    it('renders selector showing current season name when past seasons exist', async () => {
      const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
      await flushPromises()
      const selector = wrapper.find('[data-testid="season-selector"]')
      expect(selector.exists()).toBe(true)
      expect(selector.text()).toContain('2026 Season')
    })

    it('hides selector when there are no past seasons', async () => {
      vi.stubGlobal('fetch', vi.fn(async (url: string) => {
        if (url.startsWith('/api/leaderboard')) return { ok: true, json: () => Promise.resolve(mockData) } as unknown as Response
        if (url === '/api/seasons') return { ok: true, json: () => Promise.resolve({ seasons: [{ id: '2026', name: '2026 Season', isCurrent: true }] }) } as unknown as Response
        return { ok: false } as Response
      }))
      const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
      await flushPromises()
      expect(wrapper.find('[data-testid="season-selector"]').exists()).toBe(false)
    })

    it('clicking selector opens a flyout listing past seasons', async () => {
      const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
      await flushPromises()
      await wrapper.find('[data-testid="season-selector"]').trigger('click')
      const flyout = wrapper.find('[data-testid="season-flyout"]')
      expect(flyout.exists()).toBe(true)
      expect(flyout.text()).toContain('2025 Season')
    })

    it('flyout includes a Current Season link to /leaderboard', async () => {
      const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
      await flushPromises()
      await wrapper.find('[data-testid="season-selector"]').trigger('click')
      const currentLink = wrapper.find('[data-testid="season-flyout"] a[href="/leaderboard"]')
      expect(currentLink.exists()).toBe(true)
    })

    it('past-season links point to /leaderboard/:id', async () => {
      const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
      await flushPromises()
      await wrapper.find('[data-testid="season-selector"]').trigger('click')
      const pastLink = wrapper.find('[data-testid="season-flyout"] a[href="/leaderboard/2025"]')
      expect(pastLink.exists()).toBe(true)
    })

    it('shows the historical season name when viewing /leaderboard/:season', async () => {
      const r = makeLeaderboardRouter()
      await r.push('/leaderboard/2025')
      const wrapper = mount(LeaderboardView, { global: { plugins: [r] } })
      await flushPromises()
      const selector = wrapper.find('[data-testid="season-selector"]')
      expect(selector.text()).toContain('2025 Season')
    })

    it('clicking a past-season link closes the flyout', async () => {
      const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
      await flushPromises()
      await wrapper.find('[data-testid="season-selector"]').trigger('click')
      await wrapper.find('[data-testid="season-flyout"] a').trigger('click')
      expect(wrapper.find('[data-testid="season-flyout"]').exists()).toBe(false)
    })
  })
})
