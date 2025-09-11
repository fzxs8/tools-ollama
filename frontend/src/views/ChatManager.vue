<template>
  <div class="chat-interface">
    <el-row :gutter="20" style="height: 100%;">
      <el-col :span="6" style="height: 100%;">
        <ModelSelector
          v-model:selectedModel="selectedModel"
          v-model:selectedServer="selectedServer"
          v-model:modelParams="modelParams"
          :localModels="localModels"
          :availableServers="availableServers"
          @load-model="loadModel"
          @save-model-params="saveModelParams"
          @reset-model-params="resetModelParams"
          @server-change="onServerChange"
        />
      </el-col>

      <el-col :span="18" style="height: 100%; display: flex; flex-direction: column;">
        <ChatContainer
          :messages="messages"
          :is-thinking="isThinking"
          :active-system-prompt="activeSystemPrompt"
          :conversations="conversations"
          :active-conversation-id="activeConversationId"
          @clear-chat="clearChat"
          @open-system-prompt="openSystemPromptDrawer"
          @copy-message="copyMessage"
          @regenerate-message="regenerateMessage"
          @new-conversation="newConversation"
          @load-conversation="loadConversation"
          @edit-conversation-title="editConversationTitle"
          @delete-conversation="deleteConversation"
        >
          <template #input>
            <ChatInput
              v-model="inputMessage"
              :disabled="isThinking"
              @send="sendMessage"
              @keydown="handleKeydown"
            />
          </template>
        </ChatContainer>
      </el-col>
    </el-row>

    <PromptListDrawer
      v-model:visible="systemPromptDrawerVisible"
      :prompts="systemPromptList"
      mode="select"
      :selected-id="activeSystemPrompt?.id"
      @select="handleApplySystemPrompt"
    />
  </div>
</template>

<script setup lang="ts">
import {onMounted, ref} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {
  ChatMessage,
  DeleteConversation,
  GetActiveServer,
  GetConversation,
  GetOllamaServerConfig,
  GetRemoteServers,
  ListConversations,
  ListModelsByServer,
  ListPrompts,
  SaveConversation
} from '../../wailsjs/go/main/App'
import {EventsOff, EventsOn} from '../../wailsjs/runtime'
import MarkdownIt from 'markdown-it'
import ChatInput from "./ChatManager/components/ChatInput.vue";
import ModelSelector from "./ChatManager/components/ModelSelector.vue";
import ChatContainer from "./ChatManager/components/ChatContainer.vue";
import PromptListDrawer from "../components/commons/PromptListDrawer.vue";
import {main} from "../../wailsjs/go/models";

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
})

type Prompt = main.Prompt;

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
  role: 'user' | 'assistant' | 'system'
  content: string
  timestamp?: number
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

interface Conversation {
  id: string
  title: string
  messages: Message[]
  modelName: string
  systemPrompt: string
  modelParams: string
  timestamp: number
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
const systemPromptDrawerVisible = ref(false)
const activeSystemPrompt = ref<Prompt | null>(null)
const systemPromptList = ref<Prompt[]>([])
const conversations = ref<Conversation[]>([])
const activeConversationId = ref('')
const currentConversation = ref<Conversation | null>(null)

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
    systemPromptList.value = await ListPrompts();
  } catch (error) {
    console.error('加载系统提示词失败:', error)
    systemPromptList.value = []
  }
}

// 应用系统提示词
const handleApplySystemPrompt = (prompt: Prompt) => {
  activeSystemPrompt.value = prompt;
  ElMessage.success(`已应用系统提示词: “${prompt.name}”`);
  systemPromptDrawerVisible.value = false;
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
    let messagesWithSystemPrompt: Message[] = [
      { role: "user", content: message, timestamp: Date.now() }
    ]

    // 如果有激活的系统提示词，添加到消息历史最前面
    if (activeSystemPrompt.value) {
      messagesWithSystemPrompt.unshift({
        role: "system",
        content: activeSystemPrompt.value.content,
        timestamp: Date.now()
      })
    }

    // 根据输出方式选择不同的处理方式
    if (modelParams.value.outputMode === 'stream') {
      // 流式输出
      const assistantMessageIndex = messages.value.length
      messages.value.push({
        role: 'assistant',
        content: '',
        timestamp: Date.now()
      })

      // 定义事件监听器
      const streamListener = (chunk: string) => {
        try {
          if (messages.value && messages.value[assistantMessageIndex]) {
            messages.value[assistantMessageIndex].content += chunk
            setTimeout(() => scrollToBottom(), 0)
          }
        } catch (e) {
          console.error('更新消息时出错:', e)
        }
      }

      // 注册事件监听
      EventsOn('chat_stream_chunk', streamListener)

      try {
        // 调用后端API进行流式传输，现在第三个参数是布尔值
        await ChatMessage(selectedModel.value, messagesWithSystemPrompt, true)
      } finally {
        // 确保无论成功还是失败都注销监听器
        EventsOff('chat_stream_chunk', streamListener)
      }

    } else {
      // 阻塞输出
      const assistantMessageIndex = messages.value.length
      messages.value.push({
        role: 'assistant',
        content: '',
        timestamp: Date.now()
      })
      
      try {
        // 调用后端API，第三个参数为false
        const response: string = await ChatMessage(selectedModel.value, messagesWithSystemPrompt, false)
        if (messages.value && messages.value[assistantMessageIndex]) {
          messages.value[assistantMessageIndex].content = response
          setTimeout(() => scrollToBottom(), 0)
        }
      } catch (error) {
        console.error('阻塞式输出发生错误:', error)
        throw error
      }
    }
    
    // 保存对话
    await saveCurrentConversation()
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

// 新建对话
const newConversation = () => {
  // 清空消息列表
  messages.value = [{
    role: 'assistant',
    content: '你好！我是Ollama助手，请选择一个模型开始对话。',
    timestamp: Date.now()
  }]
  
  // 清空当前对话ID和对话对象
  activeConversationId.value = ''
  currentConversation.value = null
}

// 加载对话历史列表
const loadConversations = async () => {
  try {
    conversations.value = await ListConversations()
  } catch (error) {
    ElMessage.error('加载对话列表失败: ' + (error as Error).message)
  }
}

// 加载对话
const loadConversation = async (conv: Conversation) => {
  try {
    const conversation = await GetConversation(conv.id)
    // 加载对话到界面
    messages.value = conversation.messages
    selectedModel.value = conversation.modelName
    activeConversationId.value = conversation.id
    currentConversation.value = conversation
  } catch (error) {
    ElMessage.error('加载对话失败: ' + (error as Error).message)
  }
}

// 编辑对话标题
const editConversationTitle = async (conv: Conversation) => {
  try {
    const newTitle = await ElMessageBox.prompt('请输入新的对话标题', '编辑标题', {
      inputValue: conv.title,
      confirmButtonText: '确定',
      cancelButtonText: '取消'
    })
    
    if (newTitle.value) {
      const updatedConv = { ...conv, title: newTitle.value }
      await SaveConversation(updatedConv)
      ElMessage.success('标题更新成功')
      loadConversations() // 重新加载列表
    }
  } catch (error) {
    // 用户取消操作
  }
}

// 删除对话
const deleteConversation = async (id: string) => {
  try {
    await ElMessageBox.confirm('确定要删除这个对话吗？', '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await DeleteConversation(id)
    ElMessage.success('对话删除成功')
    loadConversations() // 重新加载列表
    
    // 如果删除的是当前激活的对话，需要清除激活状态
    if (activeConversationId.value === id) {
      newConversation() // 创建新对话
    }
  } catch (error) {
    if (!(error as Error).message?.includes('cancel')) {
      ElMessage.error('删除对话失败: ' + (error as Error).message)
    }
  }
}

// 保存当前对话
const saveCurrentConversation = async () => {
  try {
    // 构建对话对象
    const conversation: any = {
      messages: messages.value,
      modelName: selectedModel.value,
      systemPrompt: activeSystemPrompt.value ? JSON.stringify(activeSystemPrompt.value) : '',
      modelParams: JSON.stringify(modelParams.value),
      timestamp: Date.now()
    }
    
    // 如果是新对话，生成标题
    if (!activeConversationId.value) {
      // 使用第一条用户消息作为标题
      const firstUserMessage = messages.value.find(m => m.role === 'user')
      conversation.title = firstUserMessage ? firstUserMessage.content.substring(0, 20) + (firstUserMessage.content.length > 20 ? '...' : '') : '新对话'
    } else {
      // 如果是已有对话，获取原对话信息
      conversation.id = activeConversationId.value
      conversation.title = currentConversation.value?.title || '新对话'
    }
    
    // 保存对话
    const savedConversation = await SaveConversation(conversation)
    activeConversationId.value = savedConversation.id
    currentConversation.value = savedConversation
    
    // 重新加载对话列表
    loadConversations()
  } catch (error) {
    console.error('保存对话失败:', error)
  }
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
      let messagesWithSystemPrompt: Message[] = [
        {role: "user", content: userMessage, timestamp: Date.now()}
      ]

      // 如果有激活的系统提示词，添加到消息历史最前面
      if (activeSystemPrompt.value) {
        messagesWithSystemPrompt.unshift({
          role: "system",
          content: activeSystemPrompt.value.content,
          timestamp: Date.now()
        })
      }

      // 添加一个空的助手消息用于更新
      const assistantMessageIndex = messages.value.length
      messages.value.push({
        role: 'assistant',
        content: '',
        timestamp: Date.now()
      })

      // 重新生成时，我们默认使用阻塞模式以简化逻辑
      const response: string = await ChatMessage(selectedModel.value, messagesWithSystemPrompt, false)

      // 更新助手消息
      if (messages.value && messages.value[assistantMessageIndex]) {
        messages.value[assistantMessageIndex].content = response
      }
      
      // 保存对话
      await saveCurrentConversation()
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
  // 这里需要实现滚动到底部的逻辑
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
    repeatPenalty: 1.1,
    outputMode: 'stream'
  }
  ElMessage.info('参数已重置为默认值')
}

onMounted(async () => {
  await loadAvailableServers()
  await getModels()
  await loadSystemPrompts()
  await loadConversations()
  
  // 监听流式传输事件
  if (window && (window as any).runtime) {
    (window as any).runtime.EventsOn("chat_stream_chunk", (data: any) => {
      // 更新最后一条助手消息
      const lastMessageIndex = messages.value.length - 1
      if (lastMessageIndex >= 0 && messages.value[lastMessageIndex].role === 'assistant') {
        messages.value[lastMessageIndex].content += data
        // 滚动到底部
        scrollToBottom()
      }
    })
  }
})
</script>

<style scoped>
.chat-interface {
  background-color: #f0f4f9;
  padding: 0;
  height: 100%;
  box-sizing: border box;
}
</style>