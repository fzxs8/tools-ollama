<template>
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

    <slot name="input"></slot>
  </el-card>
</template>

<script setup lang="ts">
import {ref} from 'vue'
import MarkdownIt from 'markdown-it'

// 初始化Markdown解析器
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
})

interface Message {
  role: 'user' | 'assistant' | 'system'
  content: string
  timestamp?: number
}

interface SystemPrompt {
  id: string
  title: string
  prompt: string
  createdAt: number
}

const props = defineProps<{
  messages: Message[]
  isThinking: boolean
  activeSystemPrompt: SystemPrompt | null
}>()

const emit = defineEmits<{
  (e: 'clear-chat'): void
  (e: 'open-system-prompt'): void
  (e: 'copy-message', content: string): void
  (e: 'regenerate-message', index: number): void
}>()

const chatHistoryRef = ref<HTMLElement | null>(null)

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
  emit('copy-message', content)
}

// 打开系统提示词抽屉
const openSystemPromptDrawer = () => {
  emit('open-system-prompt')
}

// 清空聊天
const clearChat = () => {
  emit('clear-chat')
}

// 重新生成消息
const regenerateMessage = (index: number) => {
  emit('regenerate-message', index)
}
</script>

<style scoped>
.chat-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  border: none;
  box-shadow: none;
  width: 100%;
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
  padding: 12px 24px;
  border-bottom: 1px solid #e0e0e0;
  background-color: #ffffff;
}

.chat-history {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  background-color: #f0f4f9;
  display: flex;
  flex-direction: column;
}

.chat-message {
  display: flex;
  align-items: flex-start;
  margin-bottom: 24px;
  max-width: 85%;
  width: -moz-fit-content;
  width: fit-content;
}

.chat-message:last-child {
  margin-bottom: 0;
}

.user-message {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.assistant-message {
  align-self: flex-start;
}

.message-content-area {
  padding: 12px 16px;
  border-radius: 12px;
}

.user-message .message-content-area {
  background-color: #cce5ff; /* A softer, more professional blue */
}

.assistant-message .message-content-area {
  background-color: #ffffff;
}

.message-avatar {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-avatar {
  background-color: #4a5568;
}

.assistant-avatar {
  background: linear-gradient(135deg, #4285f4, #9b59b6);
}

.user-message .message-avatar {
  margin-left: 16px;
}

.assistant-message .message-avatar {
  margin-right: 16px;
}

.message-header {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.sender-name {
  font-weight: 600;
  font-size: 15px;
  color: #2d3748;
}

.message-time {
  font-size: 12px;
  color: #718096;
  margin-left: 12px;
}

.message-content {
  /* No flex: 1 here to ensure height is adaptive */
}

.message-body {
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.7;
  font-size: 16px;
  color: #1a202c;
  text-align: left !important;
}

/* Reset margin for the first element rendered by v-html */
.message-body :deep(> *:first-child) {
  margin-top: 0;
}

/* Reset margin for the last element rendered by v-html */
.message-body :deep(> *:last-child) {
  margin-bottom: 0;
}

.message-body :deep(p) {
  margin: 12px 0;
}

.message-body :deep(h1),
.message-body :deep(h2),
.message-body :deep(h3) {
  margin: 20px 0 12px 0;
  font-weight: 600;
  color: #2d3748;
}

.message-body :deep(h1) { font-size: 26px; }
.message-body :deep(h2) { font-size: 22px; }
.message-body :deep(h3) { font-size: 18px; }

.message-body :deep(ul),
.message-body :deep(ol) {
  padding-left: 24px;
  margin: 12px 0;
}

.message-body :deep(li) {
  margin: 8px 0;
}

.message-body :deep(code) {
  background-color: #edf2f7;
  padding: 3px 6px;
  border-radius: 6px;
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, Courier, monospace;
  font-size: 14px;
  color: #2d3748;
}

.message-body :deep(pre) {
  background-color: #1a202c; /* Dark background for code blocks */
  color: #e2e8f0;
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
  border-left: 4px solid #a0aec0;
  padding-left: 16px;
  margin: 16px 0;
  color: #4a5568;
}

.message-actions {
  display: flex;
  justify-content: flex-start;
  margin-top: 12px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.chat-message:hover .message-actions {
  opacity: 1;
}

.message-actions .el-button {
  font-size: 13px;
  padding: 5px 10px;
  margin-right: 10px;
  color: #718096;
}

.input-area {
  padding: 24px;
  border-top: 1px solid #e0e0e0;
  background-color: #ffffff;
}

.thinking-indicator {
  display: flex;
  align-items: center;
  font-size: 16px;
  color: #4a5568;
}

.thinking-indicator .dot {
  width: 8px;
  height: 8px;
  background-color: #a0aec0;
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
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-5px); }
}
</style>