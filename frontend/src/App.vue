<template>
  <div id="app">
    <div class="app-container">
      <!-- Modern Sidebar -->
      <div class="sidebar">
        <div class="sidebar-header">
          <div class="logo">
            <img src="./assets/images/logo.png" alt="Ollama Tools" class="logo-img" />
            <div class="logo-text">
              <h3>Ollama Tools</h3>
              <p>Model Management</p>
            </div>
          </div>
          <div class="language-switcher-container">
            <LanguageSwitcher />
          </div>
        </div>
        
        <nav class="sidebar-nav">
          <router-link to="/models" class="nav-item" :class="{ active: activeRoute === '/models' }">
            <svg class="nav-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M20 6L9 17L4 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
              <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
              <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
            </svg>
            <span>{{ $t('modelManager.title') }}</span>
          </router-link>
          
          <router-link to="/market" class="nav-item" :class="{ active: activeRoute === '/market' }">
            <svg class="nav-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M6 2L3 6V20A2 2 0 0 0 5 22H19A2 2 0 0 0 21 20V6L18 2H6Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M3 6H21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M16 10A4 4 0 0 1 8 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <span>Model Market</span>
          </router-link>
          
          <router-link to="/chat" class="nav-item" :class="{ active: activeRoute === '/chat' }">
            <svg class="nav-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M21 15A2 2 0 0 1 19 17H7L4 20V5A2 2 0 0 1 6 3H19A2 2 0 0 1 21 5Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <span>{{ $t('chatManager.title') }}</span>
          </router-link>
          
          <router-link to="/prompt" class="nav-item" :class="{ active: activeRoute === '/prompt' }">
            <svg class="nav-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
              <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
              <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
            </svg>
            <span>{{ $t('promptEngineering.title') }}</span>
          </router-link>
          
          <router-link to="/settings" class="nav-item" :class="{ active: activeRoute === '/settings' }">
            <svg class="nav-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
              <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
              <path d="M19.4 15A1.65 1.65 0 0 0 21 13.09A1.65 1.65 0 0 0 19.4 9A1.65 1.65 0 0 0 21 6.91A1.65 1.65 0 0 0 19.4 3" stroke="currentColor" stroke-width="2"/>
              <path d="M4.6 9A1.65 1.65 0 0 0 3 10.91A1.65 1.65 0 0 0 4.6 15A1.65 1.65 0 0 0 3 17.09A1.65 1.65 0 0 0 4.6 21" stroke="currentColor" stroke-width="2"/>
            </svg>
            <span>{{ $t('serverSettings.title') }}</span>
          </router-link>
          
          <router-link to="/api-debugger" class="nav-item" :class="{ active: activeRoute === '/api-debugger' }">
            <svg class="nav-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M14.7 6.3A1 1 0 0 0 13 5H5A2 2 0 0 0 3 7V19A2 2 0 0 0 5 21H19A2 2 0 0 0 21 19V11A1 1 0 0 0 19.7 9.3L14.7 6.3Z" stroke="currentColor" stroke-width="2"/>
              <path d="M14 6V10H18" stroke="currentColor" stroke-width="2"/>
              <path d="M16 13H8" stroke="currentColor" stroke-width="2"/>
              <path d="M16 17H8" stroke="currentColor" stroke-width="2"/>
            </svg>
            <span>API Debugger</span>
          </router-link>
          
          <router-link to="/openai-adapter" class="nav-item" :class="{ active: activeRoute === '/openai-adapter' }">
            <svg class="nav-icon" width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
              <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
              <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
            </svg>
            <span>{{ $t('openaiAdapter.title') }}</span>
          </router-link>
        </nav>
      </div>
      
      <!-- Main Content -->
      <div class="main-content">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref, watch} from 'vue'
import {useRoute} from 'vue-router'
import {ChatDotRound, Collection, MagicStick, Setting, Shop, Tools} from '@element-plus/icons-vue' // Added Tools icon
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
}

.sidebar-header {
  padding: 2rem;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.language-switcher-container {
  display: flex;
  justify-content: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 1rem;
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