
import './style.css';
import './app.css';

import {createApp} from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import i18n from './i18n'
import { useI18nStore } from './stores/i18n'

const app = createApp(App)
const pinia = createPinia()

// 注册所有Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

app.use(pinia) // 注册 Pinia
app.use(router)
app.use(ElementPlus)
app.use(i18n)

// Initialize i18n store to ensure proper locale sync
const i18nStore = useI18nStore()

app.mount('#app')