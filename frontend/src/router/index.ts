import {createRouter, createWebHashHistory} from 'vue-router'
import ChatManager from '../views/ChatManager.vue'
import SystemMonitor from '../views/SystemMonitor.vue'
import PromptPilot from '../views/PromptPilot/PromptPilot.vue'
import ModelManager from "../views/ModelManager.vue";
import ModelMarket from "../views/ModelMarket.vue";
import OllamaSettings from "../views/OllamaSettings.vue";
import OllamaApiDebugger from "../views/OllamaApiDebugger.vue";
import OpenAIAdapterSettings from "../views/OpenAIAdapterSettings.vue";

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
    },
    {
        path: '/api-debugger',
        name: 'OllamaApiDebugger',
        component: OllamaApiDebugger
    }
    ,
    {
        path: '/openai-adapter',
        name: 'OpenAIAdapterSettings',
        component: OpenAIAdapterSettings
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router