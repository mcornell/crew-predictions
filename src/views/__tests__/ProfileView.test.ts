import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { ref } from 'vue'
import ProfileView from '../ProfileView.vue'
import { makeRouter } from '../../test-utils/router'

const profileData = {
  userID: 'firebase:abc',
  handle: 'CrewFan',
  location: 'Columbus, OH',
  predictionCount: 3,
  acesRadio: { points: 15, rank: 1 },
  upper90Club: { points: 2, rank: 2 },
  grouchy: { points: 1, rank: 1 },
}

async function mountProfile(userID: string, currentUserID: string | null = null, data = profileData) {
  vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve(data) }))
  const r = makeRouter()
  r.addRoute({ path: '/profile/:userID', component: ProfileView })
  await r.push(`/profile/${userID}`)
  const currentUser = ref(currentUserID ? { userID: currentUserID, handle: 'me', emailVerified: true } : null)
  return mount(ProfileView, { global: { plugins: [r], provide: { currentUser } } })
}

describe('ProfileView', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('fetches and displays handle and location', async () => {
    const wrapper = await mountProfile('firebase:abc')
    await flushPromises()
    expect(wrapper.text()).toContain('CrewFan')
    expect(wrapper.text()).toContain('Columbus, OH')
  })

  it('shows prediction count and Aces Radio points', async () => {
    const wrapper = await mountProfile('firebase:abc')
    await flushPromises()
    expect(wrapper.text()).toContain('3')
    expect(wrapper.text()).toContain('15')
  })

  it('shows edit form when viewing own profile', async () => {
    const wrapper = await mountProfile('firebase:abc', 'firebase:abc')
    await flushPromises()
    expect(wrapper.find('form[data-testid="profile-form"]').exists()).toBe(true)
  })

  it('hides edit form when viewing another user profile', async () => {
    const wrapper = await mountProfile('firebase:abc', 'firebase:other')
    await flushPromises()
    expect(wrapper.find('form[data-testid="profile-form"]').exists()).toBe(false)
  })

  it('pre-populates handle and location inputs with current values', async () => {
    const wrapper = await mountProfile('firebase:abc', 'firebase:abc')
    await flushPromises()
    expect((wrapper.find('input[data-testid="display-name-input"]').element as HTMLInputElement).value).toBe('CrewFan')
    expect((wrapper.find('input[data-testid="location-input"]').element as HTMLInputElement).value).toBe('Columbus, OH')
  })

  it('posts handle and location on save, then navigates to /matches', async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve(profileData) })
      .mockResolvedValueOnce({ ok: true })
    vi.stubGlobal('fetch', fetchMock)

    const r = makeRouter()
    r.addRoute({ path: '/profile/:userID', component: ProfileView })
    await r.push('/profile/firebase:abc')
    const currentUser = ref({ userID: 'firebase:abc', handle: 'CrewFan', emailVerified: true })
    const wrapper = mount(ProfileView, { global: { plugins: [r], provide: { currentUser } } })
    await flushPromises()

    await wrapper.find('input[data-testid="display-name-input"]').setValue('NewHandle')
    await wrapper.find('input[data-testid="location-input"]').setValue('Parts Unknown')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    const handleCall = fetchMock.mock.calls.find(c => c[0] === '/auth/handle')
    expect(handleCall).toBeTruthy()
    const body = handleCall![1].body as URLSearchParams
    expect(body.get('handle')).toBe('NewHandle')
    expect(body.get('location')).toBe('Parts Unknown')
    expect(r.currentRoute.value.path).toBe('/matches')
  })

  it('shows loading state before fetch resolves', async () => {
    vi.stubGlobal('fetch', vi.fn().mockReturnValue(new Promise(() => {})))
    const r = makeRouter()
    r.addRoute({ path: '/profile/:userID', component: ProfileView })
    await r.push('/profile/firebase:abc')
    const wrapper = mount(ProfileView, { global: { plugins: [r], provide: { currentUser: ref(null) } } })
    expect(wrapper.find('[data-testid="loading"]').exists()).toBe(true)
  })

  it('shows error state when fetch fails', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: false }))
    const r = makeRouter()
    r.addRoute({ path: '/profile/:userID', component: ProfileView })
    await r.push('/profile/firebase:abc')
    const wrapper = mount(ProfileView, { global: { plugins: [r], provide: { currentUser: ref(null) } } })
    await flushPromises()
    expect(wrapper.find('[data-testid="error"]').exists()).toBe(true)
  })

  it('shows Grouchy points and rank on profile', async () => {
    const wrapper = await mountProfile('firebase:abc')
    await flushPromises()
    expect(wrapper.find('[data-testid="grouchy-points"]').text()).toBe('1')
    expect(wrapper.find('[data-testid="grouchy-rank"]').text()).toBe('#1')
  })

  it('hides Grouchy rank when not yet ranked', async () => {
    const wrapper = await mountProfile('firebase:abc', null, { ...profileData, grouchy: { points: 0, rank: 0 } })
    await flushPromises()
    expect(wrapper.find('[data-testid="grouchy-rank"]').exists()).toBe(false)
  })

  it('shows error when save fails', async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve(profileData) })
      .mockResolvedValueOnce({ ok: false })
    vi.stubGlobal('fetch', fetchMock)
    const r = makeRouter()
    r.addRoute({ path: '/profile/:userID', component: ProfileView })
    await r.push('/profile/firebase:abc')
    const currentUser = ref({ userID: 'firebase:abc', handle: 'CrewFan', emailVerified: true })
    const wrapper = mount(ProfileView, { global: { plugins: [r], provide: { currentUser } } })
    await flushPromises()

    await wrapper.find('form').trigger('submit')
    await flushPromises()
    expect(wrapper.find('.form-error').exists()).toBe(true)
  })
})
