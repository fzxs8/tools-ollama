<template>
  <div class="chat-interface">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-icon">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M21 15A2 2 0 0 1 19 17H7L4 20V5A2 2 0 0 1 6 3H19A2 2 0 0 1 21 5Z" stroke="currentColor"
                  stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <div class="header-text">
          <h1>{{ $t('chatManager.title') }}</h1>
          <p>{{ $t('chatManager.description') }}</p>
        </div>
      </div>
    </div>

    <!-- Main Chat Layout -->
    <div class="chat-layout">
      <!-- Sidebar -->
      <div class="chat-sidebar">
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
      </div>

      <!-- Chat Area -->
      <div class="chat-main">
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
      </div>
    </div>

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
import {useI18n} from 'vue-i18n'
import {
  ChatMessage,
  DeleteConversation,
  GetActiveServer,
  GetConversation,
  GetServers,
  ListConversations,
  ListModelsByServer,
  ListPrompts,
  SaveConversation,
  SetActiveServer
} from '../../../wailsjs/go/main/App'
import {EventsOn} from '../../../wailsjs/runtime'
import ChatInput from "./components/ChatInput.vue";
import ModelSelector from "./components/ModelSelector.vue";
import ChatContainer from "./components/ChatContainer.vue";
import PromptListDrawer from "../../components/commons/PromptListDrawer.vue";
import {types} from "../../../wailsjs/go/models";
import {ModelParams} from "../../classes/types";
import Conversation = types.Conversation;
import Message = types.Message;
import OllamaServerConfig = types.OllamaServerConfig;
import Model = types.Model;
import Prompt = types.Prompt;

const {t} = useI18n();

const localModels = ref<Model[]>([])
const selectedModel = ref('')
const availableServers = ref<OllamaServerConfig[]>([])
const selectedServer = ref('')
const inputMessage = ref('')
const messages = ref<Message[]>([
  {role: 'assistant', content: t('chatManager.welcomeMessage'), timestamp: Date.now()}
])
const isThinking = ref(false)
const systemPromptDrawerVisible = ref(false)
const activeSystemPrompt = ref<Prompt | null>(null)
const systemPromptList = ref<Prompt[]>([])
const conversations = ref<Conversation[]>([])
const activeConversationId = ref('')
const currentConversation = ref<Conversation | null>(null)

// Model parameters
const modelParams = ref<ModelParams>({
  temperature: 0.8,
  topP: 0.9,
  context: 2048,
  numPredict: 512,
  topK: 40,
  repeatPenalty: 1.1,
  outputMode: 'stream' // Default to stream output
})

// 复制消息内容
const copyMessage = (content: string) => {
  navigator.clipboard.writeText(content).then(() => {
    ElMessage.success(t('chatManager.messageCopied'))
  }).catch(() => {
    ElMessage.error(t('chatManager.copyFailed'))
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
    console.error(t('chatManager.loadSystemPromptsFailed'), error)
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
    availableServers.value = await GetServers();

    if (availableServers.value.length === 0) {
      ElMessage.warning(t('chatManager.noOllamaServices'));
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
      await SetActiveServer(serverToSelect);
    }
    selectedServer.value = serverToSelect;

  } catch (error) {
    console.error(t('chatManager.loadServersFailed'), error);
    ElMessage.error(t('chatManager.loadServersFailed') + ': ' + (error as Error).message);
    availableServers.value = [];
    selectedServer.value = '';
  }
};

const onServerChange = async () => {
  if (!selectedServer.value) {
    localModels.value = [];
    return;
  }
  try {
    await SetActiveServer(selectedServer.value);
    await getModels();
  } catch (error: any) {
    ElMessage.error(t('chatManager.switchServiceFailed') + ': ' + error.message);
  }
}

// 获取模型列表
const getModels = async () => {
  if (!selectedServer.value) {
    localModels.value = [];
    return;
  }
  try {
    localModels.value = await ListModelsByServer(selectedServer.value)
    if (localModels.value.length > 0) {
      selectedModel.value = localModels.value[0].name
      loadModelParams(localModels.value[0].name)
    } else {
      selectedModel.value = ''
    }
  } catch (error: any) {
    ElMessage.error(t('chatManager.getModelListFailed') + ': ' + error.message)
    localModels.value = [] // 清空模型列表
  }
}

// 加载模型
const loadModel = () => {
  if (selectedModel.value) {
    ElMessage.success(`${t('chatManager.modelLoaded')} ${selectedModel.value}`)
  } else {
    ElMessage.warning(t('chatManager.selectModelFirst'))
  }
}

// 发送消息
const sendMessage = async () => {
  const message = inputMessage.value.trim()
  console.log('sendMessage called with:', {message, isThinking: isThinking.value, selectedModel: selectedModel.value})

  if (!message || isThinking.value) return

  if (!selectedModel.value) {
    messages.value.push({
      role: 'assistant',
      content: t('chatManager.selectModelFirst'),
      timestamp: Date.now()
    })
    return
  }

  messages.value.push({
    role: 'user',
    content: message,
    timestamp: Date.now()
  })
  inputMessage.value = ''

  try {
    isThinking.value = true
    scrollToBottom()

    let messagesWithSystemPrompt: Message[] = [...messages.value];

    if (activeSystemPrompt.value) {
      messagesWithSystemPrompt.unshift({
        role: "system",
        content: activeSystemPrompt.value.content,
        timestamp: Date.now()
      })
    }

    if (modelParams.value.outputMode === 'stream') {
      messages.value.push({
        role: 'assistant',
        content: '',
        timestamp: Date.now()
      })

      try {
        console.log('Calling ChatMessage with:', {
          model: selectedModel.value,
          messageCount: messagesWithSystemPrompt.length,
          stream: true
        })
        const result = await ChatMessage(selectedModel.value, messagesWithSystemPrompt, true)
        console.log('ChatMessage result:', result)
      } catch (error) {
        console.error('ChatMessage error:', error)
        const lastMessageIndex = messages.value.length - 1
        if (lastMessageIndex >= 0 && messages.value[lastMessageIndex].role === 'assistant') {
          messages.value[lastMessageIndex].content = t('chatManager.sorry') + ': ' + (error as Error).message
        }
      }

    } else {
      const assistantMessageIndex = messages.value.length
      messages.value.push({
        role: 'assistant',
        content: '',
        timestamp: Date.now()
      })

      try {
        console.log('Calling ChatMessage (blocking) with:', {
          model: selectedModel.value,
          messageCount: messagesWithSystemPrompt.length,
          stream: false
        })
        const response: string = await ChatMessage(selectedModel.value, messagesWithSystemPrompt, false)
        console.log('ChatMessage (blocking) result:', response)
        if (messages.value && messages.value[assistantMessageIndex]) {
          messages.value[assistantMessageIndex].content = response
          setTimeout(() => scrollToBottom(), 0)
        }
      } catch (error) {
        console.error('ChatMessage (blocking) error:', error)
        throw error
      }
    }

    await saveCurrentConversation()
  } catch (error: any) {
    console.error(t('chatManager.sendMessageError'), error)
    let errorMessage = t('chatManager.sorry')
    if (error && error.message) {
      errorMessage += ': ' + error.message
    } else if (error && typeof error === 'string') {
      errorMessage += ': ' + error
    } else if (error && error.toString && error.toString() !== '[object Object]') {
      errorMessage += ': ' + error.toString()
    } else {
      errorMessage += ': ' + t('chatManager.unknownError')
    }

    // 确保最后一条消息是助手消息，如果不是则添加一条
    const lastMessageIndex = messages.value.length - 1
    if (lastMessageIndex >= 0 && messages.value[lastMessageIndex].role === 'assistant') {
      messages.value[lastMessageIndex].content = errorMessage
    } else {
      messages.value.push({
        role: 'assistant',
        content: errorMessage,
        timestamp: Date.now()
      })
    }
  } finally {
    isThinking.value = false
    scrollToBottom()
  }
}

// 清空聊天
const clearChat = () => {
  messages.value = [{
    role: 'assistant',
    content: t('chatManager.welcomeMessage'),
    timestamp: Date.now()
  }]
}

// 新建对话
const newConversation = () => {
  messages.value = [{
    role: 'assistant',
    content: t('chatManager.welcomeMessage'),
    timestamp: Date.now()
  }]
  activeConversationId.value = ''
  currentConversation.value = null
}

// 加载对话历史列表
const loadConversations = async () => {
  try {
    conversations.value = await ListConversations()
  } catch (error) {
    ElMessage.error(t('chatManager.loadConversationsFailed') + ': ' + (error as Error).message)
  }
}

// 加载对话
const loadConversation = async (conv: Conversation) => {
  try {
    const conversation = await GetConversation(conv.id)
    messages.value = conversation.messages
    selectedModel.value = conversation.modelName
    activeConversationId.value = conversation.id
    currentConversation.value = conversation
  } catch (error) {
    ElMessage.error(t('chatManager.loadConversationFailed') + ': ' + (error as Error).message)
  }
}

// 编辑对话标题
const editConversationTitle = async (conv: Conversation) => {
  try {
    const newTitle = await ElMessageBox.prompt(t('chatManager.enterNewTitle'), t('chatManager.editTitle'), {
      inputValue: conv.title,
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel')
    })

    if (newTitle.value) {
      const updatedConv = {...conv, title: newTitle.value}
      await SaveConversation(updatedConv)
      ElMessage.success(t('chatManager.titleUpdated'))
      await loadConversations() // 重新加载列表
    }
  } catch (error) {
    // 用户取消操作
  }
}

// 删除对话
const deleteConversation = async (id: string) => {
  try {
    await DeleteConversation(id)
    ElMessage.success(t('chatManager.conversationDeleted'))
    await loadConversations() // 重新加载列表

    if (activeConversationId.value === id) {
      newConversation() // 创建新对话
    }
  } catch (error) {
    ElMessage.error(t('chatManager.deleteConversationFailed') + ': ' + (error as Error).message)
  }
}

// 保存当前对话
const saveCurrentConversation = async () => {
  try {
    const conversation: any = {
      messages: messages.value,
      modelName: selectedModel.value,
      systemPrompt: activeSystemPrompt.value ? JSON.stringify(activeSystemPrompt.value) : '',
      modelParams: JSON.stringify(modelParams.value),
      timestamp: Date.now()
    }

    if (!activeConversationId.value) {
      const firstUserMessage = messages.value.find(m => m.role === 'user')
      conversation.title = firstUserMessage ? firstUserMessage.content.substring(0, 20) + (firstUserMessage.content.length > 20 ? '...' : '') : '新对话'
    } else {
      conversation.id = activeConversationId.value
      conversation.title = currentConversation.value?.title || '新对话'
    }

    const savedConversation = await SaveConversation(conversation)
    activeConversationId.value = savedConversation.id
    currentConversation.value = savedConversation

    await loadConversations()
  } catch (error) {
    console.error('保存对话失败:', error)
  }
}

// 重新生成消息
const regenerateMessage = async (index: number) => {
  if (index > 0 && messages.value[index].role === 'assistant' && messages.value[index - 1].role === 'user') {
    messages.value.splice(index, 1);

    try {
      isThinking.value = true;
      scrollToBottom();

      let messagesWithSystemPrompt: Message[] = [...messages.value];

      if (activeSystemPrompt.value) {
        messagesWithSystemPrompt.unshift({
          role: "system",
          content: activeSystemPrompt.value.content,
          timestamp: Date.now()
        })
      }

      const assistantMessageIndex = messages.value.length
      messages.value.push({
        role: 'assistant',
        content: '',
        timestamp: Date.now()
      })

      if (modelParams.value.outputMode === 'stream') {
        await ChatMessage(selectedModel.value, messagesWithSystemPrompt, true)
      } else {
        const response: string = await ChatMessage(selectedModel.value, messagesWithSystemPrompt, false)
        if (messages.value && messages.value[assistantMessageIndex]) {
          messages.value[assistantMessageIndex].content = response
        }
      }

      await saveCurrentConversation()
    } catch (error) {
      const lastMessageIndex = messages.value.length - 1
      if (lastMessageIndex >= 0 && messages.value[lastMessageIndex].role === 'assistant') {
        messages.value[lastMessageIndex].content = t('chatManager.sorry') + ': ' + (error as Error).message
      } else {
        messages.value.push({
          role: 'assistant',
          content: t('chatManager.sorry') + ': ' + (error as Error).message,
          timestamp: Date.now()
        })
      }
    } finally {
      isThinking.value = false;
      scrollToBottom();
    }
  }
}

// 滚动到底部
const scrollToBottom = () => {
  setTimeout(() => {
    const chatHistory = document.querySelector('.chat-history')
    if (chatHistory) {
      chatHistory.scrollTop = chatHistory.scrollHeight
    }
  }, 100)
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
  await loadAvailableServers();
  // Ensure backend state is synchronized before loading models or performing other operations
  if (selectedServer.value) {
    try {
      await SetActiveServer(selectedServer.value);
    } catch (error: any) {
      ElMessage.error(`同步活动服务器失败: ${error.message}`);
      return; // Don't continue if sync fails
    }
  }

  await getModels();
  await loadSystemPrompts();
  await loadConversations();

  // Listen for streaming events
  EventsOn("chat_stream_chunk", (data: any) => {
    const lastMessageIndex = messages.value.length - 1;
    if (lastMessageIndex >= 0 && messages.value[lastMessageIndex].role === 'assistant') {
      messages.value[lastMessageIndex].content += data;
      scrollToBottom();
    }
  });

  // Listen for streaming completion events
  EventsOn("chat_stream_done", () => {
    isThinking.value = false;
    scrollToBottom();
  });

  // Listen for streaming error events
  EventsOn("chat_stream_error", (error: string) => {
    isThinking.value = false;
    const lastMessageIndex = messages.value.length - 1;
    if (lastMessageIndex >= 0 && messages.value[lastMessageIndex].role === 'assistant') {
      if (messages.value[lastMessageIndex].content === '') {
        messages.value[lastMessageIndex].content = t('chatManager.sorry') + ': ' + error;
      } else {
        messages.value[lastMessageIndex].content += '\n\n' + t('chatManager.sorry') + ': ' + error;
      }
    }
    scrollToBottom();
  });
});
</script>

<style scoped>
.chat-interface {
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

.chat-layout {
  flex: 1;
  display: flex;
  width: 100%;
  padding: 0 2rem 2rem 2rem;
  gap: 1.5rem;
  min-height: 0;
}

.chat-sidebar {
  width: 350px;
  flex-shrink: 0;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
}

.chat-main {
  flex: 1;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
}

@media (max-width: 1200px) {
  .chat-layout {
    flex-direction: column;
    padding: 0;
  }

  .chat-sidebar {
    width: 100%;
    height: auto;
  }

  .chat-interface {
    padding: 1.5rem;
  }

  .header-text h1 {
    font-size: 1.75rem;
  }
}

@media (max-width: 768px) {
  .chat-interface {
    padding: 1rem;
  }

  .header-content {
    flex-direction: column;
    text-align: center;
  }

  .header-text h1 {
    font-size: 1.5rem;
  }

  .header-text p {
    font-size: 1rem;
  }
}
</style>
