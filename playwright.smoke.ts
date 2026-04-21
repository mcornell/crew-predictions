import { defineConfig } from '@playwright/test'
import { defineBddConfig } from 'playwright-bdd'

const testDir = defineBddConfig({
  paths: ['e2e/smoke/features/*.feature'],
  require: ['e2e/smoke/steps/*.ts'],
})

export default defineConfig({
  testDir,
  timeout: 30000,
  workers: 1,
  reporter: [['list'], ['html', { outputFolder: 'smoke-report', open: 'never' }]],
  globalSetup: './e2e/smoke/global-setup.ts',
  globalTeardown: './e2e/smoke/global-teardown.ts',
  use: {
    baseURL: process.env.STAGING_URL ?? 'https://crew-predictions-staging.web.app',
    headless: true,
  },
})
