import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  test: {
    include: ['src/**/*.test.ts'],
    environment: 'jsdom',
    reporters: process.env.CI
      ? [['junit', { outputFile: 'vitest-results.xml' }], 'verbose']
      : ['verbose'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'cobertura'],
      include: ['src/**/*.{ts,vue}'],
      exclude: ['src/main.ts'],
    },
  },
})
