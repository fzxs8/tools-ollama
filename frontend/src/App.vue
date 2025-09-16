<template>
  <div id="app">
    <div class="app-container">
      <!-- Modern Sidebar -->
      <div class="sidebar">
        <div class="sidebar-header">
          <div class="logo">
            <img src="./assets/images/logo.png" alt="Ollama Tools" class="logo-img"/>
            <div class="logo-text">
              <h3>Ollama Tools</h3>
              <p>Model Management</p>
            </div>
          </div>
        </div>

        <nav class="sidebar-nav">
          <router-link to="/models" class="nav-item" :class="{ active: activeRoute === '/models' }">
            <svg class="nav-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M4 7V4C4 2.9 4.9 2 6 2H18C19.1 2 20 2.9 20 4V7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <rect x="2" y="7" width="20" height="15" rx="2" ry="2" stroke="currentColor" stroke-width="2"/>
              <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
            </svg>
            <span>{{ $t('modelManager.title') }}</span>
          </router-link>

          <router-link to="/market" class="nav-item" :class="{ active: activeRoute === '/market' }">
            <svg class="nav-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M6 2L3 6V20A2 2 0 0 0 5 22H19A2 2 0 0 0 21 20V6L18 2H6Z" stroke="currentColor" stroke-width="2"
                    stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M3 6H21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M16 10A4 4 0 0 1 8 10" stroke="currentColor" stroke-width="2" stroke-linecap="round"
                    stroke-linejoin="round"/>
            </svg>
            <span>{{ $t('modelMarket.title') }}</span>
          </router-link>

          <router-link to="/chat" class="nav-item" :class="{ active: activeRoute === '/chat' }">
            <svg class="nav-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M21 15A2 2 0 0 1 19 17H7L4 20V5A2 2 0 0 1 6 3H19A2 2 0 0 1 21 5Z" stroke="currentColor"
                    stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <span>{{ $t('chatManager.title') }}</span>
          </router-link>

          <router-link to="/prompt" class="nav-item" :class="{ active: activeRoute === '/prompt' }">
            <svg class="nav-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M14 2H6C4.9 2 4 2.9 4 4V20C4 21.1 4.9 22 6 22H18C19.1 22 20 21.1 20 20V8L14 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M14 2V8H20" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M16 13H8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M16 17H8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <span>{{ $t('promptEngineering.title') }}</span>
          </router-link>
          <router-link to="/api-debugger" class="nav-item" :class="{ active: activeRoute === '/api-debugger' }">
            <svg class="nav-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M12 15L8 11L16 11L12 15Z" fill="currentColor"/>
              <path d="M21 3H3C1.9 3 1 3.9 1 5V19C1 20.1 1.9 21 3 21H21C22.1 21 23 20.1 23 19V5C23 3.9 22.1 3 21 3ZM21 19H3V5H21V19Z" fill="currentColor"/>
              <path d="M5 7H19V9H5V7Z" fill="currentColor"/>
              <path d="M5 11H11V13H5V11Z" fill="currentColor"/>
            </svg>
            <span>{{ $t('apiDebugger.title') }}</span>
          </router-link>

          <router-link to="/openai-adapter" class="nav-item" :class="{ active: activeRoute === '/openai-adapter' }">
            <svg class="nav-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M12 2L13.09 8.26L22 9L13.09 9.74L12 16L10.91 9.74L2 9L10.91 8.26L12 2Z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
              <path d="M12 12C14.21 12 16 10.21 16 8S14.21 4 12 4S8 5.79 8 8S9.79 12 12 12Z" stroke="currentColor" stroke-width="2"/>
              <path d="M8 21V19C8 16.79 9.79 15 12 15S16 16.79 16 19V21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <span>{{ $t('openaiAdapter.title') }}</span>
          </router-link>

          <router-link to="/settings" class="nav-item" :class="{ active: activeRoute === '/settings' }">
            <svg class="nav-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M12 20H21L19 16H14L12 20Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M12 20V4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M6 8L12 4L18 8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M6 16L12 20L18 16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <circle cx="12" cy="12" r="2" stroke="currentColor" stroke-width="2"/>
            </svg>
            <span>{{ $t('serverSettings.title') }}</span>
          </router-link>
        </nav>
        
        <div class="language-switcher-container">
          <LanguageSwitcher/>
          <a href="https://github.com/fzxs8/tools-ollama" target="_blank" class="github-link">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M12 0C5.374 0 0 5.373 0 12 0 17.302 3.438 21.8 8.207 23.387c.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z" fill="currentColor"/>
            </svg>
          </a>
        </div>
      </div>

      <!-- Main Content -->
      <div class="main-content">
        <router-view/>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref, watch} from 'vue'
import {useRoute} from 'vue-router'
import LanguageSwitcher from './components/LanguageSwitcher.vue'

const route = useRoute()
const activeRoute = ref(route.path)

watch(route, (newRoute) => {
  activeRoute.value = newRoute.path
})
</script>

<style scoped>
#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  height: 100vh;
  overflow: hidden;
}

.app-container {
  display: flex;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.sidebar {
  width: 280px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-right: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  flex-direction: column;
  box-shadow: 4px 0 20px rgba(0, 0, 0, 0.1);
  justify-content: space-between;
}

.sidebar-header {
  padding: 2rem;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  flex-shrink: 0;
}

.logo {
  display: flex;
  align-items: center;
  gap: 1rem;
  white-space: nowrap;
}

.logo-img {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.logo-text h3 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 700;
  color: #1a202c;
  line-height: 1.2;
}

.logo-text p {
  margin: 0;
  font-size: 0.875rem;
  color: #718096;
  font-weight: 500;
}

.sidebar-nav {
  flex: 1;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  overflow-y: auto;
}

.language-switcher-container {
  padding: 1rem 2rem 2rem 2rem;
  border-top: 1px solid rgba(0, 0, 0, 0.05);
  flex-shrink: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
}

.github-link {
  color: #4a5568;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem;
  border-radius: 8px;
}

.github-link:hover {
  color: #667eea;
  background: rgba(102, 126, 234, 0.1);
  transform: translateY(-1px);
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1rem;
  border-radius: 12px;
  text-decoration: none;
  color: #4a5568;
  font-weight: 500;
  font-size: 0.95rem;
  transition: all 0.3s ease;
  position: relative;
}

.nav-item:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  transform: translateX(4px);
}

.nav-item.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.nav-item.active:hover {
  transform: translateX(0);
}

.nav-icon {
  flex-shrink: 0;
  opacity: 0.8;
}

.nav-item.active .nav-icon {
  opacity: 1;
}

.main-content {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

@media (max-width: 768px) {
  .sidebar {
    width: 240px;
  }

  .sidebar-header {
    padding: 1.5rem;
  }

  .logo-text h3 {
    font-size: 1.1rem;
  }

  .nav-item {
    padding: 0.75rem;
    font-size: 0.9rem;
  }
}
</style>