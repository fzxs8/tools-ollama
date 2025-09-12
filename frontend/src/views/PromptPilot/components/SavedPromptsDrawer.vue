<template>
  <el-drawer
    :model-value="props.visible"
    title="我的提示词"
    direction="rtl"
    size="40%"
    @update:modelValue="closeDrawer"
  >
    <div class="saved-prompts-container">
      <el-empty v-if="!prompts || prompts.length === 0" description="没有提示词呢" />
      
      <div v-else class="prompts-list">
        <div 
          v-for="prompt in prompts" 
          :key="prompt.id"
          class="prompt-item"
        >
          <div class="prompt-header">
            <div class="prompt-title">{{ prompt.name }}</div>
            <div class="prompt-actions">
              <el-button 
                type="primary" 
                link 
                @click="previewPrompt(prompt)"
                size="small"
              >
                预览
              </el-button>
              <el-button 
                type="primary" 
                link 
                @click="editPrompt(prompt)"
                size="small"
              >
                编辑
              </el-button>
              <el-button 
                type="danger" 
                link 
                @click="deletePrompt(prompt)"
                size="small"
              >
                删除
              </el-button>
            </div>
          </div>
          
          <div class="prompt-meta">
            <div class="prompt-description" v-if="prompt.description">
              {{ prompt.description }}
            </div>
            <div class="prompt-tags" v-if="prompt.tags && prompt.tags.length > 0">
              <el-tag 
                v-for="tag in prompt.tags" 
                :key="tag" 
                size="small" 
                type="success"
                style="margin-right: 5px; margin-bottom: 5px;"
              >
                {{ tag }}
              </el-tag>
            </div>
            <div class="prompt-models">
              <el-tag 
                v-for="model in prompt.models" 
                :key="model" 
                type="info" 
                size="small"
                style="margin-right: 5px; margin-bottom: 5px;"
              >
                {{ model }}
              </el-tag>
            </div>
            <div class="prompt-time">
              更新于 {{ formatTime(prompt.updatedAt) }}
            </div>
          </div>
          
          <div class="prompt-content-preview-wrapper">
            <pre class="prompt-content-preview">{{ truncateContent(prompt.content, 100) }}</pre>
            <el-button class="preview-copy-btn" type="primary" link @click="copyPromptContent(prompt.content)">复制</el-button>
          </div>
        </div>
      </div>
    </div>
  </el-drawer>
</template>

<script setup lang="ts">
import {ElMessage} from "element-plus";

// 类型定义
type Prompt = main.Prompt;

const props = defineProps<{
  visible: boolean
  prompts: Prompt[]
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'preview', prompt: Prompt): void
  (e: 'edit', prompt: Prompt): void
  (e: 'delete', prompt: Prompt): void
}>()

const closeDrawer = () => {
  emit('update:visible', false);
}

// 格式化时间
const formatTime = (timestamp: number) => {
  const date = new Date(timestamp)
  return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')} ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

// 截断内容
const truncateContent = (content: string, length: number) => {
  if (!content) return '';
  return content.length > length ? content.substring(0, length) + '...' : content
}

const copyPromptContent = (content: string) => {
  if (content) {
    navigator.clipboard.writeText(content).then(() => {
      ElMessage.success('内容已复制到剪贴板')
    }).catch(err => {
      ElMessage.error('复制失败: ' + err)
    })
  }
}

// 预览Prompt
const previewPrompt = (prompt: Prompt) => {
  emit('preview', prompt)
}

// 编辑Prompt
const editPrompt = (prompt: Prompt) => {
  emit('edit', prompt)
}

// 删除Prompt
const deletePrompt = (prompt: Prompt) => {
  emit('delete', prompt)
}

</script>

<style scoped>
.saved-prompts-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.prompts-list {
  flex: 1;
  overflow-y: auto;
  padding-right: 10px; /* for scrollbar */
}

.prompt-item {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 15px;
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  transition: box-shadow 0.3s ease;
}

.prompt-item:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.prompt-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 10px;
}

.prompt-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  flex: 1;
  text-align: left;
}

.prompt-actions {
  display: flex;
  gap: 10px;
}

.prompt-meta {
  margin-bottom: 10px;
  text-align: left;
}

.prompt-description {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
  line-height: 1.4;
}

.prompt-tags {
  margin-bottom: 8px;
}

.prompt-models {
  margin-bottom: 8px;
}

.prompt-time {
  font-size: 12px;
  color: #999;
}

.prompt-content-preview-wrapper {
  position: relative;
}

pre.prompt-content-preview {
  font-size: 13px;
  color: #555;
  background-color: #f8f9fa;
  padding: 12px;
  border-radius: 6px;
  white-space: pre-wrap;
  line-height: 1.5;
  border: 1px solid #eee;
  margin: 0;
  text-align: left;
}

.preview-copy-btn {
  position: absolute;
  top: 5px;
  right: 5px; /* Changed from left to right */
  padding: 2px;
  font-size: 12px;
}
</style>