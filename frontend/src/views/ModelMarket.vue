<template>
  <div class="model-market">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>模型市场</span>
          <div>
            <el-input
              v-model="searchText"
              placeholder="搜索模型..."
              style="width: 300px; margin-right: 10px"
              clearable
              @keyup.enter="searchModels"
            />
            <el-button type="primary" @click="searchModels">搜索</el-button>
          </div>
        </div>
      </template>
      
      <el-table :data="onlineModels" style="width: 100%" v-loading="loading" empty-text="暂无数据">
        <el-table-column prop="model_name" label="模型名称" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="pulls" label="下载次数">
          <template #default="scope">
            {{ formatPullCount(scope.row.pulls) }}
          </template>
        </el-table-column>
        <el-table-column prop="last_updated" label="更新时间" />
        <el-table-column label="操作">
          <template #default="scope">
            <el-button size="small" @click="showModelDetails(scope.row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 模型详情抽屉 -->
    <el-drawer
        v-model="drawerVisible"
        :title="selectedModel?.model_name"
        direction="rtl"
        size="40%"
    >
      <div v-if="selectedModel" class="drawer-content">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="模型名称">{{ selectedModel.model_name }}</el-descriptions-item>
          <el-descriptions-item label="模型标识">{{ selectedModel.model_identifier }}</el-descriptions-item>
          <el-descriptions-item label="描述">{{ selectedModel.description || '暂无描述' }}</el-descriptions-item>
          <el-descriptions-item label="下载次数">{{ formatPullCount(selectedModel.pulls) }}</el-descriptions-item>
          <el-descriptions-item label="标签数量">{{ selectedModel.tags }}</el-descriptions-item>
          <el-descriptions-item label="最后更新">{{ selectedModel.last_updated }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ selectedModel.last_updated_str }}</el-descriptions-item>
        </el-descriptions>
        
        <div style="margin-top: 20px;">
          <el-alert
            title="提示"
            type="info"
            description="由于目标网站的安全策略，无法在当前页面直接显示模型详情。您可以通过下面的按钮在新窗口中查看详细信息。"
            show-icon
            :closable="false"
          />
        </div>
        
        <div style="margin-top: 20px;">
          <el-button type="primary" @click="openModelPage(selectedModel.url)">在新窗口中查看</el-button>
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

interface OnlineModel {
  model_identifier: string
  namespace: string | null
  model_name: string
  model_type: string
  description: string
  capability: string | null
  labels: string[]
  pulls: number
  tags: number
  last_updated: string
  last_updated_str: string
  url: string
}

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
    ElMessage.error('搜索模型失败: ' + error.message)
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
    ElMessage.error('加载在线模型列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadOnlineModels() // 组件挂载时加载在线模型列表
})
</script>

<style scoped>
.model-market {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.drawer-content {
  height: 100%;
}
</style>