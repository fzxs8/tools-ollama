<template>
  <el-card class="model-selector">
    <template #header>
      <div class="card-header">
        <span>服务和模型选择</span>
      </div>
    </template>
    
    <!-- 服务选择 -->
    <div class="server-selector">
      <div class="section-title">选择服务</div>
      <el-select 
        v-model="selectedServerId" 
        placeholder="请选择服务" 
        style="width: 100%" 
        @change="onServerChange"
      >
        <el-option
          v-for="server in servers"
          :key="server.id"
          :label="server.name"
          :value="server.id"
        >
          <span style="float: left">{{ server.name }}</span>
          <span style="float: right; color: #8492a6; font-size: 13px">{{ server.baseUrl }}</span>
        </el-option>
      </el-select>
    </div>
    
    <!-- 模型选择 -->
    <div class="model-selector-section">
      <div class="section-title">选择模型 (最多3个)</div>
      <div v-if="isLoadingModels" class="loading-models">
        <el-skeleton :rows="3" animated />
      </div>
      <div v-else>
        <div 
          v-for="model in models" 
          :key="model.name"
          class="model-item"
          :class="{ selected: selectedModels.includes(model.name), disabled: selectedModels.length >= 3 && !selectedModels.includes(model.name) }"
          @click="toggleModelSelection(model.name)"
        >
          <div class="model-info">
            <div class="model-name">{{ model.name }}</div>
            <div class="model-details">
              <span class="model-size">{{ formatSize(model.size) }}</span>
              <span class="model-date">{{ formatDate(model.modified_at) }}</span>
            </div>
          </div>
          <div class="selection-indicator">
            <el-icon v-if="selectedModels.includes(model.name)"><Check /></el-icon>
          </div>
        </div>
      </div>
    </div>
    
    <div class="selection-summary">
      <el-tag 
        v-for="model in selectedModels" 
        :key="model" 
        closable
        @close="removeModel(model)"
        style="margin-right: 5px; margin-bottom: 5px;"
      >
        {{ model }}
      </el-tag>
    </div>
    
    <div class="model-limit-tip">
      最多可选择3个模型进行对比测试
    </div>
    
    <!-- 超过限制时的提醒对话框 -->
    <el-dialog
      v-model="showLimitDialog"
      title="提示"
      width="300px"
    >
      <span>最多只能选择3个模型进行对比测试</span>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="showLimitDialog = false">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Check } from '@element-plus/icons-vue'

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

const props = defineProps<{
  models: Model[]
  servers: Server[]
  selectedModels: string[]
  selectedServerId: string
  isLoadingModels: boolean
}>()

const emit = defineEmits<{
  (e: 'update:selectedModels', models: string[]): void
  (e: 'update:selectedServerId', serverId: string): void
  (e: 'loadModels', serverId: string): void
}>()

// 响应式数据
const showLimitDialog = ref(false)
const selectedServerId = ref(props.selectedServerId)

// 格式化文件大小
const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 格式化日期
const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

// 切换模型选择
const toggleModelSelection = (modelName: string) => {
  let newSelectedModels: string[]
  
  if (props.selectedModels.includes(modelName)) {
    // 如果已选择，则移除
    newSelectedModels = props.selectedModels.filter(name => name !== modelName)
  } else {
    // 如果未选择，检查是否已达到上限
    if (props.selectedModels.length >= 3) {
      // 达到上限时，显示提醒对话框
      showLimitDialog.value = true
      return
    }
    // 未达到上限，直接添加
    newSelectedModels = [...props.selectedModels, modelName]
  }
  
  emit('update:selectedModels', newSelectedModels)
}

// 移除模型
const removeModel = (modelName: string) => {
  const newSelectedModels = props.selectedModels.filter(name => name !== modelName)
  emit('update:selectedModels', newSelectedModels)
}

// 服务变更处理
const onServerChange = (serverId: string) => {
  selectedServerId.value = serverId
  emit('update:selectedServerId', serverId)
  emit('loadModels', serverId)
}

onMounted(() => {
  // 如果有选中的服务，加载模型
  if (props.selectedServerId) {
    emit('loadModels', props.selectedServerId)
  }
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

.section-title {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 10px;
  color: #606266;
}

.server-selector {
  margin-bottom: 20px;
}

.model-selector-section {
  margin-bottom: 20px;
}

.model-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 15px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  margin-bottom: 10px;
  cursor: pointer;
  transition: all 0.3s;
}

.model-item:hover {
  border-color: #409eff;
  background-color: #f5f9ff;
}

.model-item.selected {
  border-color: #409eff;
  background-color: #ecf5ff;
}

.model-item.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.model-item.disabled:hover {
  border-color: #e0e0e0;
  background-color: transparent;
}

.model-info {
  flex: 1;
}

.model-name {
  font-weight: 500;
  margin-bottom: 5px;
  color: #333;
}

.model-details {
  display: flex;
  font-size: 12px;
  color: #666;
}

.model-size {
  margin-right: 15px;
}

.selection-indicator {
  width: 20px;
  height: 20px;
  border: 1px solid #d0d0d0;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.model-item.selected .selection-indicator {
  background-color: #409eff;
  border-color: #409eff;
  color: white;
}

.selection-summary {
  min-height: 30px;
  margin-bottom: 15px;
}

.model-limit-tip {
  font-size: 12px;
  color: #999;
  text-align: center;
}

.loading-models {
  padding: 10px 0;
}
</style>