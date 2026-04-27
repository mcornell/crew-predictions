import { createBdd } from 'playwright-bdd'
import { expect } from '@playwright/test'

const { Then } = createBdd()

Then('the API at {string} returns a JSON array', async ({ request }, path: string) => {
  const response = await request.get(path)
  expect(response.ok()).toBe(true)
  const body = await response.json()
  expect(Array.isArray(body)).toBe(true)
})
