<template>
  <div class="model-market">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-icon">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
            <path d="M8 14S9.5 16 12 16S16 14 16 14" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <line x1="9" y1="9" x2="9.01" y2="9" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            <line x1="15" y1="9" x2="15.01" y2="9" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
        </div>
        <div class="header-text">
          <h1>{{ t('modelMarket.title') }}</h1>
          <p>{{ t('modelMarket.description') }}</p>
        </div>
      </div>
    </div>

    <el-card class="market-card">
      <template #header>
        <div class="card-header">
          <span>{{ t('modelMarket.cardTitle') }}</span>
          <div>
            <el-input
              v-model="searchText"
              :placeholder="t('modelMarket.searchPlaceholder')"
              style="width: 300px; margin-right: 10px"
              clearable
              @keyup.enter="searchModels"
            />
            <el-button type="primary" @click="searchModels">{{ t('modelMarket.searchButton') }}</el-button>
          </div>
        </div>
      </template>
      
      <el-table :data="onlineModels" style="width: 100%" v-loading="loading" :empty-text="t('modelMarket.noData')">
        <el-table-column prop="model_name" :label="t('modelMarket.tableName')" />
        <el-table-column prop="description" :label="t('modelMarket.tableDescription')" />
        <el-table-column prop="pulls" :label="t('modelMarket.tableDownloads')">
          <template #default="scope">
            {{ formatPullCount(scope.row.pulls) }}
          </template>
        </el-table-column>
        <el-table-column prop="last_updated" :label="t('modelMarket.tableLastUpdated')" />
        <el-table-column :label="t('modelMarket.tableActions')">
          <template #default="scope">
            <el-button size="small" @click="showModelDetails(scope.row)">{{ t('modelMarket.detailsButton') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- Model Details Drawer -->
    <el-drawer
        v-model="drawerVisible"
        :title="selectedModel?.model_name || t('modelMarket.drawerTitle')"
        direction="rtl"
        size="40%"
        :close-on-click-modal="false"
    >
      <div v-if="selectedModel" class="drawer-content">
        <el-descriptions :column="1" border>
          <el-descriptions-item :label="t('modelMarket.drawerModelName')">{{ selectedModel.model_name }}</el-descriptions-item>
          <el-descriptions-item :label="t('modelMarket.drawerModelId')">{{ selectedModel.model_identifier }}</el-descriptions-item>
          <el-descriptions-item :label="t('modelMarket.drawerDescription')">{{ selectedModel.description || t('modelMarket.drawerNoDescription') }}</el-descriptions-item>
          <el-descriptions-item :label="t('modelMarket.drawerDownloads')">{{ formatPullCount(selectedModel.pulls) }}</el-descriptions-item>
          <el-descriptions-item :label="t('modelMarket.drawerTagsCount')">{{ selectedModel.tags }}</el-descriptions-item>
          <el-descriptions-item :label="t('modelMarket.drawerLastUpdated')">{{ selectedModel.last_updated }}</el-descriptions-item>
          <el-descriptions-item :label="t('modelMarket.drawerUpdateTime')">{{ selectedModel.last_updated_str }}</el-descriptions-item>
        </el-descriptions>
        
        <div style="margin-top: 20px;">
          <el-alert
            :title="t('modelMarket.noticeTitle')"
            type="info"
            :description="t('modelMarket.noticeDescription')"
            show-icon
            :closable="false"
          />
        </div>
        
        <div style="margin-top: 20px;">
          <el-button type="primary" @click="openModelPage(selectedModel.url)">{{ t('modelMarket.viewInNewWindow') }}</el-button>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { SearchOnlineModels } from '../../wailsjs/go/main/App'
import { OpenInBrowser } from '../../wailsjs/go/main/App'
import {OnlineModel} from "../classes/types";

const { t } = useI18n();


const searchText = ref('')
const onlineModels = ref<OnlineModel[]>([])
const loading = ref(false)

// 模型详情抽屉相关变量
const drawerVisible = ref(false)
const selectedModel = ref<OnlineModel | null>(null)

// 格式化下载次数
const formatPullCount = (count: number) => {
  if (count > 1000000) {
    return (count / 1000000).toFixed(1) + 'M'
  }
  if (count > 1000) {
    return (count / 1000).toFixed(1) + 'K'
  }
  return count
}

// 搜索模型
const searchModels = async () => {
  loading.value = true
  try {
    const results = await SearchOnlineModels(searchText.value)
    onlineModels.value = results as unknown as OnlineModel[]
  } catch (error: any) {
    ElMessage.error(t('modelMarket.searchFailed') + ': ' + error.message)
  } finally {
    loading.value = false
  }
}

// 显示模型详情
const showModelDetails = (model: OnlineModel) => {
  selectedModel.value = model
  drawerVisible.value = true
}

// 在浏览器中打开模型页面
const openModelPage = (url: string) => {
  OpenInBrowser(url)
}

// 加载在线模型列表
const loadOnlineModels = async () => {
  loading.value = true
  try {
    const results = await SearchOnlineModels('')
    onlineModels.value = results as unknown as OnlineModel[]
  } catch (error: any) {
    ElMessage.error(t('modelMarket.loadFailed') + ': ' + error.message)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadOnlineModels() // Load online model list when component is mounted
})
</script>

<style scoped>
.model-market {
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 2rem;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.page-header {
  margin-bottom: 2rem;
  flex-shrink: 0;
  position: sticky;
  top: 0;
  z-index: 10;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.header-content {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-icon {
  width: 60px;
  height: 60px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  backdrop-filter: blur(10px);
}

.header-text h1 {
  margin: 0;
  font-size: 2rem;
  font-weight: 700;
  color: white;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-text p {
  margin: 0.5rem 0 0 0;
  color: rgba(255, 255, 255, 0.8);
  font-size: 1.1rem;
}

.market-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.market-card :deep(.el-card__body) {
  height: calc(100% - 60px);
  overflow-y: auto;
}

.drawer-content {
  height: 100%;
}
</style>
