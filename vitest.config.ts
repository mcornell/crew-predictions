import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  test: {
    include: ['src/**/*.test.ts'],
    environment: 'jsdom',
    reporters: process.env.CI
      ? [
          // GitHub Actions reporter: inline failure annotations on the diff
          // view + auto-generated job summary (test counts, flaky tests with
          // permalinks). Auto-enabled by Vitest only when no reporters are
          // configured, so we list it explicitly alongside our other ones.
          'github-actions',
          // JUnit XML for dorny/test-reporter (separate "Vue Tests" check).
          ['junit', { outputFile: 'vitest-results.xml' }],
          // Default reporter for readable CI logs.
          'default',
        ]
      : ['verbose'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'cobertura'],
      include: ['src/**/*.{ts,vue}'],
      exclude: ['src/main.ts'],
    },
  },
})
