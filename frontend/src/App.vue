<template>
  <div id="app">
    <el-container style="height: 100vh">
      <el-header style="background-color: #409eff; color: white; padding: 0">
        <div class="header-content">
          <div class="logo">
            <img src="./assets/images/logo-universal.png" alt="Logo" class="logo-img" />
            <span class="logo-text">Ollama 客户端</span>
          </div>
          <div class="header-title">Ollama 本地模型管理工具</div>
        </div>
      </el-header>
      
      <el-container>
        <el-aside width="200px" style="background-color: #f5f7fa">
          <el-menu
            :default-active="activeRoute"
            class="el-menu-vertical"
            router
          >
            <el-menu-item index="/models">
              <el-icon><Collection /></el-icon>
              <span>模型管理</span>
            </el-menu-item>
            <el-menu-item index="/market">
              <el-icon><Shop /></el-icon>
              <span>模型市场</span>
            </el-menu-item>
            <el-menu-item index="/chat">
              <el-icon><ChatDotRound /></el-icon>
              <span>聊天测试</span>
            </el-menu-item>
            <el-menu-item index="/system">
              <el-icon><Monitor /></el-icon>
              <span>系统监控</span>
            </el-menu-item>
            <el-menu-item index="/settings">
              <el-icon><Setting /></el-icon>
              <span>服务设置</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        
        <el-main style="padding: 0">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import {
  Collection,
  Shop,
  ChatDotRound,
  Monitor,
  Setting
} from '@element-plus/icons-vue'

const route = useRoute()
const activeRoute = ref(route.path)

watch(route, (newRoute) => {
  activeRoute.value = newRoute.path
})
</script>

<style scoped>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

.header-content {
  display: flex;
  align-items: center;
  height: 100%;
  padding: 0 20px;
}

.logo {
  display: flex;
  align-items: center;
  margin-right: 30px;
}

.logo-img {
  width: 32px;
  height: 32px;
  margin-right: 10px;
}

.logo-text {
  font-size: 20px;
  font-weight: bold;
}

.header-title {
  font-size: 18px;
}

.el-menu-vertical {
  height: 100%;
  border-right: none;
}

.el-menu-item {
  font-size: 16px;
}

.el-menu-item .el-icon {
  margin-right: 10px;
}
</style>