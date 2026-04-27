import { createBdd } from 'playwright-bdd'
import { expect } from '@playwright/test'

const { Then } = createBdd()

Then('the API at {string} returns a JSON object with key {string}', async ({ request }, path: string, key: string) => {
  const response = await request.get(path)
  expect(response.ok()).toBe(true)
  const body = await response.json()
  expect(typeof body).toBe('object')
  expect(body).toHaveProperty(key)
})
