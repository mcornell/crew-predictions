import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import LeaderboardView from '../LeaderboardView.vue'
import { makeRouter } from '../../test-utils/router'

const mockData = {
  acesRadio: [{ userID: 'firebase:abc', handle: 'BlackAndGold@bsky.mock', points: 15 }],
  upper90Club: [{ userID: 'firebase:def', handle: 'ColumbusNordecke@bsky.mock', points: 2 }],
}

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
    ok: true,
    json: () => Promise.resolve(mockData),
  }))
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
      json: () => Promise.resolve({ acesRadio: [], upper90Club: [] }),
    }))
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    expect(wrapper.text()).toContain('No predictions scored yet')
  })

  it('renders leaderboard rows for Aces Radio', async () => {
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    const rows = wrapper.findAll('[data-testid="leaderboard-row"]')
    expect(rows.length).toBeGreaterThan(0)
    expect(rows[0].text()).toContain('BlackAndGold@bsky.mock')
    expect(rows[0].find('[data-testid="leaderboard-points"]').text()).toBe('15')
  })

  it('handle links to profile page', async () => {
    const wrapper = mount(LeaderboardView, { global: { plugins: [makeRouter()] } })
    await flushPromises()
    const link = wrapper.find('[data-testid="leaderboard-row"] a')
    expect(link.attributes('href')).toContain('firebase:abc')
  })
})
