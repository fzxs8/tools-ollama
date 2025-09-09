<template>
  <el-drawer
    v-model="visible"
    title="已保存的Prompt"
    direction="rtl"
    size="40%"
  >
    <div class="saved-prompts-container">
      <el-empty v-if="prompts.length === 0" description="暂无保存的Prompt" />
      
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
                @click="deletePrompt(prompt.id)"
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
              创建于 {{ formatTime(prompt.createdAt) }}
            </div>
          </div>
          
          <div class="prompt-content-preview">
            {{ truncateContent(prompt.content, 100) }}
          </div>
        </div>
      </div>
    </div>
  </el-drawer>
  
  <!-- 编辑Prompt对话框 -->
  <el-dialog
    v-model="showEditDialog"
    :title="editingPrompt ? '编辑Prompt' : '新增Prompt'"
    width="600px"
  >
    <el-form :model="promptForm" label-width="80px">
      <el-form-item label="名称">
        <el-input v-model="promptForm.name" placeholder="请输入Prompt名称" />
      </el-form-item>
      <el-form-item label="描述">
        <el-input 
          v-model="promptForm.description" 
          type="textarea" 
          placeholder="请输入描述（可选）" 
          :rows="2" />
      </el-form-item>
      <el-form-item label="标签">
        <el-select
          v-model="promptForm.tags"
          multiple
          filterable
          allow-create
          default-first-option
          placeholder="请输入标签，可创建新标签"
          style="width: 100%">
          <el-option
            v-for="tag in allTags"
            :key="tag"
            :label="tag"
            :value="tag">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="内容">
        <el-input
          v-model="promptForm.content"
          type="textarea"
          placeholder="请输入Prompt内容"
          :rows="6"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="saveEditedPrompt">保存</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'

interface Prompt {
  id: string
  name: string
  content: string
  description: string
  createdAt: number
  updatedAt: number
  models: string[]
  version: number
  tags: string[]
  createdBy: string
}

const props = defineProps<{
  visible: boolean
  prompts: Prompt[]
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'preview', prompt: Prompt): void
  (e: 'edit', prompt: Prompt): void
  (e: 'delete', id: string): void
  (e: 'save', prompt: Prompt): void
}>()

// 响应式数据
const visible = ref(false)
const showEditDialog = ref(false)
const editingPrompt = ref<Prompt | null>(null)

const promptForm = ref({
  name: '',
  description: '',
  tags: [] as string[],
  content: ''
})

// 计算所有标签的集合
const allTags = computed(() => {
  const tagsSet = new Set<string>()
  props.prompts.forEach(prompt => {
    prompt.tags.forEach(tag => tagsSet.add(tag))
  })
  return Array.from(tagsSet)
})

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

// 截断内容
const truncateContent = (content: string, length: number) => {
  return content.length > length ? content.substring(0, length) + '...' : content
}

// 预览Prompt
const previewPrompt = (prompt: Prompt) => {
  emit('preview', prompt)
}

// 编辑Prompt
const editPrompt = (prompt: Prompt) => {
  editingPrompt.value = prompt
  promptForm.value = {
    name: prompt.name,
    description: prompt.description,
    tags: [...prompt.tags],
    content: prompt.content
  }
  showEditDialog.value = true
}

// 删除Prompt
const deletePrompt = (id: string) => {
  emit('delete', id)
}

// 保存编辑的Prompt
const saveEditedPrompt = () => {
  if (!promptForm.value.name.trim()) {
    ElMessage.warning('请输入Prompt名称')
    return
  }
  
  if (!promptForm.value.content.trim()) {
    ElMessage.warning('请输入Prompt内容')
    return
  }
  
  const updatedPrompt: Prompt = {
    id: editingPrompt.value?.id || Date.now().toString(),
    name: promptForm.value.name,
    description: promptForm.value.description,
    tags: [...promptForm.value.tags],
    content: promptForm.value.content,
    models: editingPrompt.value?.models || [],
    createdAt: editingPrompt.value?.createdAt || Date.now(),
    updatedAt: Date.now(),
    version: editingPrompt.value?.version ? editingPrompt.value.version + 1 : 1,
    createdBy: editingPrompt.value?.createdBy || 'user'
  }
  
  emit('save', updatedPrompt)
  showEditDialog.value = false
  ElMessage.success(editingPrompt.value ? 'Prompt更新成功' : 'Prompt创建成功')
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
}

.prompt-actions {
  display: flex;
  gap: 10px;
}

.prompt-meta {
  margin-bottom: 10px;
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

.prompt-content-preview {
  font-size: 13px;
  color: #555;
  background-color: #f8f9fa;
  padding: 12px;
  border-radius: 6px;
  white-space: pre-wrap;
  line-height: 1.5;
  border: 1px solid #eee;
}
</style>