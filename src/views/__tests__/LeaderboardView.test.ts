import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import LeaderboardView from '../LeaderboardView.vue'
import { makeRouter } from '../../test-utils/router'

const mockEntries = [
  { userID: 'firebase:abc', handle: 'BlackAndGold@bsky.mock', acesRadioPoints: 15, upper90ClubPoints: 1, hasProfile: true },
  { userID: 'firebase:def', handle: 'ColumbusNordecke@bsky.mock', acesRadioPoints: 10, upper90ClubPoints: 2, hasProfile: true },
]

const mockData = { entries: mockEntries }

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
    ok: true,
    json: () => Promise.resolve(mockData),
  }))
})

afterEach(() => {
  vi.restoreAllMocks()
})

describe('LeaderboardView', () => {
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

  it('sorts by Upper 90 Club when that column header is clicked', async () => {
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    await wrapper.find('[data-testid="sort-upper90"]').trigger('click')
    const rows = wrapper.findAll('[data-testid="leaderboard-row"]')
    expect(rows[0].find('[data-testid="leaderboard-upper90-points"]').text()).toBe('2')
    expect(rows[1].find('[data-testid="leaderboard-upper90-points"]').text()).toBe('1')
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
})
