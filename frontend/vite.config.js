import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { configDefaults } from 'vitest/config'

export default defineConfig({
  plugins: [react()],
  test: {
    exclude: [...configDefaults.exclude],
    coverage: {
      reporter: ['json-summary', 'text', 'lcov'], // JSON summary is needed
    },
  },
})
