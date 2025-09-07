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
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

interface MarketModel {
  name: string
  description: string
  size: number
  tags?: string[]
}

const searchText = ref('')
const marketModels = ref<MarketModel[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const selectedModel = ref<MarketModel | null>(null)

// 模拟模型市场数据
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

// 格式化文件大小
const formatSize = (size: number) => {
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(2) + ' KB'
  if (size < 1024 * 1024 * 1024) return (size / (1024 * 1024)).toFixed(2) + ' MB'
  return (size / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
}

// 搜索模型
const searchModels = () => {
  loading.value = true
  setTimeout(() => {
    if (searchText.value) {
      marketModels.value = mockMarketModels.filter(model => 
        model.name.toLowerCase().includes(searchText.value.toLowerCase()) ||
        model.description.toLowerCase().includes(searchText.value.toLowerCase())
      )
    } else {
      marketModels.value = [...mockMarketModels]
    }
    loading.value = false
  }, 500)
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

onMounted(() => {
  marketModels.value = [...mockMarketModels]
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