<template>
  <div class="model-manager">
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <div>
                <el-select 
                  v-model="selectedServer" 
                  placeholder="选择服务" 
                  style="width: 200px; margin-right: 10px;"
                  @change="onServerChange"
                >
                  <el-option
                    v-for="server in availableServers"
                    :key="server.id"
                    :label="server.name"
                    :value="server.id"
                  />
                </el-select>
                <el-button type="primary" @click="openModelSearch">搜索模型</el-button>
                <el-button @click="refreshModels">刷新</el-button>
              </div>
            </div>
          </template>
          <el-table :data="localModels" style="width: 100%" v-loading="loading">
            <el-table-column type="index" label="#" width="60" />
            <el-table-column prop="name" label="模型名称" />
            <el-table-column prop="size" label="大小">
              <template #default="scope">
                {{ formatSize(scope.row.size) }}
              </template>
            </el-table-column>
            <el-table-column prop="modified_at" label="修改时间">
              <template #default="scope">
                {{ formatDate(scope.row.modified_at) }}
              </template>
            </el-table-column>
            <el-table-column label="运行状态" width="100">
              <template #default="scope">
                <el-tag v-if="scope.row.is_running" type="success">运行中</el-tag>
                <el-tag v-else type="info">未运行</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="scope">
                <el-button size="small" @click="viewModelDetails(scope.row)">查看</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 模型详情抽屉 -->
    <el-drawer
      v-model="drawerVisible"
      title="模型详情"
      direction="rtl"
      size="40%"
    >
      <div v-if="selectedModel">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="模型名称">{{ selectedModel.name }}</el-descriptions-item>
          <el-descriptions-item label="大小">{{ formatSize(selectedModel.size) }}</el-descriptions-item>
          <el-descriptions-item label="修改时间">{{ formatDate(selectedModel.modified_at) }}</el-descriptions-item>
          <el-descriptions-item label="运行状态">
            <el-tag v-if="selectedModel.is_running" type="success">运行中</el-tag>
            <el-tag v-else type="info">未运行</el-tag>
          </el-descriptions-item>
        </el-descriptions>
        
        <div style="margin-top: 20px">
          <el-button type="primary" @click="runModel">运行</el-button>
          <el-button @click="testModel">测试</el-button>
          <el-button type="danger" @click="deleteModel(selectedModel)">删除</el-button>
          <el-button v-if="selectedModel.is_running" @click="stopModel">停止</el-button>
        </div>
        
        <div style="margin-top: 20px">
          <h4>模型参数</h4>
          <el-form :model="modelParams" label-width="80px" size="small">
            <el-form-item label="温度">
              <el-slider v-model="modelParams.temperature" :min="0" :max="1" :step="0.1" />
            </el-form-item>
            <el-form-item label="Top P">
              <el-slider v-model="modelParams.topP" :min="0" :max="1" :step="0.1" />
            </el-form-item>
            <el-form-item label="上下文">
              <el-input-number v-model="modelParams.context" :min="1" :max="32768" />
            </el-form-item>
          </el-form>
          <div style="margin-top: 10px">
            <el-button @click="saveModelParams" type="primary" size="small">保存参数</el-button>
          </div>
        </div>
      </div>
      <div v-else>
        <p>请选择一个模型查看详情</p>
      </div>
    </el-drawer>
    
    <!-- 搜索模型抽屉 -->
    <el-drawer
      v-model="searchDrawerVisible"
      title="搜索模型"
      direction="rtl"
      size="40%"
    >
      <div style="margin-bottom: 20px">
        <el-input
          v-model="searchModelQuery"
          placeholder="请输入模型名称"
          clearable
          @keyup.enter="searchModels"
        >
          <template #append>
            <el-button @click="searchModels">搜索</el-button>
          </template>
        </el-input>
      </div>
      
      <div v-if="searchResults.length > 0">
        <el-collapse v-model="activeNames">
          <el-collapse-item 
            v-for="(group, index) in modelGroups" 
            :key="index" 
            :name="index"
            :title="group.name"
          >
            <el-table :data="group.models" style="width: 100%">
              <el-table-column prop="name" label="模型名称" />
              <el-table-column prop="size" label="大小">
                <template #default="scope">
                  {{ formatSize(scope.row.size) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="100">
                <template #default="scope">
                  <el-button size="small" type="success" @click="downloadModelFromSearch(scope.row)">下载</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-collapse-item>
        </el-collapse>
      </div>
      
      <div v-else-if="searchModelQuery && searchResults.length === 0">
        <el-empty description="未找到相关模型" />
      </div>
      
      <div v-else>
        <el-alert
          title="提示"
          description="请输入模型名称进行搜索，例如 'llama' 或 'mistral'"
          type="info"
          show-icon
        />
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  ListModelsByServer, 
  DeleteModel, 
  GetRemoteServers,
  GetOllamaServerConfig,
  RunModel,
  StopModel,
  TestModel,
  SetModelParams,
  GetModelParams,
  SearchModels,
  GetModelFamilies,
  GetModelTags
} from '../../wailsjs/go/main/App'

interface Model {
  name: string
  size: number
  modified_at: string
  is_running?: boolean
}

interface Server {
  id: string
  name: string
  baseUrl: string
  apiKey: string
  isActive: boolean
}

// 模型参数接口
interface ModelParams {
  temperature: number
  topP: number
  topK: number
  context: number
  numPredict: number
  repeatPenalty: number
}

const localModels = ref<Model[]>([])
const selectedModel = ref<Model | null>(null)
const loading = ref(false)
const availableServers = ref<Server[]>([])
const selectedServer = ref<string>('local')
const drawerVisible = ref(false)
const searchDrawerVisible = ref(false)
const searchModelQuery = ref('')
const searchResults = ref<Model[]>([])
const activeNames = ref([])
const downloading = ref(false)
const downloadProgress = ref(0)

// 模型参数
const modelParams = reactive<ModelParams>({
  temperature: 0.8,
  topP: 0.9,
  context: 2048
})

// 格式化文件大小
const formatSize = (size: number) => {
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(2) + ' KB'
  if (size < 1024 * 1024 * 1024) return (size / (1024 * 1024)).toFixed(2) + ' MB'
  return (size / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
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

// 模型分组
const modelGroups = computed(() => {
  const groups: { name: string; models: Model[] }[] = []
  const groupMap: { [key: string]: Model[] } = {}
  
  searchResults.value.forEach(model => {
    // 以模型名称的第一部分作为分组依据（例如：llama3:8b 中的 llama3）
    const groupName = model.name.split(':')[0]
    if (!groupMap[groupName]) {
      groupMap[groupName] = []
    }
    groupMap[groupName].push(model)
  })
  
  Object.keys(groupMap).forEach(key => {
    groups.push({
      name: key,
      models: groupMap[key]
    })
  })
  
  return groups
})

// 获取可用的服务列表
const loadAvailableServers = async () => {
  try {
    // 获取本地配置
    const localBaseUrl = await GetOllamaServerConfig()
    const localServer: Server = {
      id: 'local',
      name: '本地服务',
      baseUrl: localBaseUrl,
      apiKey: '',
      isActive: true
    }
    
    let remoteServers: Server[] = []
    
    try {
      // 获取远程服务器列表
      remoteServers = await GetRemoteServers()
    } catch (remoteError) {
      console.error('获取远程服务器列表失败:', remoteError)
      ElMessage.warning('无法获取远程服务器列表，将仅使用本地服务')
    }
    
    // 合并本地和远程服务器
    availableServers.value = [localServer, ...remoteServers]
    
    // 设置默认选中的服务器
    // 首先检查是否有活动的远程服务器
    const activeRemoteServer = remoteServers.find((server: Server) => server.isActive)
    if (activeRemoteServer) {
      selectedServer.value = activeRemoteServer.id
    } else {
      // 如果没有活动的远程服务器，默认选择本地服务器
      selectedServer.value = 'local'
    }
  } catch (error) {
    console.error('加载服务列表失败:', error)
    ElMessage.error('加载服务列表失败: ' + (error as Error).message)
    
    // 出现错误时默认设置为本地服务
    selectedServer.value = 'local'
  }
}

// 当服务器选择改变时
const onServerChange = () => {
  getModels()
}

// 获取模型列表
const getModels = async () => {
  try {
    loading.value = true
    
    // 根据选择的服务获取模型列表
    const models: Model[] = await ListModelsByServer(selectedServer.value)
    // 添加默认状态
    localModels.value = models.map(model => ({
      ...model,
      is_running: false // 默认未运行
    }))
  } catch (error) {
    ElMessage.error('获取模型列表失败: ' + (error as Error).message)
  } finally {
    loading.value = false
  }
}

// 刷新模型列表
const refreshModels = () => {
  getModels()
}

// 查看模型详情
const viewModelDetails = (model: Model) => {
  selectedModel.value = {
    ...model,
    is_running: model.is_running !== undefined ? model.is_running : false
  }
  drawerVisible.value = true
  
  // 加载模型参数
  loadModelParams(model.name)
}

// 运行模型
const runModel = async () => {
  if (selectedModel.value) {
    try {
      await RunModel(selectedModel.value.name, {
        temperature: modelParams.temperature,
        top_p: modelParams.topP,
        top_k: modelParams.topK,
        context: modelParams.context,
        num_predict: modelParams.numPredict,
        repeat_penalty: modelParams.repeatPenalty
      })
      
      selectedModel.value.is_running = true
      // 更新本地模型列表中的状态
      const index = localModels.value.findIndex(m => m.name === selectedModel.value!.name)
      if (index !== -1) {
        localModels.value[index].is_running = true
      }
      ElMessage.success(`模型 "${selectedModel.value.name}" 已启动`)
    } catch (error) {
      console.error('启动模型失败:', error)
      ElMessage.error('启动模型失败: ' + (error as Error).message)
    }
  } else {
    ElMessage.warning('请先选择一个模型')
  }
}

// 停止模型
const stopModel = async () => {
  if (selectedModel.value) {
    try {
      await StopModel(selectedModel.value.name)
      
      selectedModel.value.is_running = false
      // 更新本地模型列表中的状态
      const index = localModels.value.findIndex(m => m.name === selectedModel.value!.name)
      if (index !== -1) {
        localModels.value[index].is_running = false
      }
      ElMessage.success(`模型 "${selectedModel.value.name}" 已停止`)
    } catch (error) {
      console.error('停止模型失败:', error)
      ElMessage.error('停止模型失败: ' + (error as Error).message)
    }
  } else {
    ElMessage.warning('请先选择一个模型')
  }
}

// 删除模型
const deleteModel = (model: Model) => {
  ElMessageBox.confirm(
    `确定要删除模型 "${model.name}" 吗？此操作不可恢复。`,
    '删除模型',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await DeleteModel(model.name)
      ElMessage.success('模型删除成功')
      refreshModels()
      // 关闭抽屉
      drawerVisible.value = false
    } catch (error) {
      ElMessage.error('删除模型失败: ' + (error as Error).message)
    }
  }).catch(() => {
    ElMessage.info('已取消删除')
  })
}

// 测试模型
const testModel = async () => {
  if (selectedModel.value) {
    try {
      ElMessage.info(`正在测试模型: ${selectedModel.value.name}`)
      const response = await TestModel(selectedModel.value.name)
      ElMessage.success(`测试完成: ${response.substring(0, 100)}...`)
    } catch (error) {
      console.error('测试模型失败:', error)
      ElMessage.error('测试模型失败: ' + (error as Error).message)
    }
  } else {
    ElMessage.warning('请先选择一个模型')
  }
}

// 下载模型
const downloadModel = () => {
  if (selectedModel.value) {
    downloading.value = true
    downloadProgress.value = 0
    
    // 模拟下载进度
    const timer = setInterval(() => {
      downloadProgress.value += 10
      if (downloadProgress.value >= 100) {
        clearInterval(timer)
        downloading.value = false
        if (selectedModel.value) {
          selectedModel.value.is_running = false
        }
        ElMessage.success('模型下载完成')
      }
    }, 300)
  } else {
    ElMessage.warning('请先选择一个模型')
  }
}

// 打开模型搜索
const openModelSearch = () => {
  searchDrawerVisible.value = true
  searchModelQuery.value = ''
  searchResults.value = []
  selectedFamilies.value = []
  selectedTags.value = []
}

// 加载搜索模型数据
const loadSearchModels = async () => {
  try {
    // 加载模型家族和标签
    modelFamilies.value = await GetModelFamilies()
    modelTags.value = await GetModelTags()
    
    // 加载所有模型
    await searchModels()
  } catch (error) {
    ElMessage.error('加载搜索数据失败: ' + (error as Error).message)
  }
}

// 搜索模型
const searchModels = async () => {
  try {
    const params: any = {
      families: selectedFamilies.value,
      tags: selectedTags.value
    }
    
    if (searchModelQuery.value.trim()) {
      params.query = searchModelQuery.value.trim()
    }
    
    const results: Model[] = await SearchModels(params)
    searchResults.value = results
  } catch (error) {
    ElMessage.error('搜索模型失败: ' + (error as Error).message)
  }
}

// 从搜索结果下载模型
const downloadModelFromSearch = async (model: Model) => {
  try {
    ElMessage.success(`开始在 ${getCurrentServerName()} 中下载模型: ${model.name}`)
    // 关闭搜索抽屉
    searchDrawerVisible.value = false
    // 刷新模型列表
    setTimeout(() => {
      refreshModels()
    }, 1000)
  } catch (error) {
    ElMessage.error('下载模型失败: ' + (error as Error).message)
  }
}

// 获取当前选中服务的名称
const getCurrentServerName = () => {
  const server = availableServers.value.find(s => s.id === selectedServer.value)
  return server ? server.name : '未知服务'
}

// 加载模型参数
const loadModelParams = async (modelName: string) => {
  try {
    const params = await GetModelParams(modelName)
    modelParams.temperature = params.temperature
    modelParams.topP = params.top_p
    modelParams.topK = params.top_k
    modelParams.context = params.context
    modelParams.numPredict = params.num_predict
    modelParams.repeatPenalty = params.repeat_penalty
  } catch (error) {
    console.error('加载模型参数失败:', error)
    // 使用默认参数
    modelParams.temperature = 0.8
    modelParams.topP = 0.9
    modelParams.topK = 40
    modelParams.context = 2048
    modelParams.numPredict = 512
    modelParams.repeatPenalty = 1.1
  }
}

// 保存模型参数
const saveModelParams = async () => {
  if (!selectedModel.value) {
    ElMessage.warning('请先选择一个模型')
    return
  }
  
  try {
    await SetModelParams(selectedModel.value.name, {
      temperature: modelParams.temperature,
      top_p: modelParams.topP,
      top_k: modelParams.topK,
      context: modelParams.context,
      num_predict: modelParams.numPredict,
      repeat_penalty: modelParams.repeatPenalty
    })
    ElMessage.success('参数保存成功')
  } catch (error) {
    ElMessage.error('参数保存失败: ' + (error as Error).message)
  }
}

// 重置模型参数
const resetModelParams = () => {
  modelParams.temperature = 0.8
  modelParams.topP = 0.9
  modelParams.topK = 40
  modelParams.context = 2048
  modelParams.numPredict = 512
  modelParams.repeatPenalty = 1.1
  ElMessage.info('参数已重置为默认值')
}

onMounted(() => {
  loadAvailableServers()
  getModels()
  loadSearchModels()
})
</script>

<style scoped>
.model-manager {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>