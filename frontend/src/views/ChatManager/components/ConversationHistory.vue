<template>
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
        @click="handleLoadConversation(conv)"
      >
        <div class="conversation-title">{{ conv.title }}</div>
        <div class="conversation-time">{{ formatTime(conv.timestamp) }}</div>
        <div class="conversation-actions">
          <el-button 
            size="small" 
            type="primary" 
            link
            @click.stop="handleEditTitle(conv)"
          >
            编辑
          </el-button>
          <el-button 
            size="small" 
            type="danger" 
            link
            @click.stop="handleDeleteConversation(conv.id)"
          >
            删除
          </el-button>
        </div>
      </div>
    </div>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessageBox } from 'element-plus'
import {types} from "../../../../wailsjs/go/models";
import Conversation = types.Conversation;

const props = defineProps<{
  visible: boolean
  conversations: Conversation[]
  activeConversationId: string
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'load-conversation', conversation: Conversation): void
  (e: 'edit-title', conversation: Conversation): void
  (e: 'delete-conversation', id: string): void
}>()

const visible = ref(false)

// 监听visible属性变化
watch(() => props.visible, (newVal) => {
  visible.value = newVal
})

// 监听visible值变化并更新父组件
watch(visible, (newVal) => {
  emit('update:visible', newVal)
})

// 格式化时间
const formatTime = (timestamp: number) => {
  const date = new Date(timestamp)
  return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')} ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

// 处理加载对话
const handleLoadConversation = (conv: Conversation) => {
  emit('load-conversation', conv)
}

// 处理编辑标题
const handleEditTitle = (conv: Conversation) => {
  emit('edit-title', conv)
}

// 处理删除对话
const handleDeleteConversation = async (id: string) => {
  try {
    await ElMessageBox.confirm('确定要删除这个对话吗？', '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    emit('delete-conversation', id)
  } catch (error) {
    // 用户取消操作
  }
}
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