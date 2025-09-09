import {createRouter, createWebHashHistory, RouteRecordRaw} from 'vue-router'
import ModelManager from '../views/ModelManager.vue'
import ModelMarket from '../views/ModelMarket.vue'
import SystemMonitor from '../views/SystemMonitor.vue'
import OllamaSettings from '../views/OllamaSettings.vue'
import ChatManager from "../views/ChatManager.vue";

const routes: Array<RouteRecordRaw> = [
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
    name: 'ChatManager',
    component: ChatManager
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