import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import LeaderboardView from '../LeaderboardView.vue'

const mockData = {
  acesRadio: [{ handle: 'BlackAndGold@bsky.mock', points: 15 }],
  upper90Club: [{ handle: 'ColumbusNordecke@bsky.mock', points: 2 }],
}

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
    ok: true,
    json: () => Promise.resolve(mockData),
  }))
})

describe('LeaderboardView', () => {
  it('sets document title to Leaderboard — Crew Predictions', async () => {
    mount(LeaderboardView)
    await flushPromises()
    expect(document.title).toBe('Leaderboard — Crew Predictions')
  })

  it('shows a helpful message when no predictions scored yet', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ acesRadio: [], upper90Club: [] }),
    }))
    const wrapper = mount(LeaderboardView)
    await flushPromises()
    expect(wrapper.text()).toContain('No predictions scored yet')
    expect(wrapper.text()).toContain('rules')
  })

  it('renders leaderboard rows for Aces Radio', async () => {
    const wrapper = mount(LeaderboardView)
    await flushPromises()
    const rows = wrapper.findAll('[data-testid="leaderboard-row"]')
    expect(rows.length).toBeGreaterThan(0)
    expect(rows[0].text()).toContain('BlackAndGold@bsky.mock')
    expect(rows[0].find('[data-testid="leaderboard-points"]').text()).toBe('15')
  })

  it('renders leaderboard rows for Upper90Club', async () => {
    const wrapper = mount(LeaderboardView)
    await flushPromises()
    const rows = wrapper.findAll('[data-testid="leaderboard-row"]')
    const u90Row = rows[rows.length - 1]
    expect(u90Row.text()).toContain('ColumbusNordecke@bsky.mock')
    expect(u90Row.find('[data-testid="leaderboard-points"]').text()).toBe('2')
  })
})
