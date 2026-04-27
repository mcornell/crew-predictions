import { defineConfig } from '@playwright/test';
import { defineBddConfig } from 'playwright-bdd';

const testDir = defineBddConfig({
  features: 'e2e/features/**/*.feature',
  steps: 'e2e/steps/**/*.steps.ts',
});

export default defineConfig({
  testDir,
  workers: 2,
  globalSetup: './e2e/global-setup.ts',
  reporter: process.env.CI
    ? [['list'], ['github'], ['html', { open: 'never' }], ['junit', { outputFile: 'playwright-results.xml' }]]
    : [['list'], ['html', { open: 'never' }]],
  use: {
    baseURL: 'http://localhost:8083',
    screenshot: 'only-on-failure',
    trace: 'retain-on-failure',
  },
  projects: [
    {
      name: 'auth',
      grepInvert: /@reset/,
    },
    {
      name: 'app',
      grep: /@reset/,
      workers: 1,
    },
  ],
  webServer: [
    {
      command: 'go run ./cmd/server',
      port: 8082,
      reuseExistingServer: !process.env.CI,
      env: {
        PORT: '8082',
        FIREBASE_AUTH_EMULATOR_HOST: 'localhost:9099',
        FIREBASE_PROJECT_ID: 'crew-predictions',
        FIREBASE_API_KEY: 'fake-api-key',
        FIREBASE_AUTH_DOMAIN: 'localhost',
        TEST_MODE: '1',
      },
    },
    {
      command: 'vite build --logLevel silent && vite preview --port 8083',
      port: 8083,
      reuseExistingServer: !process.env.CI,
      timeout: 120000,
    },
  ],
});
