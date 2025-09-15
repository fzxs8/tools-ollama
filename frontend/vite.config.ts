import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    watch: {
      usePolling: true,
      interval: 1000,
      ignored: [
        '**/wailsjs/**', 
        '**/*.db*', 
        '**/*.log',
        '**/data/**',
        '**/logs/**',
        '**/tmp/**',
        '**/temp/**',
        '**/node_modules/**',
        '**/dist/**',
        '**/build/**',
        '**/.git/**',
        '**/package-lock.json',
        '**/yarn.lock',
        '**/package.json.md5'
      ]
    }
  },
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