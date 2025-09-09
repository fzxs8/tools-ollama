<template>
  <div class="chat-interface">
    <el-row :gutter="20" style="height: 100%;">
      <el-col :span="6" style="height: 100%;">
        <el-card class="model-selector">
          <template #header>
            <div class="card-header">
              <span>模型选择</span>
            </div>
          </template>
          <el-select v-model="selectedServer" placeholder="选择服务" style="width: 100%; margin-bottom: 10px;"
                     @change="onServerChange">
            <el-option
                v-for="server in availableServers"
                :key="server.id"
                :label="server.name"
                :value="server.id"
            />
          </el-select>
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
                <el-slider v-model="modelParams.temperature" :min="0" :max="1" :step="0.1"/>
              </el-form-item>
              <el-form-item label="Top P">
                <el-slider v-model="modelParams.topP" :min="0" :max="1" :step="0.1"/>
              </el-form-item>
              <el-form-item label="上下文">
                <el-input-number v-model="modelParams.context" :min="1" :max="32768"/>
              </el-form-item>
              <el-form-item label="预测数">
                <el-input-number v-model="modelParams.numPredict" :min="1" :max="4096"/>
              </el-form-item>
              <el-form-item label="Top K">
                <el-input-number v-model="modelParams.topK" :min="1" :max="100"/>
              </el-form-item>
              <el-form-item label="重复惩罚">
                <el-input-number v-model="modelParams.repeatPenalty" :min="0.1" :max="2" :step="0.1"/>
              </el-form-item>
              <el-form-item label="输出方式">
                <el-select v-model="modelParams.outputMode" placeholder="选择输出方式">
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
      </el-col>

      <el-col :span="18" style="height: 100%; display: flex; flex-direction: column;">
        <el-card class="chat-container" style="flex: 1; display: flex; flex-direction: column;">
          <template #header>
            <div class="card-header">
              <span>聊天界面{{ activeSystemPrompt ? ` - ${activeSystemPrompt.title}` : '' }}</span>
              <div>
                <el-button @click="clearChat" style="margin-right: 10px;">清空聊天</el-button>
                <el-button @click="openSystemPromptDrawer">系统提示词</el-button>
              </div>
            </div>
          </template>

          <div class="chat-history" ref="chatHistoryRef">
            <div
                v-for="(message, index) in messages"
                :key="index"
                class="chat-message"
                :class="{ 'user-message': message.role === 'user', 'assistant-message': message.role === 'assistant' }"
            >
              <div class="message-avatar">
                <div v-if="message.role === 'user'" class="user-avatar">
                  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path
                        d="M12 12C14.2091 12 16 10.2091 16 8C16 5.79086 14.2091 4 12 4C9.79086 4 8 5.79086 8 8C8 10.2091 9.79086 12 12 12Z"
                        fill="white"/>
                    <path d="M6 20C6 12 12 12 12 12C12 12 18 12 18 20" stroke="white" stroke-width="2"
                          stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </div>
                <div v-else class="assistant-avatar">
                  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="white" stroke-width="2" stroke-linecap="round"
                          stroke-linejoin="round"/>
                    <path d="M2 17L12 22L22 17" stroke="white" stroke-width="2" stroke-linecap="round"
                          stroke-linejoin="round"/>
                    <path d="M2 12L12 17L22 12" stroke="white" stroke-width="2" stroke-linecap="round"
                          stroke-linejoin="round"/>
                  </svg>
                </div>
              </div>
              <div class="message-content-area">
                <div class="message-header">
                  <span class="sender-name">{{ message.role === 'user' ? 'You' : 'Assistant' }}</span>
                  <span class="message-time">{{ formatTime(message.timestamp || Date.now()) }}</span>
                </div>
                <div class="message-content">
                  <div class="message-body" v-html="renderMarkdown(message.content)"></div>
                  <div class="message-actions" v-if="message.role === 'assistant'">
                    <el-button
                        size="small"
                        type="primary"
                        @click="copyMessage(message.content)"
                        link
                    >
                      复制
                    </el-button>
                    <el-button
                        size="small"
                        type="primary"
                        @click="regenerateMessage(index)"
                        link
                    >
                      重新生成
                    </el-button>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="isThinking" class="chat-message assistant-message">
              <div class="message-avatar">
                <div class="assistant-avatar">
                  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="white" stroke-width="2" stroke-linecap="round"
                          stroke-linejoin="round"/>
                    <path d="M2 17L12 22L22 17" stroke="white" stroke-width="2" stroke-linecap="round"
                          stroke-linejoin="round"/>
                    <path d="M2 12L12 17L22 12" stroke="white" stroke-width="2" stroke-linecap="round"
                          stroke-linejoin="round"/>
                  </svg>
                </div>
              </div>
              <div class="message-content-area">
                <div class="message-header">
                  <span class="sender-name">Assistant</span>
                  <span class="message-time">{{ getCurrentTime() }}</span>
                </div>
                <div class="message-content">
                  <div class="thinking-indicator">
                    <span>正在思考</span>
                    <div class="dot"></div>
                    <div class="dot"></div>
                    <div class="dot"></div>
                  </div>
                </div>
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
              <el-button type="primary" @click="sendMessage" :disabled="isThinking">
                发送 (Enter)
              </el-button>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 系统提示词设置抽屉 -->
    <el-drawer
        v-model="systemPromptDrawerVisible"
        title="系统提示词设置"
        direction="rtl"
        size="40%"
    >
      <div class="system-prompt-content">
        <h3>添加/编辑提示词</h3>
        <el-form :model="systemPromptForm" label-width="80px">
          <el-form-item label="标题">
            <el-input
                v-model="systemPromptForm.title"
                placeholder="请输入提示词标题"
            />
          </el-form-item>
          <el-form-item label="提示词">
            <el-input
                v-model="systemPromptForm.prompt"
                type="textarea"
                :rows="6"
                placeholder="请输入系统提示词，用于指导AI助手的行为"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary"
                       @click="systemPromptForm.title && systemPromptList.some(p => p.title === systemPromptForm.title) ? updateSystemPrompt() : saveSystemPrompt()">
              保存
            </el-button>
            <el-button @click="resetSystemPromptForm">重置</el-button>
          </el-form-item>
        </el-form>

        <h3>提示词列表</h3>
        <el-table :data="systemPromptList" style="width: 100%" empty-text="暂无提示词">
          <el-table-column prop="title" label="标题"/>
          <el-table-column prop="createdAt" label="创建时间">
            <template #default="scope">
              {{ new Date(scope.row.createdAt).toLocaleString() }}
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template #default="scope">
              <el-button
                  size="small"
                  type="primary"
                  @click="setActiveSystemPrompt(scope.row)"
                  :disabled="activeSystemPrompt && activeSystemPrompt.id === scope.row.id"
              >
                {{ activeSystemPrompt && activeSystemPrompt.id === scope.row.id ? '已启用' : '启用' }}
              </el-button>
              <el-button size="small" @click="editSystemPrompt(scope.row)">修改</el-button>
              <el-button size="small" type="danger" @click="deleteSystemPrompt(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="prompt-examples" style="margin-top: 30px;">
          <h4>提示词示例</h4>
          <el-collapse>
            <el-collapse-item title="通用助手" name="1">
              <p>你是一个有用的助手，请以简洁明了的方式回答用户的问题。</p>
            </el-collapse-item>
            <el-collapse-item title="技术专家" name="2">
              <p>你是一位技术专家，请以专业的角度回答用户的技术问题，并提供详细的解释。</p>
            </el-collapse-item>
            <el-collapse-item title="创意写作" name="3">
              <p>你是一位创意写作专家，请帮助用户创作富有想象力和吸引力的内容。</p>
            </el-collapse-item>
          </el-collapse>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import {nextTick, onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {
  ChatMessage,
  GetActiveServer,
  GetOllamaServerConfig,
  GetRemoteServers,
  KVDelete,
  KVGet,
  KVList,
  KVSet,
  ListModelsByServer
} from '../../wailsjs/go/main/App'
import MarkdownIt from 'markdown-it'

// 初始化Markdown解析器
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
})

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

interface Message {
  role: 'user' | 'assistant'
  content: string
  timestamp?: number
}

// 系统提示词接口
interface SystemPrompt {
  id: string
  title: string
  prompt: string
  createdAt: number
}

// 模型参数接口
interface ModelParams {
  temperature: number
  topP: number
  context: number
  numPredict: number
  topK: number
  repeatPenalty: number
  outputMode: 'stream' | 'blocking' // 添加输出方式选项
}

// 系统提示词表单
interface SystemPromptForm {
  title: string
  prompt: string
}

const localModels = ref<Model[]>([])
const selectedModel = ref('')
const availableServers = ref<Server[]>([])
const selectedServer = ref<string>('local')
const inputMessage = ref('')
const messages = ref<Message[]>([
  {role: 'assistant', content: '你好！我是Ollama助手，请选择一个模型开始对话。', timestamp: Date.now()}
])
const isThinking = ref(false)
const chatHistoryRef = ref<HTMLElement | null>(null)
const systemPromptDrawerVisible = ref(false)
const activeSystemPrompt = ref<SystemPrompt | null>(null)
const systemPromptList = ref<SystemPrompt[]>([])

// 系统提示词表单
const systemPromptForm = reactive<SystemPromptForm>({
  title: '',
  prompt: ''
})

// 模型参数
const modelParams = ref<ModelParams>({
  temperature: 0.8,
  topP: 0.9,
  context: 2048,
  numPredict: 512,
  topK: 40,
  repeatPenalty: 1.1,
  outputMode: 'stream' // 默认使用流式输出
})

// 格式化时间
const formatTime = (timestamp: number) => {
  const date = new Date(timestamp)
  const hours = date.getHours().toString().padStart(2, '0')
  const minutes = date.getMinutes().toString().padStart(2, '0')
  return `${hours}:${minutes}`
}

// 获取当前时间
const getCurrentTime = () => {
  const now = new Date()
  const hours = now.getHours().toString().padStart(2, '0')
  const minutes = now.getMinutes().toString().padStart(2, '0')
  return `${hours}:${minutes}`
}

// 渲染Markdown内容
const renderMarkdown = (content: string) => {
  return md.render(content)
}

// 复制消息内容
const copyMessage = (content: string) => {
  navigator.clipboard.writeText(content).then(() => {
    ElMessage.success('消息已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

// 打开系统提示词抽屉
const openSystemPromptDrawer = async () => {
  systemPromptDrawerVisible.value = true
  await loadSystemPrompts()
}

// 加载系统提示词列表
const loadSystemPrompts = async () => {
  try {
    const promptListStr = await KVList("system_prompts")
    if (promptListStr) {
      systemPromptList.value = JSON.parse(promptListStr)
    } else {
      systemPromptList.value = []
    }

    // 加载当前激活的提示词
    const activePromptStr = await KVGet("active_system_prompt")
    if (activePromptStr) {
      activeSystemPrompt.value = JSON.parse(activePromptStr)
    }
  } catch (error) {
    console.error('加载系统提示词失败:', error)
    systemPromptList.value = []
  }
}

// 保存系统提示词
const saveSystemPrompt = async () => {
  if (!systemPromptForm.title.trim()) {
    ElMessage.warning('请输入提示词标题')
    return
  }

  if (!systemPromptForm.prompt.trim()) {
    ElMessage.warning('请输入提示词内容')
    return
  }

  try {
    const newPrompt: SystemPrompt = {
      id: Date.now().toString(),
      title: systemPromptForm.title,
      prompt: systemPromptForm.prompt,
      createdAt: Date.now()
    }

    // 添加到列表
    systemPromptList.value.push(newPrompt)

    // 保存到存储
    await KVSet("system_prompts", JSON.stringify(systemPromptList.value))

    ElMessage.success('系统提示词已保存')
    resetSystemPromptForm()
  } catch (error) {
    ElMessage.error('保存失败: ' + (error as Error).message)
  }
}

// 重置系统提示词表单
const resetSystemPromptForm = () => {
  systemPromptForm.title = ''
  systemPromptForm.prompt = ''
}

// 设置激活的系统提示词
const setActiveSystemPrompt = async (prompt: SystemPrompt) => {
  try {
    activeSystemPrompt.value = prompt
    await KVSet("active_system_prompt", JSON.stringify(prompt))
    ElMessage.success(`已启用提示词: ${prompt.title}`)
  } catch (error) {
    ElMessage.error('设置失败: ' + (error as Error).message)
  }
}

// 删除系统提示词
const deleteSystemPrompt = async (prompt: SystemPrompt) => {
  try {
    // 从列表中移除
    systemPromptList.value = systemPromptList.value.filter(p => p.id !== prompt.id)

    // 如果删除的是当前激活的提示词，则清除激活状态
    if (activeSystemPrompt.value?.id === prompt.id) {
      activeSystemPrompt.value = null
      await KVDelete("active_system_prompt")
    }

    // 保存到存储
    await KVSet("system_prompts", JSON.stringify(systemPromptList.value))

    ElMessage.success('提示词已删除')
  } catch (error) {
    ElMessage.error('删除失败: ' + (error as Error).message)
  }
}

// 编辑系统提示词
const editSystemPrompt = (prompt: SystemPrompt) => {
  systemPromptForm.title = prompt.title
  systemPromptForm.prompt = prompt.prompt
}

// 更新系统提示词
const updateSystemPrompt = async () => {
  if (!systemPromptForm.title.trim()) {
    ElMessage.warning('请输入提示词标题')
    return
  }

  if (!systemPromptForm.prompt.trim()) {
    ElMessage.warning('请输入提示词内容')
    return
  }

  try {
    // 查找并更新提示词
    const index = systemPromptList.value.findIndex(p =>
        p.title === systemPromptForm.title && p.prompt !== systemPromptForm.prompt)

    if (index !== -1) {
      systemPromptList.value[index].prompt = systemPromptForm.prompt
      systemPromptList.value[index].title = systemPromptForm.title

      // 如果更新的是当前激活的提示词，则同步更新
      if (activeSystemPrompt.value?.id === systemPromptList.value[index].id) {
        activeSystemPrompt.value = {...systemPromptList.value[index]}
        await KVSet("active_system_prompt", JSON.stringify(activeSystemPrompt.value))
      }

      // 保存到存储
      await KVSet("system_prompts", JSON.stringify(systemPromptList.value))

      ElMessage.success('系统提示词已更新')
      resetSystemPromptForm()
    }
  } catch (error) {
    ElMessage.error('更新失败: ' + (error as Error).message)
  }
}

const loadAvailableServers = async () => {
  try {
    const localBaseUrl = await GetOllamaServerConfig()
    const localServer = {
      id: 'local',
      name: '本地服务',
      baseUrl: localBaseUrl,
      apiKey: '',
      isActive: true,
      testStatus: '',
      type: 'local'
    }
    let remoteServers: Server[] = []
    try {
      const remoteList: any[] = await GetRemoteServers()
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
    availableServers.value = [localServer, ...remoteServers]

    try {
      const activeServer = await GetActiveServer()
      const activeServerExists = activeServer && activeServer.id && availableServers.value.some(s => s.id === activeServer.id)
      if (activeServerExists) {
        selectedServer.value = activeServer.id
      } else {
        selectedServer.value = 'local'
      }
    } catch (e) {
      selectedServer.value = 'local'
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
    availableServers.value = [localServer];
    selectedServer.value = 'local'
  }
}

const onServerChange = () => {
  getModels()
}

// 获取模型列表
const getModels = async () => {
  try {
    const models: Model[] = await ListModelsByServer(selectedServer.value)
    localModels.value = models
    if (models.length > 0) {
      selectedModel.value = models[0].name
      // 加载选中模型的参数设置
      loadModelParams(models[0].name)
    } else {
      selectedModel.value = ''
    }
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
    localModels.value = [] // 清空模型列表
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
  if (!message || isThinking.value) return

  // 检查是否选择了模型
  if (!selectedModel.value) {
    messages.value.push({
      role: 'assistant',
      content: '请先选择一个模型',
      timestamp: Date.now()
    })
    return
  }

  // 添加用户消息
  messages.value.push({
    role: 'user',
    content: message,
    timestamp: Date.now()
  })
  inputMessage.value = ''

  try {
    isThinking.value = true
    // 滚动到底部
    scrollToBottom()

    // 构建包含系统提示词的消息
    let messagesWithSystemPrompt = [
      { role: "user", content: message }
    ]

    // 如果有激活的系统提示词，添加到消息历史最前面
    if (activeSystemPrompt.value) {
      messagesWithSystemPrompt.unshift({
        role: "system",
        content: activeSystemPrompt.value.prompt
      })
    }

    // 根据输出方式选择不同的处理方式
    if (modelParams.value.outputMode === 'stream') {
      // 流式输出
      // 添加一个空的助手消息用于流式更新
      const assistantMessageIndex = messages.value.length
      messages.value.push({
        role: 'assistant',
        content: '',
        timestamp: Date.now()
      })

      try {
        // 调用后端API进行流式传输
        await ChatMessage(selectedModel.value, messagesWithSystemPrompt, (chunk: string) => {
          // 流式更新助手消息
          try {
            console.log('流式输出接收到数据块，长度:', chunk.length)
            console.log('流式输出接收到数据块内容:', chunk)
            // 确保消息仍然存在并且可以安全更新
            if (messages.value && messages.value[assistantMessageIndex]) {
              messages.value[assistantMessageIndex].content += chunk
              // 使用setTimeout确保DOM更新不会阻塞
              setTimeout(() => scrollToBottom(), 0)
            }
          } catch (e) {
            console.error('更新消息时出错:', e)
          }
        })
      } catch (error) {
        console.error('流式输出发生错误:', error)
        throw error
      }
    } else {
      // 阻塞输出
      // 添加一个空的助手消息用于更新
      const assistantMessageIndex = messages.value.length
      messages.value.push({
        role: 'assistant',
        content: '',
        timestamp: Date.now()
      })
      
      try {
        const response: string = await ChatMessage(selectedModel.value, messagesWithSystemPrompt, null)
        console.log('阻塞式输出接收到响应，长度:', response.length)
        console.log('阻塞式输出接收到响应，前100个字符:', response.substring(0, Math.min(100, response.length)))
        // 更新助手消息
        if (messages.value && messages.value[assistantMessageIndex]) {
          messages.value[assistantMessageIndex].content = response
          // 确保滚动到底部以显示新内容
          setTimeout(() => scrollToBottom(), 0)
        }
      } catch (error) {
        console.error('阻塞式输出发生错误:', error)
        throw error // 重新抛出错误以在下面的catch块中处理
      }
    }
  } catch (error: any) {
    console.error('发送消息时出现错误:', error)
    let errorMessage = '抱歉，出现错误'
    if (error.message) {
      errorMessage += ': ' + error.message
    } else if (error.toString() !== '[object Object]') {
      errorMessage += ': ' + error.toString()
    } else {
      errorMessage += ': 未知错误'
    }
    // 添加错误消息
    messages.value.push({
      role: 'assistant',
      content: errorMessage,
      timestamp: Date.now()
    })
  } finally {
    isThinking.value = false
    scrollToBottom()
  }
}

// 清空聊天
const clearChat = () => {
  messages.value = [{
    role: 'assistant',
    content: '你好！我是Ollama助手，请选择一个模型开始对话。',
    timestamp: Date.now()
  }]
}

// 重新生成消息
const regenerateMessage = async (index: number) => {
  // 确保这是助手消息，并且有前一条用户消息
  if (index > 0 && messages.value[index].role === 'assistant' && messages.value[index - 1].role === 'user') {
    // 获取前一条用户消息
    const userMessage = messages.value[index - 1].content;

    // 移除当前助手消息
    messages.value.splice(index, 1);

    try {
      isThinking.value = true;
      // 滚动到底部
      scrollToBottom();

      // 构建包含系统提示词的消息
      let messagesWithSystemPrompt = [
        {role: "user", content: userMessage}
      ]

      // 如果有激活的系统提示词，添加到消息历史最前面
      if (activeSystemPrompt.value) {
        messagesWithSystemPrompt.unshift({
          role: "system",
          content: activeSystemPrompt.value.prompt
        })
      }

      // 添加一个空的助手消息用于更新
      const assistantMessageIndex = messages.value.length
      messages.value.push({
        role: 'assistant',
        content: '',
        timestamp: Date.now()
      })

      // 调用后端API
      const response: string = await ChatMessage(selectedModel.value, messagesWithSystemPrompt, null)

      // 更新助手消息
      if (messages.value && messages.value[assistantMessageIndex]) {
        messages.value[assistantMessageIndex].content = response
      }
    } catch (error) {
      messages.value.push({
        role: 'assistant',
        content: '抱歉，重新生成消息时出现错误: ' + (error as Error).message,
        timestamp: Date.now()
      })
    } finally {
      isThinking.value = false;
      scrollToBottom();
    }
  }
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

onMounted(async () => {
  await loadAvailableServers()
  await getModels()
  await loadSystemPrompts()
})
</script>

<style scoped>
.chat-interface {
  padding: 20px;
  height: 100%;
  box-sizing: border-box;
}

.model-selector {
  height: 100%;
}

.model-selector :deep(.el-card__body) {
  height: calc(100% - 60px);
  overflow-y: auto;
}

.chat-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.chat-container :deep(.el-card__body) {
  height: calc(100% - 60px);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid #e5e5e5;
}

.chat-history {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  background-color: #f7f7f8;
}

.chat-message {
  display: flex;
  padding: 24px 0;
  border-bottom: 1px solid #e5e5e5;
}

.chat-message:last-child {
  border-bottom: none;
}

.user-message {
  background-color: #ffffff;
}

.assistant-message {
  background-color: #f7f7f8;
}

.message-avatar {
  flex-shrink: 0;
  width: 30px;
  height: 30px;
  border-radius: 2px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 20px;
}

.user-avatar {
  background-color: #5436da;
  border-radius: 2px;
}

.assistant-avatar {
  background-color: #19c37d;
  border-radius: 2px;
}

.message-content-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding-right: 30px;
}

.message-header {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.sender-name {
  font-weight: 600;
  font-size: 14px;
  color: #333;
  margin-right: 12px;
}

.message-time {
  font-size: 12px;
  color: #888;
}

.message-content {
  flex: 1;
}

.message-body {
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.6;
  font-size: 16px;
  color: #333;
  text-align: left;
}

.message-body :deep(p) {
  margin: 12px 0;
}

.message-body :deep(h1),
.message-body :deep(h2),
.message-body :deep(h3) {
  margin: 18px 0 12px 0;
  font-weight: 600;
}

.message-body :deep(h1) {
  font-size: 28px;
}

.message-body :deep(h2) {
  font-size: 24px;
}

.message-body :deep(h3) {
  font-size: 20px;
}

.message-body :deep(ul),
.message-body :deep(ol) {
  padding-left: 24px;
  margin: 12px 0;
}

.message-body :deep(li) {
  margin: 6px 0;
}

.message-body :deep(code) {
  background-color: rgba(0, 0, 0, 0.05);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 14px;
}

.message-body :deep(pre) {
  background-color: #2d2d2d;
  color: #f8f8f2;
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 16px 0;
}

.message-body :deep(pre code) {
  background-color: transparent;
  padding: 0;
  color: inherit;
  font-size: 14px;
}

.message-body :deep(blockquote) {
  border-left: 4px solid #19c37d;
  padding-left: 16px;
  margin: 12px 0;
  color: #666;
}

.message-body :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 16px 0;
}

.message-body :deep(th),
.message-body :deep(td) {
  border: 1px solid #ddd;
  padding: 10px 14px;
  text-align: left;
}

.message-body :deep(th) {
  background-color: #f5f5f5;
  font-weight: 600;
}

.message-body :deep(img) {
  max-width: 100%;
  height: auto;
  margin: 12px 0;
}

.message-body :deep(a) {
  color: #19c37d;
  text-decoration: none;
}

.message-body :deep(a:hover) {
  text-decoration: underline;
}

.thinking-indicator {
  display: flex;
  align-items: center;
  font-size: 16px;
}

.thinking-indicator .dot {
  width: 8px;
  height: 8px;
  background-color: #999;
  border-radius: 50%;
  margin: 0 4px;
  animation: bounce 1.5s infinite;
}

.thinking-indicator .dot:nth-child(2) {
  animation-delay: 0.2s;
}

.thinking-indicator .dot:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes bounce {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-5px);
  }
}

.message-actions {
  display: flex;
  justify-content: flex-start;
  margin-top: 12px;
}

.message-actions .el-button {
  font-size: 13px;
  padding: 5px 10px;
  margin-right: 10px;
}

.input-area {
  padding: 24px;
  border-top: 1px solid #e5e5e5;
  background-color: white;
}

.system-prompt-content {
  padding: 20px;
  height: 100%;
  box-sizing: border-box;
}

.prompt-examples {
  margin-top: 30px;
}

.prompt-examples h4 {
  margin-bottom: 15px;
}
</style>