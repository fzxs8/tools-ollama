import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  build: {
    rollupOptions: {
      external: ['../wailsjs/go/main/App', '../wailsjs/runtime'],
      output: {
        globals: {
          vue: 'Vue'
        }
      }
    }
  },
  resolve: {
    conditions: ['module', 'browser', 'development|production']
  },
  ssr: {
    resolve: {
      conditions: ['module', 'node', 'development|production']
    }
  }
})