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
                <el-button type="primary" @click="openDownloadDialog">下载模型</el-button>
                <el-button @click="openDownloadQueueDialog">下载队列</el-button>
                <el-button @click="refreshModels">刷新</el-button>
              </div>
            </div>
          </template>
          <el-table :data="localModels" style="width: 100%" v-loading="loading" empty-text="暂无数据">
            <el-table-column type="index" label="#" width="60"/>
            <el-table-column prop="name" label="模型名称"/>
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

    <!-- 下载模型对话框 -->
    <el-dialog v-model="downloadDialogVisible" title="下载模型" width="30%">
      <el-alert
          title="提示"
          type="info"
          show-icon
          :closable="false"
          style="margin-bottom: 20px;"
      >
        请从 <a href="#" @click.prevent="openOllamaLibrary"
                class="el-link el-link--primary">https://ollama.com/library</a> 复制模型名称 (例如: llama3:8b) 并粘贴到下方。
      </el-alert>
      <el-input v-model="modelNameToDownload" placeholder="请输入要下载的模型名称"></el-input>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="downloadDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleDownload" :loading="isDownloading" :disabled="isDownloading">
            开始下载
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 下载队列抽屉 -->
    <el-drawer
        v-model="downloadQueueDrawerVisible"
        title="下载队列"
        direction="rtl"
        size="40%"
    >
      <el-table :data="downloadQueue" style="width: 100%" empty-text="暂无数据">
        <el-table-column prop="model" label="模型名称"/>
        <el-table-column prop="status" label="状态"/>
        <el-table-column label="进度">
          <template #default="scope">
            <el-progress :percentage="scope.row.percentage"/>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>

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
          <el-button v-if="!selectedModel.is_running" type="primary" @click="runModel" :loading="isRunningModel"
                     :disabled="isRunningModel">运行
          </el-button>
          <el-button type="danger" @click="deleteModel(selectedModel)">删除</el-button>
          <el-button v-if="selectedModel.is_running" @click="stopModel" :loading="isStoppingModel"
                     :disabled="isStoppingModel">停止
          </el-button>
        </div>

        <el-divider/>

        <div>
          <h4>模型参数</h4>
          <el-form :model="modelParams" label-width="80px" size="small">
            <el-form-item label="温度">
              <el-slider v-model="modelParams.temperature" :min="0" :max="1" :step="0.1"/>
            </el-form-item>
            <el-form-item label="Top P">
              <el-slider v-model="modelParams.topP" :min="0" :max="1" :step="0.1"/>
            </el-form-item>
            <el-form-item label="上下文">
              <el-input-number v-model="modelParams.context" :min="1" :max="32768"/>
            </el-form-item>
          </el-form>
          <div style="margin-top: 10px">
            <el-button @click="saveModelParams" type="primary" size="small">保存参数</el-button>
            <el-button @click="resetModelParams" size="small">重置参数</el-button>
          </div>
        </div>

        <el-divider/>

        <div>
          <h4>模型测试</h4>
          <el-form label-position="top">
            <el-form-item label="测试内容">
              <el-input
                  v-model="testPrompt"
                  type="textarea"
                  :rows="4"
                  placeholder="请输入要发送给模型的测试内容"
              />
            </el-form-item>
          </el-form>
          <el-button @click="testModel" :loading="isTestingModel" type="primary" size="small">发送测试</el-button>

          <div v-if="testResult || isTestingModel" style="margin-top: 15px;">
            <el-card shadow="never" v-loading="isTestingModel">
              <div style="white-space: pre-wrap; font-family: monospace; font-size: 14px;">
                {{ testResult }}
              </div>
            </el-card>
          </div>
        </div>

      </div>
      <div v-else>
        <p>请选择一个模型查看详情</p>
      </div>
    </el-drawer>

  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, onUnmounted, reactive, ref} from 'vue'
import {ElMessage, ElMessageBox, ElNotification} from 'element-plus'
import * as runtime from '../../wailsjs/runtime/runtime.js'
import {
  DeleteModel,
  DownloadModel,
  GetActiveServer,
  GetLocalServerTestStatus,
  GetModelParams,
  GetOllamaServerConfig,
  GetRemoteServers,
  ListModelsByServer,
  OpenInBrowser,
  RunModel,
  SetActiveServer,
  SetModelParams,
  StopModel,
  TestModel
} from '../../wailsjs/go/main/App'

const openOllamaLibrary = () => {
  OpenInBrowser('https://ollama.com/library')
}

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
  test_status?: 'unknown' | 'success' | 'failed' | string
}

interface ModelParams {
  temperature: number
  topP: number
  topK: number
  context: number
  numPredict: number
  repeatPenalty: number
}

interface DownloadProgress {
  model: string
  status: string
  percentage: number
  notification?: any
}

const localModels = ref<Model[]>([])
const selectedModel = ref<Model | null>(null)
const loading = ref(false)
const availableServers = ref<Server[]>([])
const selectedServer = ref<string>('local')
const drawerVisible = ref(false)

const downloadDialogVisible = ref(false)
const modelNameToDownload = ref('')
const isDownloading = ref(false)
const downloadProgresses = reactive<Record<string, DownloadProgress>>({})
const downloadQueueDrawerVisible = ref(false)  // 新增的抽屉可见性变量
const downloadQueue = computed(() => Object.values(downloadProgresses))

const isRunningModel = ref(false)
const isStoppingModel = ref(false)
const isTestingModel = ref(false)
const testPrompt = ref('你好，请用中文简单介绍一下自己。')
const testResult = ref('')

const modelParams = reactive<ModelParams>({
  temperature: 0.8,
  topP: 0.9,
  topK: 40,
  context: 2048,
  numPredict: 512,
  repeatPenalty: 1.1
})

const formatSize = (size: number) => {
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(2) + ' KB'
  if (size < 1024 * 1024 * 1024) return (size / (1024 * 1024)).toFixed(2) + ' MB'
  return (size / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.getFullYear() + '年' +
      (date.getMonth() + 1).toString().padStart(2, '0') + '月' +
      date.getDate().toString().padStart(2, '0') + '日 ' +
      date.getHours().toString().padStart(2, '0') + ':' +
      date.getMinutes().toString().padStart(2, '0') + ':' +
      date.getSeconds().toString().padStart(2, '0')
}

const resetModelParams = () => {
  modelParams.temperature = 0.8
  modelParams.topP = 0.9
  modelParams.topK = 40
  modelParams.context = 2048
  modelParams.numPredict = 512
  modelParams.repeatPenalty = 1.1
  ElMessage.info('参数已重置为默认值')
}

const loadAvailableServers = async () => {
  try {
    const localConfig = await GetOllamaServerConfig();
    const remoteList = await GetRemoteServers();
    const localTestStatus = await GetLocalServerTestStatus();
    let activeServerId = 'local'; // 默认本地

    // 正确映射 OllamaServerConfig 到 Server 结构
    const allServersRaw = [
      {id: 'local', name: '本地服务', baseUrl: localConfig, apiKey: '', isActive: false, test_status: localTestStatus},
      ...remoteList.map(server => ({
        id: server.id,
        name: server.name,
        baseUrl: server.base_url,
        apiKey: server.api_key,
        isActive: server.is_active,
        test_status: server.test_status
      }))
    ];

    const testedServers = allServersRaw.filter(s => s.test_status === 'success');

    availableServers.value = testedServers as Server[];

    try {
      const activeServer = await GetActiveServer();
      if (activeServer && activeServer.id && availableServers.value.some(s => s.id === activeServer.id)) {
        activeServerId = activeServer.id;
      } else if (availableServers.value.length > 0) {
        activeServerId = availableServers.value[0].id
        await SetActiveServer(activeServerId)
      } else {
        activeServerId = ''
      }
    } catch (e) {
      console.error("无法获取活动服务器，将默认使用第一个可用服务器。", e)
      if (availableServers.value.length > 0) {
        activeServerId = availableServers.value[0].id
      } else {
        activeServerId = ''
      }
    }

    // 更新 isActive 标志
    availableServers.value.forEach(s => {
      s.isActive = s.id === activeServerId;
    });

    selectedServer.value = activeServerId;

  } catch (error) {
    console.error('加载可用服务器列表失败:', error);
    availableServers.value = [];
    selectedServer.value = '';
  }
};

const onServerChange = () => {
  getModels()
}

const getModels = async () => {
  try {
    loading.value = true
    localModels.value = await ListModelsByServer(selectedServer.value)
  } catch (error: any) {
    ElMessage.error('获取模型列表失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const refreshModels = () => {
  getModels()
}

const viewModelDetails = (model: Model) => {
  selectedModel.value = {...model}
  drawerVisible.value = true
  loadModelParams(model.name)
}

const runModel = async () => {
  if (!selectedModel.value) return
  isRunningModel.value = true
  ElMessage.info(`正在启动模型 "${selectedModel.value.name}"...`)
  try {
    await RunModel(selectedModel.value.name, modelParams)
    selectedModel.value.is_running = true
    const index = localModels.value.findIndex(m => m.name === selectedModel.value!.name)
    if (index !== -1) localModels.value[index].is_running = true
    ElMessage.success(`模型 "${selectedModel.value.name}" 已启动`)
  } catch (error: any) {
    ElMessage.error('启动模型失败: ' + error.message)
    if (selectedModel.value) selectedModel.value.is_running = false
  } finally {
    isRunningModel.value = false
  }
}

const stopModel = async () => {
  if (!selectedModel.value) return
  isStoppingModel.value = true
  ElMessage.info(`正在停止模型 "${selectedModel.value.name}"...`)
  try {
    await StopModel(selectedModel.value.name)
    selectedModel.value.is_running = false
    const index = localModels.value.findIndex(m => m.name === selectedModel.value!.name)
    if (index !== -1) localModels.value[index].is_running = false
    ElMessage.success(`模型 "${selectedModel.value.name}" 已停止`)
  } catch (error: any) {
    ElMessage.error('停止模型失败: ' + error.message)
  } finally {
    isStoppingModel.value = false
  }
}

const deleteModel = (model: Model) => {
  ElMessageBox.confirm(
      `确定要删除模型 "${model.name}" 吗？此操作不可恢复。`,
      '删除模型',
      {confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning'}
  ).then(async () => {
    try {
      await DeleteModel(model.name)
      ElMessage.success('模型删除成功')
      refreshModels()
      drawerVisible.value = false
    } catch (error: any) {
      ElMessage.error('删除模型失败: ' + error.message)
    }
  }).catch(() => {
    ElMessage.info('已取消删除')
  })
}

const testModel = async () => {
  if (!selectedModel.value) return
  if (!testPrompt.value.trim()) {
    ElMessage.warning('请输入测试内容')
    return
  }
  isTestingModel.value = true
  testResult.value = ''
  try {
    const response = await TestModel(selectedModel.value.name, testPrompt.value)
    testResult.value = response
    ElMessage.success(`测试完成`)
  } catch (error: any) {
    testResult.value = '测试失败: ' + error.message
    ElMessage.error('测试模型失败')
  } finally {
    isTestingModel.value = false
  }
}

const openDownloadDialog = () => {
  modelNameToDownload.value = ''
  downloadDialogVisible.value = true
}

const handleDownload = async () => {
  const modelName = modelNameToDownload.value.trim()
  if (!modelName) {
    ElMessage.warning('请输入要下载的模型名称')
    return
  }
  if (downloadProgresses[modelName]) {
    ElMessage.warning(`模型 "${modelName}" 已在下载队列中。`)
    return
  }

  downloadProgresses[modelName] = reactive({
    model: modelName,
    status: '正在准备...',
    percentage: 0,
  }) as DownloadProgress

  isDownloading.value = true
  downloadDialogVisible.value = false

  // 自动拉出下载队列抽屉
  downloadQueueDrawerVisible.value = true

  DownloadModel(selectedServer.value, modelName)
  ElMessage.info(`已将模型 "${modelName}" 添加到下载队列。`)
  isDownloading.value = false
}

// 修改打开下载队列的函数
const openDownloadQueueDialog = () => {
  downloadQueueDrawerVisible.value = true
}

const loadModelParams = async (modelName: string) => {
  try {
    const params: any = await GetModelParams(modelName)
    Object.assign(modelParams, params)
  } catch (error) {
    resetModelParams()
  }
}

const saveModelParams = async () => {
  if (!selectedModel.value) return
  try {
    await SetModelParams(selectedModel.value.name, modelParams)
    ElMessage.success('参数保存成功')
  } catch (error: any) {
    ElMessage.error('参数保存失败: ' + error.message)
  }
}

const setupDownloadListeners = () => {
  // 保存事件处理函数的引用，以便后续移除
  const progressHandler = (data: any) => {
    console.log('收到下载进度:', data);
    const {model, status, completed, total} = data;

    // 确保下载进度对象存在
    if (!downloadProgresses[model]) {
      downloadProgresses[model] = reactive({
        model: model,
        status: '',
        percentage: 0,
      }) as DownloadProgress;
    }

    // 计算进度百分比
    let progress = 0;
    if (total && total > 0) {
      progress = Math.round((completed / total) * 100);
    } else if (status === "success" || status.includes("success")) {
      progress = 100;
    }

    // 更新状态和进度
    downloadProgresses[model].status = status || '下载中';
    downloadProgresses[model].percentage = progress;

    console.log(`模型 ${model} 下载进度: ${progress}%, 状态: ${status}`);
  };

  const doneHandler = (data: any) => {
    console.log('下载完成:', data);
    const {model} = data;

    // 从下载进度中移除
    if (downloadProgresses[model]) {
      delete downloadProgresses[model];
    }

    ElNotification.success({
      title: '下载完成',
      message: `模型 "${model}" 已成功下载。`,
      duration: 3000
    });

    // 刷新模型列表
    refreshModels();
  };

  const errorHandler = (data: any) => {
    console.error('下载出错:', data);
    const {model, error} = data;

    // 从下载进度中移除
    if (downloadProgresses[model]) {
      delete downloadProgresses[model];
    }

    ElNotification.error({
      title: '下载失败',
      message: `模型 "${model}" 下载失败: ${error}`,
      duration: 0
    });

    console.error(`模型 "${model}" 下载失败:`, error);
  };

  // 添加事件监听器
  runtime.EventsOn(`model:download:progress`, progressHandler);
  runtime.EventsOn('model:download:done', doneHandler);
  runtime.EventsOn('model:download:error', errorHandler);

  // 返回清理函数
  return () => {
    runtime.EventsOff('model:download:progress');
    runtime.EventsOff('model:download:done');
    runtime.EventsOff('model:download:error');

  };
}

let cleanupDownloadListeners: Function | null = null;

onMounted(async () => {
  await loadAvailableServers()
  await getModels()
  cleanupDownloadListeners = setupDownloadListeners()
})

onUnmounted(() => {
  // 清理所有事件监听器
  Object.keys(downloadProgresses).forEach(modelName => {
    if (downloadProgresses[modelName] && downloadProgresses[modelName].notification) {
      downloadProgresses[modelName].notification.close()
    }
  });

  // 调用清理函数移除事件监听器
  if (cleanupDownloadListeners && typeof cleanupDownloadListeners === 'function') {
    cleanupDownloadListeners();
  }
});

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
