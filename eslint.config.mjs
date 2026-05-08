import js from '@eslint/js'
import tseslint from 'typescript-eslint'
import pluginVue from 'eslint-plugin-vue'
import globals from 'globals'

export default tseslint.config(
  {
    ignores: [
      'dist/**',
      'node_modules/**',
      '.features-gen/**',
      'playwright-report/**',
      'playwright-report*/**',
      'test-results/**',
      'smoke-report/**',
      'prod-smoke-report/**',
      'coverage/**',
      'infra/**',
    ],
  },
  js.configs.recommended,
  ...tseslint.configs.recommended,
  ...pluginVue.configs['flat/recommended'],
  {
    files: ['**/*.vue'],
    languageOptions: {
      parserOptions: {
        parser: tseslint.parser,
      },
      globals: {
        ...globals.browser,
      },
    },
  },
  {
    files: ['**/*.{ts,tsx}'],
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.node,
      },
    },
  },
  {
    files: ['e2e/**/*.ts', 'src/**/__tests__/**/*.ts', 'src/__tests__/**/*.ts'],
    languageOptions: {
      globals: {
        ...globals.node,
      },
    },
    rules: {
      // Test/e2e code legitimately needs `any` for dynamic API response shapes,
      // Cucumber data tables, and `window`-as-global escape hatches.
      '@typescript-eslint/no-explicit-any': 'off',
    },
  },
)
