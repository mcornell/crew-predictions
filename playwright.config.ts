import { defineConfig } from '@playwright/test';
import { defineBddConfig } from 'playwright-bdd';

const testDir = defineBddConfig({
  features: 'e2e/features/**/*.feature',
  steps: 'e2e/steps/**/*.steps.ts',
});

export default defineConfig({
  testDir,
  use: {
    baseURL: 'http://localhost:8080',
  },
  webServer: {
    command: 'go run ./cmd/server',
    port: 8080,
    reuseExistingServer: !process.env.CI,
  },
});
