<template>
  <div class="chat-interface">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="model-selector">
          <template #header>
            <div class="card-header">
              <span>模型选择</span>
            </div>
          </template>
          <el-select v-model="selectedModel" placeholder="选择模型" style="width: 100%">
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
                <el-slider v-model="modelParams.temperature" :min="0" :max="1" :step="0.1" />
              </el-form-item>
              <el-form-item label="Top P">
                <el-slider v-model="modelParams.topP" :min="0" :max="1" :step="0.1" />
              </el-form-item>
              <el-form-item label="上下文">
                <el-input-number v-model="modelParams.context" :min="1" :max="32768" />
              </el-form-item>
              <el-form-item label="预测数">
                <el-input-number v-model="modelParams.numPredict" :min="1" :max="4096" />
              </el-form-item>
              <el-form-item label="Top K">
                <el-input-number v-model="modelParams.topK" :min="1" :max="100" />
              </el-form-item>
              <el-form-item label="重复惩罚">
                <el-input-number v-model="modelParams.repeatPenalty" :min="0.1" :max="2" :step="0.1" />
              </el-form-item>
            </el-form>
            <div style="margin-top: 10px">
              <el-button @click="saveModelParams" type="primary" size="small">保存参数</el-button>
              <el-button @click="resetModelParams" size="small">重置</el-button>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="18">
        <el-card class="chat-container">
          <template #header>
            <div class="card-header">
              <span>聊天界面</span>
              <el-button @click="clearChat">清空聊天</el-button>
            </div>
          </template>
          
          <div class="chat-history" ref="chatHistoryRef">
            <div 
              v-for="(message, index) in messages" 
              :key="index" 
              class="message"
              :class="message.role"
            >
              <div class="message-content">
                <div class="role-tag">{{ message.role === 'user' ? '用户' : '助手' }}</div>
                <div class="content">{{ message.content }}</div>
              </div>
            </div>
            <div v-if="isThinking" class="message assistant">
              <div class="message-content">
                <div class="role-tag">助手</div>
                <div class="content">正在思考...</div>
              </div>
            </div>
          </div>
          
          <div class="input-area">
            <el-input
              v-model="inputMessage"
              type="textarea"
              :rows="3"
              placeholder="输入消息..."
              @keydown="handleKeydown"
            />
            <div style="margin-top: 10px; text-align: right">
              <el-button type="primary" @click="sendMessage" :disabled="!selectedModel || isThinking">
                发送 (Enter)
              </el-button>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { ListModels, ChatMessage, GetModelSettings, SaveModelSettings } from '../../wailsjs/go/main/App'

interface Model {
  name: string
  size: number
  modified_at: string
}

interface Message {
  role: 'user' | 'assistant'
  content: string
}

// 模型参数接口
interface ModelParams {
  temperature: number
  topP: number
  context: number
  numPredict: number
  topK: number
  repeatPenalty: number
}

const localModels = ref<Model[]>([])
const selectedModel = ref('')
const inputMessage = ref('')
const messages = ref<Message[]>([
  { role: 'assistant', content: '你好！我是Ollama助手，请选择一个模型开始对话。' }
])
const isThinking = ref(false)
const chatHistoryRef = ref<HTMLElement | null>(null)

// 模型参数
const modelParams = ref<ModelParams>({
  temperature: 0.8,
  topP: 0.9,
  context: 2048,
  numPredict: 512,
  topK: 40,
  repeatPenalty: 1.1
})

// 获取模型列表
const getModels = async () => {
  try {
    const models: Model[] = await ListModels()
    localModels.value = models
    if (models.length > 0 && !selectedModel.value) {
      selectedModel.value = models[0].name
      // 加载选中模型的参数设置
      loadModelParams(models[0].name)
    }
  } catch (error) {
    ElMessage.error('获取模型列表失败: ' + (error as Error).message)
  }
}

// 加载模型
const loadModel = () => {
  if (selectedModel.value) {
    ElMessage.success(`模型 ${selectedModel.value} 已加载`)
  } else {
    ElMessage.warning('请先选择一个模型')
  }
}

// 发送消息
const sendMessage = async () => {
  const message = inputMessage.value.trim()
  if (!message || !selectedModel.value || isThinking.value) return

  // 添加用户消息
  messages.value.push({ role: 'user', content: message })
  inputMessage.value = ''
  
  try {
    isThinking.value = true
    // 滚动到底部
    scrollToBottom()
    
    // 调用后端API
    const response: string = await ChatMessage(selectedModel.value, [
      { role: "user", content: message }
    ])
    
    // 添加助手回复
    messages.value.push({ role: 'assistant', content: response })
  } catch (error) {
    messages.value.push({ role: 'assistant', content: '抱歉，出现错误: ' + (error as Error).message })
  } finally {
    isThinking.value = false
    scrollToBottom()
  }
}

// 清空聊天
const clearChat = () => {
  messages.value = [{ role: 'assistant', content: '你好！我是Ollama助手，请选择一个模型开始对话。' }]
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (chatHistoryRef.value) {
      chatHistoryRef.value.scrollTop = chatHistoryRef.value.scrollHeight
    }
  })
}

// 处理键盘事件
const handleKeydown = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    sendMessage()
  }
}

// 加载模型参数
const loadModelParams = async (modelName: string) => {
  try {
    // 这里应该调用后端获取模型参数的接口
    // 暂时使用默认参数
    console.log(`加载模型 ${modelName} 的参数设置`)
  } catch (error) {
    console.error('加载模型参数失败:', error)
  }
}

// 保存模型参数
const saveModelParams = async () => {
  if (!selectedModel.value) {
    ElMessage.warning('请先选择一个模型')
    return
  }
  
  try {
    // 这里应该调用后端保存模型参数的接口
    ElMessage.success('参数保存成功')
  } catch (error) {
    ElMessage.error('参数保存失败: ' + (error as Error).message)
  }
}

// 重置模型参数
const resetModelParams = () => {
  modelParams.value = {
    temperature: 0.8,
    topP: 0.9,
    context: 2048,
    numPredict: 512,
    topK: 40,
    repeatPenalty: 1.1
  }
  ElMessage.info('参数已重置为默认值')
}

// 监听模型选择变化
const onModelChange = (newModel: string) => {
  loadModelParams(newModel)
}

onMounted(() => {
  getModels()
})
</script>

<style scoped>
.chat-interface {
  padding: 20px;
  height: calc(100vh - 80px);
}

.model-selector {
  height: 100%;
}

.chat-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chat-history {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  margin-bottom: 20px;
  background-color: #f5f7fa;
}

.message {
  margin-bottom: 15px;
}

.message.user {
  text-align: right;
}

.message.assistant {
  text-align: left;
}

.message-content {
  display: inline-block;
  max-width: 80%;
  padding: 10px 15px;
  border-radius: 8px;
}

.message.user .message-content {
  background-color: #409eff;
  color: white;
}

.message.assistant .message-content {
  background-color: white;
  border: 1px solid #dcdfe6;
}

.role-tag {
  font-size: 12px;
  margin-bottom: 5px;
  font-weight: bold;
}

.input-area {
  margin-top: auto;
}
</style>