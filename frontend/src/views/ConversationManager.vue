<template>
  <div class="conversation-manager">
    <el-drawer
        v-model="visible"
        title="历史对话"
        direction="rtl"
        size="400px"
    >
      <div class="conversation-list">
        <el-empty v-if="conversations.length === 0" description="暂无历史对话" />
        <div 
          v-else 
          v-for="conv in conversations" 
          :key="conv.id" 
          class="conversation-item"
          :class="{ active: conv.id === activeConversationId }"
          @click="loadConversation(conv)"
        >
          <div class="conversation-title">{{ conv.title }}</div>
          <div class="conversation-time">{{ formatTime(conv.timestamp) }}</div>
          <div class="conversation-actions">
            <el-button 
              size="small" 
              type="primary" 
              link
              @click.stop="editConversationTitle(conv)"
            >
              编辑
            </el-button>
            <el-button 
              size="small" 
              type="danger" 
              link
              @click.stop="deleteConversation(conv.id)"
            >
              删除
            </el-button>
          </div>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, defineProps, defineEmits } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  ListConversations,
  GetConversation,
  DeleteConversation,
  SaveConversation
} from '../../wailsjs/go/main/App'

interface Message {
  role: 'user' | 'assistant' | 'system'
  content: string
  timestamp?: number
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

const visible = defineModel<boolean>('visible', { default: false })
const activeConversationId = defineModel<string>('activeConversationId')
const conversations = ref<Conversation[]>([])

const emit = defineEmits(['load-conversation'])

// 格式化时间
const formatTime = (timestamp: number) => {
  const date = new Date(timestamp)
  return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')} ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

// 加载对话列表
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
    // 通过事件将对话数据传递给父组件
    emit('load-conversation', conversation)
    activeConversationId.value = conv.id
    visible.value = false
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
      activeConversationId.value = ''
    }
  } catch (error) {
    if (!(error as Error).message?.includes('cancel')) {
      ElMessage.error('删除对话失败: ' + (error as Error).message)
    }
  }
}

onMounted(() => {
  loadConversations()
})

defineExpose({
  loadConversations
})
</script>

<style scoped>
.conversation-list {
  height: 100%;
  overflow-y: auto;
}

.conversation-item {
  padding: 15px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  transition: background-color 0.2s;
}

.conversation-item:hover {
  background-color: #f5f5f5;
}

.conversation-item.active {
  background-color: #e6f7ff;
  border-left: 3px solid #1890ff;
}

.conversation-title {
  font-weight: 500;
  margin-bottom: 5px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.conversation-time {
  font-size: 12px;
  color: #999;
  margin-bottom: 10px;
}

.conversation-actions {
  display: flex;
  gap: 10px;
}
</style>