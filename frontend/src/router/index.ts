import { createRouter, createWebHashHistory } from 'vue-router'
import ChatManager from '../views/ChatManager.vue'
import SystemMonitor from '../views/SystemMonitor.vue'
import PromptPilot from '../views/PromptPilot/PromptPilot.vue'
import ModelManager from "../views/ModelManager.vue";
import ModelMarket from "../views/ModelMarket.vue";
import OllamaSettings from "../views/OllamaSettings.vue";

const routes = [
  {
    path: '/',
    redirect: '/models'
  },
  {
    path: '/models',
    name: 'ModelManager',
    component: ModelManager
  },
  {
    path: '/market',
    name: 'ModelMarket',
    component: ModelMarket
  },
  {
    path: '/chat',
    component: ChatManager
  },
  {
    path: '/prompt',
    component: PromptPilot
  },
  {
    path: '/system',
    name: 'SystemMonitor',
    component: SystemMonitor
  },
  {
    path: '/settings',
    name: 'OllamaSettings',
    component: OllamaSettings
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router