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
  reporter: [['list'], ['html', { open: 'never' }]],
  use: {
    baseURL: 'http://localhost:8080',
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
  webServer: {
    command: 'vite build --logLevel silent && go run ./cmd/server',
    port: 8080,
    reuseExistingServer: !process.env.CI,
    env: {
      FIREBASE_AUTH_EMULATOR_HOST: 'localhost:9099',
      FIREBASE_PROJECT_ID: 'crew-predictions',
      FIREBASE_API_KEY: 'fake-api-key',
      FIREBASE_AUTH_DOMAIN: 'localhost',
      TEST_MODE: '1',
    },
  },
});
