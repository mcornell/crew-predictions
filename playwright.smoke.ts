import { defineConfig } from '@playwright/test'
import { defineBddConfig } from 'playwright-bdd'

const testDir = defineBddConfig({
  paths: ['e2e/features/*.feature', 'e2e/smoke/features/*.feature'],
  require: ['e2e/steps/*.ts', 'e2e/smoke/steps/*.ts'],
  tags: '@smoke',
})

const debug = process.env.SMOKE_DEBUG === '1'

export default defineConfig({
  testDir,
  timeout: 30000,
  workers: 1,
  reporter: process.env.CI
    ? [['list'], ['html', { outputFolder: 'smoke-report', open: 'never' }], ['junit', { outputFile: 'smoke-results.xml' }]]
    : [['list'], ['html', { outputFolder: 'smoke-report', open: debug ? 'always' : 'never' }]],
  globalSetup: './e2e/smoke/global-setup.ts',
  globalTeardown: './e2e/smoke/global-teardown.ts',
  expect: {
    timeout: 15000,
  },
  use: {
    baseURL: process.env.STAGING_URL ?? 'https://crew-predictions-staging.web.app',
    headless: !debug,
    video: debug ? 'on' : 'off',
    screenshot: 'on',
  },
})
