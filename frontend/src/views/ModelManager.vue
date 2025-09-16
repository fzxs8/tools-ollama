<template>
  <div class="model-manager">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-icon">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
            <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
            <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
          </svg>
        </div>
        <div class="header-text">
          <h1>{{ t('modelManager.title') }}</h1>
          <p>{{ t('modelManager.description') }}</p>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="main-content">
      <!-- Control Panel -->
      <div class="control-panel">
        <div class="server-selector">
          <label>{{ t('modelManager.server') }}:</label>
          <div class="select-wrapper">
            <select v-model="selectedServer" @change="onServerChange">
              <option value="" disabled>{{ t('modelManager.selectServer') }}</option>
              <option v-for="server in availableServers" :key="server.id" :value="server.id">
                {{ server.name }}
              </option>
            </select>
            <svg class="select-arrow" width="12" height="12" viewBox="0 0 24 24" fill="none">
              <path d="M6 9L12 15L18 9" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
        </div>
        
        <div class="action-buttons">
          <button class="btn primary" @click="openDownloadDialog">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
              <path d="M21 15V19A2 2 0 0 1 19 21H5A2 2 0 0 1 3 19V15" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <polyline points="7,10 12,15 17,10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <line x1="12" y1="15" x2="12" y2="3" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            {{ t('modelManager.downloadModel') }}
          </button>
          
          <button class="btn secondary" @click="openDownloadQueueDialog">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
              <path d="M16 4H18A2 2 0 0 1 20 6V20A2 2 0 0 1 18 22H6A2 2 0 0 1 4 20V6A2 2 0 0 1 6 4H8" stroke="currentColor" stroke-width="2"/>
              <rect x="8" y="2" width="8" height="4" rx="1" ry="1" stroke="currentColor" stroke-width="2"/>
            </svg>
            {{ t('modelManager.downloadQueue') }}
          </button>
          
          <button class="btn secondary" @click="refreshModels">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
              <polyline points="23,4 23,10 17,10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <polyline points="1,20 1,14 7,14" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M20.49 9A9 9 0 0 0 5.64 5.64L1 10M22.88 14.36A9 9 0 0 1 8.51 18.36L4 14" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            {{ t('common.refresh') }}
          </button>
        </div>
      </div>

      <!-- Models Table -->
      <div class="models-table-container">
        <div class="table-wrapper" v-loading="loading">
          <table class="models-table">
            <thead>
              <tr>
                <th>#</th>
                <th>{{ t('modelManager.modelName') }}</th>
                <th>{{ t('modelManager.size') }}</th>
                <th>{{ t('modelManager.modified') }}</th>
                <th>{{ t('modelManager.status') }}</th>
                <th>{{ t('modelManager.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="localModels.length === 0 && !loading">
                <td colspan="6" class="empty-state">
                  <div class="empty-content">
                    <svg width="48" height="48" viewBox="0 0 24 24" fill="none">
                      <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
                      <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
                      <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
                    </svg>
                    <p>{{ t('modelManager.noModels') }}</p>
                    <span>{{ t('modelManager.downloadToStart') }}</span>
                  </div>
                </td>
              </tr>
              <tr v-for="(model, index) in localModels" :key="model.name" class="model-row">
                <td>{{ index + 1 }}</td>
                <td class="model-name">
                  <div class="model-info">
                    <span class="name">{{ model.name }}</span>
                  </div>
                </td>
                <td>{{ formatSize(model.size) }}</td>
                <td>{{ formatDate(model.modified_at) }}</td>
                <td>
                  <span class="status-badge" :class="{ running: model.isRunning, stopped: !model.isRunning }">
                    <span class="status-dot"></span>
                    {{ model.isRunning ? t('common.running') : t('common.stopped') }}
                  </span>
                </td>
                <td>
                  <button class="btn-small primary" @click="viewModelDetails(model)">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
                      <path d="M1 12S5 4 12 4S23 12 23 12S19 20 12 20S1 12 1 12Z" stroke="currentColor" stroke-width="2"/>
                      <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
                    </svg>
                    {{ t('common.view') }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Download Model Drawer -->
    <el-drawer v-model="downloadDialogVisible" :title="t('modelManager.downloadModelTitle')" direction="rtl" size="40%" :close-on-click-modal="false" class="custom-drawer">
      <div class="info-box">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
          <path d="M12 16V12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M12 8H12.01" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <p>{{ t('modelManager.downloadInstruction') }} <a href="#" @click.prevent="openOllamaLibrary" class="link">https://ollama.com/library</a></p>
      </div>
      
      <div class="input-group">
        <label>{{ t('modelManager.modelName') }}</label>
        <el-input v-model="modelNameToDownload" :placeholder="t('modelManager.enterModelName')" />
      </div>
      
      <div class="drawer-footer">
        <el-button @click="downloadDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleDownload" :loading="isDownloading">
          {{ t('modelManager.startDownload') }}
        </el-button>
      </div>
    </el-drawer>

    <!-- Download Queue Drawer -->
    <div v-if="downloadQueueDrawerVisible" class="drawer-overlay">
      <div class="drawer-content" @click.stop>
        <div class="drawer-header">
          <h3>{{ t('modelManager.downloadQueue') }}</h3>
          <button class="close-btn" @click="downloadQueueDrawerVisible = false">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
              <line x1="18" y1="6" x2="6" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <line x1="6" y1="6" x2="18" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </button>
        </div>
        
        <div class="drawer-body">
          <div v-if="downloadQueue.length === 0" class="empty-state">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none">
              <path d="M16 4H18A2 2 0 0 1 20 6V20A2 2 0 0 1 18 22H6A2 2 0 0 1 4 20V6A2 2 0 0 1 6 4H8" stroke="currentColor" stroke-width="2"/>
              <rect x="8" y="2" width="8" height="4" rx="1" ry="1" stroke="currentColor" stroke-width="2"/>
            </svg>
            <p>{{ t('modelManager.noDownloads') }}</p>
          </div>
          
          <div v-for="item in downloadQueue" :key="item.model" class="download-item">
            <div class="download-info">
              <h4>{{ item.model }}</h4>
              <p>{{ item.status }}</p>
            </div>
            <div class="progress-container">
              <div class="progress-bar">
                <div class="progress-fill" :style="{ width: item.percentage + '%' }"></div>
              </div>
              <span class="progress-text">{{ item.percentage }}%</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Model Details Drawer -->
    <el-drawer
        v-model="drawerVisible"
        title="Model Details"
        direction="rtl"
        size="40%"
        :close-on-click-modal="false"
        class="custom-drawer"
    >
      <div v-if="selectedModel">
        <div>
          <h4>{{ t('modelManager.modelDetails') }}</h4>
          <el-descriptions :column="1" border size="small">
            <el-descriptions-item :label="t('modelManager.modelName')">{{ selectedModel.name }}</el-descriptions-item>
            <el-descriptions-item :label="t('modelManager.size')">{{ formatSize(selectedModel.size) }}</el-descriptions-item>
            <el-descriptions-item :label="t('modelManager.modified')">{{ formatDate(selectedModel.modified_at) }}</el-descriptions-item>
            <el-descriptions-item label="Digest">{{ selectedModel.digest }}</el-descriptions-item>
            <el-descriptions-item label="Model ID">{{ selectedModel.model }}</el-descriptions-item>
            <el-descriptions-item :label="t('modelManager.status')">
              <el-tag v-if="selectedModel.isRunning" type="success">{{ t('common.running') }}</el-tag>
              <el-tag v-else type="info">{{ t('common.stopped') }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="t('modelManager.modelFamily')" v-if="selectedModel.details && selectedModel.details.family">
              {{ selectedModel.details.family }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('modelManager.modelFormat')" v-if="selectedModel.details && selectedModel.details.format">
              {{ selectedModel.details.format }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('modelManager.parameterSize')" v-if="selectedModel.details && selectedModel.details.parameter_size">
              {{ selectedModel.details.parameter_size }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('modelManager.quantizationLevel')" v-if="selectedModel.details && selectedModel.details.quantization_level">
              {{ selectedModel.details.quantization_level }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('modelManager.parentModel')" v-if="selectedModel.details && selectedModel.details.parent_model">
              {{ selectedModel.details.parent_model }}
            </el-descriptions-item>
            <el-descriptions-item :label="t('modelManager.supportedFamilies')" v-if="selectedModel.details && selectedModel.details.families">
              {{ Array.isArray(selectedModel.details.families) ? selectedModel.details.families.join(', ') : selectedModel.details.families }}
            </el-descriptions-item>
          </el-descriptions>
          
          <div v-if="selectedModel.details" style="margin-top: 15px;">
            <h5>{{ t('modelManager.rawData') }}</h5>
            <el-card shadow="never" style="background: #f8f9fa;">
              <pre style="font-size: 12px; margin: 0; white-space: pre-wrap; word-break: break-all;">{{ JSON.stringify(selectedModel.details, null, 2) }}</pre>
            </el-card>
          </div>
        </div>

        <div style="margin-top: 20px">
          <el-button v-if="!selectedModel.isRunning" type="primary" @click="runModel" :loading="isRunningModel"
                     :disabled="isRunningModel">{{ t('common.start') }}
          </el-button>
          <el-button type="danger" @click="deleteModel(selectedModel)">{{ t('common.delete') }}</el-button>
          <el-button v-if="selectedModel.isRunning" @click="stopModel" :loading="isStoppingModel"
                     :disabled="isStoppingModel">{{ t('common.stop') }}
          </el-button>
        </div>

        <el-divider/>

        <div>
          <h4>{{ t('modelManager.modelParameters') }}</h4>
          <el-form :model="modelParams" label-width="80px" size="small">
            <el-form-item :label="t('modelManager.temperature')">
              <el-slider v-model="modelParams.temperature" :min="0" :max="1" :step="0.1"/>
            </el-form-item>
            <el-form-item label="Top P">
              <el-slider v-model="modelParams.topP" :min="0" :max="1" :step="0.1"/>
            </el-form-item>
            <el-form-item :label="t('modelManager.context')">
              <el-input-number v-model="modelParams.context" :min="1" :max="32768"/>
            </el-form-item>
          </el-form>
          <div style="margin-top: 10px">
            <el-button @click="saveModelParams" type="primary" size="small">{{ t('modelManager.saveParameters') }}</el-button>
            <el-button @click="resetModelParams" size="small">{{ t('modelManager.resetParameters') }}</el-button>
          </div>
        </div>

        <el-divider/>

        <div>
          <h4>{{ t('modelManager.modelTest') }}</h4>
          <el-form label-position="top">
            <el-form-item :label="t('modelManager.testContent')">
              <el-input
                  v-model="testPrompt"
                  type="textarea"
                  :rows="4"
                  :placeholder="t('modelManager.enterTestContent')"
              />
            </el-form-item>
          </el-form>
          <el-button @click="testModel" :loading="isTestingModel" type="primary" size="small">{{ t('modelManager.sendTest') }}</el-button>

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
        <p>{{ t('modelManager.selectModelDetails') }}</p>
      </div>
    </el-drawer>

  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, onUnmounted, reactive, ref} from 'vue'
import {ElMessage, ElMessageBox, ElNotification} from 'element-plus'
import { useI18n } from 'vue-i18n'
import * as runtime from '../../wailsjs/runtime/runtime.js'
import {
  DeleteModel,
  DownloadModel,
  GetActiveServer,
  GetModelParams,
  GetServers,
  ListModelsByServer,
  OpenInBrowser,
  RunModel,
  SetActiveServer,
  SetModelParams,
  StopModel,
  TestModel
} from '../../wailsjs/go/main/App'
import {types} from "../../wailsjs/go/models";
import {DownloadProgress, ModelParams} from "../classes/types";
import OllamaServerConfig = types.OllamaServerConfig;
import Model = types.Model;

const { t } = useI18n();

const openOllamaLibrary = () => {
  OpenInBrowser('https://ollama.com/library')
}

const localModels = ref<Model[]>([])
const selectedModel = ref<Model | null>(null)
const loading = ref(false)
const availableServers = ref<OllamaServerConfig[]>([])
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
  outputMode: 'stream',
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
  if (!dateString || dateString.trim() === '') {
    return t('common.unknown')
  }
  
  const date = new Date(dateString)
  if (isNaN(date.getTime())) {
    return t('common.unknown')
  }
  
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString()
}

const resetModelParams = () => {
  modelParams.temperature = 0.8
  modelParams.topP = 0.9
  modelParams.topK = 40
  modelParams.context = 2048
  modelParams.numPredict = 512
  modelParams.repeatPenalty = 1.1
  ElMessage.info(t('messages.parametersReset'))
}

const loadAvailableServers = async () => {
  try {
    availableServers.value = await GetServers();

    if (availableServers.value.length === 0) {
      ElMessage.warning('没有配置任何Ollama服务。请在“服务设置”页面添加一个。');
      localModels.value = [];
      selectedServer.value = '';
      return;
    }

    const activeServer = await GetActiveServer();
    const activeServerExists = activeServer && availableServers.value.some(s => s.id === activeServer.id);

    let serverToSelect = '';
    if (activeServerExists) {
      serverToSelect = activeServer.id;
    } else {
      serverToSelect = availableServers.value[0].id;
      // 如果默认的活动服务器不存在，则将列表中的第一个设置为活动状态
      await SetActiveServer(serverToSelect);
    }

    selectedServer.value = serverToSelect;

  } catch (error) {
    console.error('加载可用服务器列表失败:', error);
    ElMessage.error('加载可用服务器列表失败: ' + (error as Error).message);
    availableServers.value = [];
    selectedServer.value = '';
  }
};

const onServerChange = () => {
  getModels()
}

const getModels = async () => {
  if (!selectedServer.value) {
    localModels.value = [];
    loading.value = false;
    return;
  }
  try {
    loading.value = true
    localModels.value = await ListModelsByServer(selectedServer.value)
  } catch (error: any) {
    ElMessage.error(t('messages.getModelListFailed') + ': ' + error.message)
    localModels.value = [];
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
  ElMessage.info(`${t('messages.modelStarting')} "${selectedModel.value.name}"...`)
  try {
    await RunModel(selectedModel.value.name, modelParams)
    selectedModel.value.isRunning = true
    const index = localModels.value.findIndex(m => m.name === selectedModel.value!.name)
    if (index !== -1) localModels.value[index].isRunning = true
    ElMessage.success(`${t('messages.modelStarted')} "${selectedModel.value.name}"`)
  } catch (error: any) {
    ElMessage.error(t('messages.modelStartFailed') + ': ' + error.message)
    if (selectedModel.value) selectedModel.value.isRunning = false
  } finally {
    isRunningModel.value = false
  }
}

const stopModel = async () => {
  if (!selectedModel.value) return
  isStoppingModel.value = true
  ElMessage.info(`${t('messages.modelStopping')} "${selectedModel.value.name}"...`)
  try {
    await StopModel(selectedModel.value.name)
    selectedModel.value.isRunning = false
    const index = localModels.value.findIndex(m => m.name === selectedModel.value!.name)
    if (index !== -1) localModels.value[index].isRunning = false
    ElMessage.success(`${t('messages.modelStopped')} "${selectedModel.value.name}"`)
  } catch (error: any) {
    ElMessage.error(t('messages.modelStopFailed') + ': ' + error.message)
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
      ElMessage.success(t('messages.modelDeleted'))
      refreshModels()
      drawerVisible.value = false
    } catch (error: any) {
      ElMessage.error(t('messages.modelDeleteFailed') + ': ' + error.message)
    }
  }).catch(() => {
    ElMessage.info(t('messages.deleteCancelled'))
  })
}

const testModel = async () => {
  if (!selectedModel.value) return
  if (!testPrompt.value.trim()) {
    ElMessage.warning(t('messages.enterTestContent'))
    return
  }
  isTestingModel.value = true
  testResult.value = ''
  try {
    const response = await TestModel(selectedModel.value.name, testPrompt.value)
    testResult.value = response
    ElMessage.success(t('messages.testCompleted'))
  } catch (error: any) {
    testResult.value = t('messages.testFailed') + ': ' + error.message
    ElMessage.error(t('messages.testFailed'))
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
    ElMessage.warning(t('messages.enterModelName'))
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
    ElMessage.success(t('messages.parametersSaved'))
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

.main-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  flex: 1;
  min-height: 0;
}

.control-panel {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
}

.server-selector {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.server-selector label {
  font-weight: 600;
  color: #4a5568;
  font-size: 0.95rem;
}

.select-wrapper {
  position: relative;
  min-width: 250px;
}

.select-wrapper select {
  width: 100%;
  padding: 0.75rem 2.5rem 0.75rem 1rem;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 0.95rem;
  background: white;
  color: #2d3748;
  appearance: none;
  cursor: pointer;
  transition: all 0.3s ease;
}

.select-wrapper select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.select-arrow {
  position: absolute;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
  color: #a0aec0;
  pointer-events: none;
}

.action-buttons {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.25rem;
  border: none;
  border-radius: 8px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  text-decoration: none;
}

.btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.btn.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
}

.btn.secondary {
  background: white;
  color: #4a5568;
  border: 2px solid #e2e8f0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.btn.secondary:hover {
  background: #f7fafc;
  border-color: #cbd5e0;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.models-table-container {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.table-wrapper {
  overflow: auto;
  flex: 1;
}

.models-table {
  width: 100%;
  border-collapse: collapse;
}

.models-table th {
  background: #f8fafc;
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: #4a5568;
  border-bottom: 1px solid #e2e8f0;
  font-size: 0.9rem;
}

.models-table td {
  padding: 1rem;
  border-bottom: 1px solid #f1f5f9;
  color: #2d3748;
}

.model-row:hover {
  background: #f8fafc;
}

.model-name .name {
  font-weight: 600;
  color: #2d3748;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.375rem 0.75rem;
  border-radius: 50px;
  font-size: 0.875rem;
  font-weight: 500;
}

.status-badge.running {
  background: rgba(72, 187, 120, 0.1);
  color: #38a169;
  border: 1px solid rgba(72, 187, 120, 0.2);
}

.status-badge.stopped {
  background: rgba(160, 174, 192, 0.1);
  color: #718096;
  border: 1px solid rgba(160, 174, 192, 0.2);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.btn-small {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.5rem 0.875rem;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-small.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.btn-small.primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.empty-state {
  text-align: center;
  padding: 3rem;
}

.empty-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  color: #718096;
}

.empty-content svg {
  opacity: 0.5;
}

.empty-content p {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #4a5568;
}

.empty-content span {
  font-size: 0.95rem;
}

:deep(.custom-drawer .el-drawer__body) {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.info-box {
  display: flex;
  gap: 0.75rem;
  padding: 1rem;
  background: #ebf8ff;
  border: 1px solid #bee3f8;
  border-radius: 8px;
  margin-bottom: 1.5rem;
  color: #2b6cb0;
}

.info-box svg {
  flex-shrink: 0;
  margin-top: 0.125rem;
}

.info-box p {
  margin: 0;
  font-size: 0.9rem;
  line-height: 1.5;
}

.link {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.link:hover {
  text-decoration: underline;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 2rem;
}

.input-group label {
  font-weight: 500;
  color: #4a5568;
  font-size: 0.9rem;
}

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: auto;
  padding-top: 1rem;
  border-top: 1px solid #e2e8f0;
}



/* Drawer Styles */
.drawer-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  justify-content: flex-end;
  z-index: 1000;
}

.drawer-content {
  background: white;
  width: 400px;
  max-width: 90vw;
  height: 100vh;
  box-shadow: -10px 0 30px rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
}

.drawer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
}

.drawer-header h3 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #2d3748;
}

.drawer-body {
  flex: 1;
  padding: 1.5rem;
  overflow-y: auto;
}

.download-item {
  padding: 1rem;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  margin-bottom: 1rem;
}

.download-info h4 {
  margin: 0 0 0.5rem 0;
  font-size: 1rem;
  font-weight: 600;
  color: #2d3748;
}

.download-info p {
  margin: 0;
  font-size: 0.875rem;
  color: #718096;
}

.progress-container {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-top: 0.75rem;
}

.progress-bar {
  flex: 1;
  height: 8px;
  background: #e2e8f0;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 0.875rem;
  font-weight: 500;
  color: #4a5568;
  min-width: 40px;
  text-align: right;
}

:deep(.custom-drawer .el-drawer) {
  height: calc(100vh - 4rem) !important;
  top: 2rem !important;
}

@media (max-width: 768px) {
  .model-manager {
    padding: 1rem;
  }
  
  .control-panel {
    flex-direction: column;
    align-items: stretch;
  }
  
  .action-buttons {
    justify-content: center;
  }
  
  .models-table {
    font-size: 0.875rem;
  }
  
  .models-table th,
  .models-table td {
    padding: 0.75rem 0.5rem;
  }
  
  .drawer-content {
    width: 100vw;
  }
  
  :deep(.custom-drawer .el-drawer) {
    height: calc(100vh - 2rem) !important;
    top: 1rem !important;
  }
}
</style>
