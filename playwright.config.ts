import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  use: {
    baseURL: 'http://localhost:8080',
  },
  webServer: {
    command: 'go run ./cmd/server',
    port: 8080,
    reuseExistingServer: !process.env.CI,
  },
});
