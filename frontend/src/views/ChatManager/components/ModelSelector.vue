<template>
  <el-card class="model-selector">
    <template #header>
      <div class="card-header">
        <span>模型选择</span>
      </div>
    </template>
    <el-select
        :model-value="selectedServer"
        placeholder="选择服务"
        style="width: 100%; margin-bottom: 10px;"
        @update:modelValue="updateSelectedServer"
        @change="onServerChange">
      <el-option
          v-for="server in availableServers"
          :key="server.id"
          :label="server.name"
          :value="server.id"
      />
    </el-select>
    <el-select
        :model-value="selectedModel"
        placeholder="选择模型"
        style="width: 100%"
        @update:modelValue="updateSelectedModel">
      <el-option
          v-for="model in localModels"
          :key="model.name"
          :label="model.name"
          :value="model.name"
      />
    </el-select>

    <div style="margin-top: 20px">
      <el-button @click="loadModel" style="width: 100%">加载模型</el-button>
    </div>

    <div style="margin-top: 20px">
      <h4>参数设置</h4>
      <el-form :model="modelParams" label-width="80px" size="small">
        <el-form-item label="温度">
          <el-slider
              :model-value="modelParams.temperature"
              :min="0"
              :max="1"
              :step="0.1"
              @update:modelValue="updateTemperature"/>
        </el-form-item>
        <el-form-item label="Top P">
          <el-slider
              :model-value="modelParams.topP"
              :min="0"
              :max="1"
              :step="0.1"
              @update:modelValue="updateTopP"/>
        </el-form-item>
        <el-form-item label="上下文">
          <el-input-number
              :model-value="modelParams.context"
              :min="1"
              :max="32768"
              @update:modelValue="updateContext"/>
        </el-form-item>
        <el-form-item label="预测数">
          <el-input-number
              :model-value="modelParams.numPredict"
              :min="1"
              :max="4096"
              @update:modelValue="updateNumPredict"/>
        </el-form-item>
        <el-form-item label="Top K">
          <el-input-number
              :model-value="modelParams.topK"
              :min="1"
              :max="100"
              @update:modelValue="updateTopK"/>
        </el-form-item>
        <el-form-item label="重复惩罚">
          <el-input-number
              :model-value="modelParams.repeatPenalty"
              :min="0.1"
              :max="2"
              :step="0.1"
              @update:modelValue="updateRepeatPenalty"/>
        </el-form-item>
        <el-form-item label="输出方式">
          <el-select
              :model-value="modelParams.outputMode"
              placeholder="选择输出方式"
              @update:modelValue="updateOutputMode">
            <el-option label="流式输出" value="stream"/>
            <el-option label="阻塞输出" value="blocking"/>
          </el-select>
        </el-form-item>
      </el-form>
      <div style="margin-top: 10px">
        <el-button @click="saveModelParams" type="primary" size="small">保存参数</el-button>
        <el-button @click="resetModelParams" size="small">重置</el-button>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import {onMounted} from 'vue'
import {ElMessage} from 'element-plus'
import {GetActiveServer, GetServers, ListModelsByServer} from '../../../../wailsjs/go/main/App'

interface Model {
  name: string
  size: number
  modified_at: string
}

interface Server {
  id: string
  name: string
  baseUrl: string
  apiKey: string
  isActive: boolean
  testStatus: string
  type: string
}

interface ModelParams {
  temperature: number
  topP: number
  context: number
  numPredict: number
  topK: number
  repeatPenalty: number
  outputMode: 'stream' | 'blocking' // 添加输出方式选项
}

const props = defineProps<{
  localModels: Model[]
  selectedModel: string
  availableServers: Server[]
  selectedServer: string
  modelParams: ModelParams
}>()

const emit = defineEmits<{
  (e: 'update:selectedModel', value: string): void
  (e: 'update:selectedServer', value: string): void
  (e: 'update:modelParams', value: ModelParams): void
  (e: 'load-model'): void
  (e: 'save-model-params'): void
  (e: 'reset-model-params'): void
  (e: 'server-change'): void
}>()

// 更新模型选择
const updateSelectedModel = (value: string) => {
  emit('update:selectedModel', value)
}

// 更新服务器选择
const updateSelectedServer = (value: string) => {
  emit('update:selectedServer', value)
}

// 更新温度参数
const updateTemperature = (value: number) => {
  emit('update:modelParams', {...props.modelParams, temperature: value})
}

// 更新Top P参数
const updateTopP = (value: number) => {
  emit('update:modelParams', {...props.modelParams, topP: value})
}

// 更新上下文参数
const updateContext = (value: number) => {
  emit('update:modelParams', {...props.modelParams, context: value})
}

// 更新预测数参数
const updateNumPredict = (value: number) => {
  emit('update:modelParams', {...props.modelParams, numPredict: value})
}

// 更新Top K参数
const updateTopK = (value: number) => {
  emit('update:modelParams', {...props.modelParams, topK: value})
}

// 更新重复惩罚参数
const updateRepeatPenalty = (value: number) => {
  emit('update:modelParams', {...props.modelParams, repeatPenalty: value})
}

// 更新输出方式参数
const updateOutputMode = (value: 'stream' | 'blocking') => {
  emit('update:modelParams', {...props.modelParams, outputMode: value})
}

const loadAvailableServers = async () => {
  try {
    // const localBaseUrl = await GetOllamaServerConfig()
    // const localServer = {
    //   id: 'local',
    //   name: '本地服务',
    //   baseUrl: localBaseUrl,
    //   apiKey: '',
    //   isActive: true,
    //   testStatus: '',
    //   type: 'local'
    // }
    let remoteServers: Server[] = []
    try {
      const remoteList: any[] = await GetServers()
      remoteServers = remoteList.map(server => ({
        id: server.id || server.ID,
        name: server.name || server.Name,
        baseUrl: server.baseUrl || server.base_url || server.BaseURL,
        apiKey: server.apiKey || server.api_key || server.APIKey,
        isActive: server.isActive !== undefined ? server.isActive : (server.is_active !== undefined ? server.is_active : server.IsActive),
        testStatus: server.testStatus || server.test_status || server.TestStatus || '',
        type: server.type || server.Type || ''
      }))
    } catch (remoteError) {
      console.error('获取远程服务器列表失败:', remoteError)
    }
    emit('update:availableServers', [...remoteServers])

    try {
      const activeServer = await GetActiveServer()
      const activeServerExists = activeServer && activeServer.id && [...remoteServers].some(s => s.id === activeServer.id)
      if (activeServerExists) {
        emit('update:selectedServer', activeServer.id)
      } else {
        emit('update:selectedServer', 'local')
      }
    } catch (e) {
      emit('update:selectedServer', 'local')
    }
  } catch (error) {
    const localServer = {
      id: 'local',
      name: '本地服务',
      baseUrl: '',
      apiKey: '',
      isActive: true,
      testStatus: '',
      type: 'local'
    };
    emit('update:availableServers', [localServer])
    emit('update:selectedServer', 'local')
  }
}

const onServerChange = () => {
  emit('server-change')
}

// 获取模型列表
const getModels = async () => {
  try {
    const models: Model[] = await ListModelsByServer(props.selectedServer)
    // 更新本地模型列表
    // 这里需要通过事件通知父组件更新
  } catch (error: any) {
    console.error('获取模型列表失败:', error)
    let errorMessage = '获取模型列表失败'
    if (error.message) {
      errorMessage += ': ' + error.message
    } else if (error.toString() !== '[object Object]') {
      errorMessage += ': ' + error.toString()
    } else {
      errorMessage += ': 未知错误'
    }
    ElMessage.error(errorMessage)
  }
}

// 加载模型
const loadModel = () => {
  emit('load-model')
}

// 保存模型参数
const saveModelParams = async () => {
  emit('save-model-params')
}

// 重置模型参数
const resetModelParams = () => {
  emit('reset-model-params')
}

onMounted(async () => {
  await loadAvailableServers()
  await getModels()
})
</script>

<style scoped>
.model-selector {
  height: 100%;
}

.model-selector :deep(.el-card__body) {
  height: calc(100% - 60px);
  overflow-y: auto;
}
</style>