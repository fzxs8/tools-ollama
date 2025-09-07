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
            />
            <el-button type="primary" @click="searchModels">搜索</el-button>
            <el-button type="success" @click="openSearchDrawer" style="margin-left: 10px">搜索在线模型</el-button>
          </div>
        </div>
      </template>
      
      <el-table :data="marketModels" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="模型名称" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="size" label="大小">
          <template #default="scope">
            {{ formatSize(scope.row.size) }}
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button size="small" @click="showModelDetails(scope.row)">详情</el-button>
            <el-button size="small" type="primary" @click="downloadModel(scope.row)">下载</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 在线模型展示区域 -->
    <el-card style="margin-top: 20px">
      <template #header>
        <div class="card-header">
          <span>热门在线模型</span>
        </div>
      </template>
      
      <el-table :data="onlineModels" style="width: 100%" v-loading="loadingOnlineModels">
        <el-table-column prop="name" label="模型名称" />
        <el-table-column prop="pull_count" label="下载次数">
          <template #default="scope">
            {{ formatPullCount(scope.row.pull_count) }}
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间">
          <template #default="scope">
            {{ formatDate(scope.row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button size="small" type="primary" @click="handleOnlineModelDownload(scope.row.name)">下载</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 搜索在线模型抽屉 -->
    <el-drawer
        v-model="searchDrawerVisible"
        title="搜索在线模型"
        direction="rtl"
        size="50%"
    >
      <div class="search-drawer-content">
        <el-input
            v-model="searchQuery"
            placeholder="输入模型名称进行搜索 (例如: llama3)"
            clearable
            @input="handleSearch"
            style="margin-bottom: 20px;"
        >
          <template #prepend>
            <el-button :icon="Search"/>
          </template>
        </el-input>

        <el-table :data="searchResults" style="width: 100%" v-loading="isSearching">
          <el-table-column prop="name" label="模型名称"/>
          <el-table-column prop="pull_count" label="下载次数">
            <template #default="scope">
              {{ formatPullCount(scope.row.pull_count) }}
            </template>
          </el-table-column>
          <el-table-column prop="updated_at" label="更新时间">
            <template #default="scope">
              {{ formatDate(scope.row.updated_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120">
            <template #default="scope">
              <el-button size="small" type="primary" @click="handleDownloadFromSearch(scope.row.name)">下载</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-drawer>
    
    <!-- 模型详情对话框 -->
    <el-dialog v-model="dialogVisible" :title="selectedModel?.name" width="50%">
      <div v-if="selectedModel">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="模型名称">{{ selectedModel.name }}</el-descriptions-item>
          <el-descriptions-item label="描述">{{ selectedModel.description }}</el-descriptions-item>
          <el-descriptions-item label="大小">{{ formatSize(selectedModel.size) }}</el-descriptions-item>
          <el-descriptions-item label="标签">{{ selectedModel.tags?.join(', ') }}</el-descriptions-item>
        </el-descriptions>
        <div style="margin-top: 20px">
          <p><strong>使用说明:</strong></p>
          <p>下载后可在模型管理页面使用该模型</p>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="downloadModel(selectedModel)">下载</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { SearchOnlineModels } from '../../wailsjs/go/main/App'

interface MarketModel {
  name: string
  description: string
  size: number
  tags?: string[]
}

interface SearchResult {
  name: string
  pull_count: number
  updated_at: string
}

const searchText = ref('')
const marketModels = ref<MarketModel[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const selectedModel = ref<MarketModel | null>(null)

// 搜索在线模型相关变量
const searchDrawerVisible = ref(false)
const searchQuery = ref('')
const searchResults = ref<SearchResult[]>([])
const isSearching = ref(false)
const onlineModels = ref<SearchResult[]>([]) // 添加在线模型列表
const loadingOnlineModels = ref(false) // 添加加载状态

// 模拟模型市场数据
/*
const mockMarketModels: MarketModel[] = [
  {
    name: 'llama3',
    description: 'Meta开发的新一代大型语言模型',
    size: 4600000000,
    tags: ['Meta', 'LLM', 'Open Source']
  },
  {
    name: 'mistral',
    description: 'Mistral AI开发的高效大型语言模型',
    size: 3800000000,
    tags: ['Mistral AI', 'LLM', 'Efficient']
  },
  {
    name: 'phi3',
    description: 'Microsoft开发的小型但功能强大的模型',
    size: 1800000000,
    tags: ['Microsoft', 'Compact', 'Efficient']
  },
  {
    name: 'gemma',
    description: 'Google开发的轻量级但功能强大的模型',
    size: 2500000000,
    tags: ['Google', 'Lightweight', 'Open Source']
  }
]
*/

// 格式化文件大小
const formatSize = (size: number) => {
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(2) + ' KB'
  if (size < 1024 * 1024 * 1024) return (size / (1024 * 1024)).toFixed(2) + ' MB'
  return (size / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
}

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

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.getFullYear() + '年' +
      (date.getMonth() + 1).toString().padStart(2, '0') + '月' +
      date.getDate().toString().padStart(2, '0') + '日 ' +
      date.getHours().toString().padStart(2, '0') + ':' +
      date.getMinutes().toString().padStart(2, '0') + ':' +
      date.getSeconds().toString().padStart(2, '0')
}

// 搜索模型
const searchModels = async () => {
  if (!searchText.value) {
    // 如果搜索文本为空，加载所有在线模型
    await loadOnlineModels()
    return
  }
  
  // 如果有搜索文本，执行在线搜索
  loading.value = true
  try {
    const results = await SearchOnlineModels(searchText.value)
    onlineModels.value = results as SearchResult[]
  } catch (error: any) {
    ElMessage.error('搜索模型失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 显示模型详情
const showModelDetails = (model: MarketModel) => {
  selectedModel.value = model
  dialogVisible.value = true
}

// 下载模型
const downloadModel = (model: MarketModel | null) => {
  if (!model) {
    ElMessage.warning('请选择一个模型')
    return
  }
  
  dialogVisible.value = false
  ElMessage.success(`开始下载模型: ${model.name}`)
  // 这里应该调用实际的下载API
}

// 打开搜索抽屉
const openSearchDrawer = () => {
  searchDrawerVisible.value = true
  searchResults.value = []
  searchQuery.value = ''
}

// 处理搜索
const handleSearch = async () => {
  if (!searchQuery.value.trim()) {
    searchResults.value = []
    return
  }
  isSearching.value = true
  try {
    const results = await SearchOnlineModels(searchQuery.value)
    searchResults.value = results
  } catch (error: any) {
    ElMessage.error('搜索模型失败: ' + error.message)
  } finally {
    isSearching.value = false
  }
}

// 处理从搜索结果下载
const handleDownloadFromSearch = (modelName: string) => {
  if (!modelName) {
    ElMessage.warning('无效的模型名称')
    return
  }
  
  searchDrawerVisible.value = false
  ElMessage.success(`开始下载模型: ${modelName}`)
  // 这里应该调用实际的下载API
}

// 处理在线模型下载
const handleOnlineModelDownload = (modelName: string) => {
  if (!modelName) {
    ElMessage.warning('无效的模型名称')
    return
  }
  
  ElMessage.success(`开始下载模型: ${modelName}`)
  // 这里应该调用实际的下载API
}

// 加载在线模型列表
const loadOnlineModels = async () => {
  loadingOnlineModels.value = true
  try {
    const results = await SearchOnlineModels('')
    onlineModels.value = results
  } catch (error: any) {
    ElMessage.error('加载在线模型列表失败: ' + error.message)
  } finally {
    loadingOnlineModels.value = false
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
</style>