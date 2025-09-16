<template>
  <div class="model-market">
    <el-card class="market-card">
      <template #header>
        <div class="card-header">
          <span>Model Market</span>
          <div>
            <el-input
              v-model="searchText"
              placeholder="Search models..."
              style="width: 300px; margin-right: 10px"
              clearable
              @keyup.enter="searchModels"
            />
            <el-button type="primary" @click="searchModels">Search</el-button>
          </div>
        </div>
      </template>
      
      <el-table :data="onlineModels" style="width: 100%" v-loading="loading" empty-text="No data">
        <el-table-column prop="model_name" label="Model Name" />
        <el-table-column prop="description" label="Description" />
        <el-table-column prop="pulls" label="Downloads">
          <template #default="scope">
            {{ formatPullCount(scope.row.pulls) }}
          </template>
        </el-table-column>
        <el-table-column prop="last_updated" label="Last Updated" />
        <el-table-column label="Actions">
          <template #default="scope">
            <el-button size="small" @click="showModelDetails(scope.row)">Details</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- Model Details Drawer -->
    <el-drawer
        v-model="drawerVisible"
        :title="selectedModel?.model_name"
        direction="rtl"
        size="40%"
        :close-on-click-modal="false"
    >
      <div v-if="selectedModel" class="drawer-content">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="Model Name">{{ selectedModel.model_name }}</el-descriptions-item>
          <el-descriptions-item label="Model ID">{{ selectedModel.model_identifier }}</el-descriptions-item>
          <el-descriptions-item label="Description">{{ selectedModel.description || 'No description' }}</el-descriptions-item>
          <el-descriptions-item label="Downloads">{{ formatPullCount(selectedModel.pulls) }}</el-descriptions-item>
          <el-descriptions-item label="Tags Count">{{ selectedModel.tags }}</el-descriptions-item>
          <el-descriptions-item label="Last Updated">{{ selectedModel.last_updated }}</el-descriptions-item>
          <el-descriptions-item label="Update Time">{{ selectedModel.last_updated_str }}</el-descriptions-item>
        </el-descriptions>
        
        <div style="margin-top: 20px;">
          <el-alert
            title="Notice"
            type="info"
            description="Due to the target website's security policy, model details cannot be displayed directly on the current page. You can view detailed information in a new window using the button below."
            show-icon
            :closable="false"
          />
        </div>
        
        <div style="margin-top: 20px;">
          <el-button type="primary" @click="openModelPage(selectedModel.url)">View in New Window</el-button>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { SearchOnlineModels } from '../../wailsjs/go/main/App'
import { OpenInBrowser } from '../../wailsjs/go/main/App'
import {OnlineModel} from "../classes/types";


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
    ElMessage.error('Failed to search models: ' + error.message)
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
    ElMessage.error('Failed to load online model list: ' + error.message)
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
  padding: 2rem;
  height: 100%;
  box-sizing: border-box;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.market-card {
  height: 100%;
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