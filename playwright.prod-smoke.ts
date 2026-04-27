import { defineConfig } from '@playwright/test'
import { defineBddConfig } from 'playwright-bdd'

const testDir = defineBddConfig({
  paths: ['e2e/features/*.feature', 'e2e/smoke/features/health.feature'],
  require: ['e2e/steps/*.ts', 'e2e/smoke/steps/health.steps.ts'],
  tags: '@smoke',
})

export default defineConfig({
  testDir,
  timeout: 30000,
  workers: 1,
  reporter: process.env.CI
    ? [['list'], ['html', { outputFolder: 'prod-smoke-report', open: 'never' }], ['junit', { outputFile: 'prod-smoke-results.xml' }]]
    : [['list'], ['html', { outputFolder: 'prod-smoke-report', open: 'never' }]],
  use: {
    baseURL: process.env.PROD_URL ?? 'https://crew-predictions.web.app',
    headless: true,
    screenshot: 'only-on-failure',
  },
})
