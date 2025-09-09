<template>
  <el-drawer
      v-model="visible"
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
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { ElMessage } from 'element-plus'

interface SystemPrompt {
  id: string
  title: string
  prompt: string
  createdAt: number
}

interface SystemPromptForm {
  title: string
  prompt: string
}

const props = defineProps<{
  visible: boolean
  activeSystemPrompt: SystemPrompt | null
  systemPromptList: SystemPrompt[]
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'update:activeSystemPrompt', value: SystemPrompt | null): void
  (e: 'update:systemPromptList', value: SystemPrompt[]): void
  (e: 'save-system-prompt', prompt: SystemPrompt): void
  (e: 'update-system-prompt', prompt: SystemPrompt): void
  (e: 'delete-system-prompt', id: string): void
}>()

// 系统提示词表单
const systemPromptForm = reactive<SystemPromptForm>({
  title: '',
  prompt: ''
})

const visible = defineModel<boolean>('visible', { default: false })

// 重置系统提示词表单
const resetSystemPromptForm = () => {
  systemPromptForm.title = ''
  systemPromptForm.prompt = ''
}

// 设置激活的系统提示词
const setActiveSystemPrompt = (prompt: SystemPrompt) => {
  emit('update:activeSystemPrompt', prompt)
}

// 删除系统提示词
const deleteSystemPrompt = (prompt: SystemPrompt) => {
  emit('delete-system-prompt', prompt.id)
}

// 编辑系统提示词
const editSystemPrompt = (prompt: SystemPrompt) => {
  systemPromptForm.title = prompt.title
  systemPromptForm.prompt = prompt.prompt
}

// 保存系统提示词
const saveSystemPrompt = () => {
  if (!systemPromptForm.title.trim()) {
    ElMessage.warning('请输入提示词标题')
    return
  }

  if (!systemPromptForm.prompt.trim()) {
    ElMessage.warning('请输入提示词内容')
    return
  }

  const newPrompt: SystemPrompt = {
    id: Date.now().toString(),
    title: systemPromptForm.title,
    prompt: systemPromptForm.prompt,
    createdAt: Date.now()
  }

  emit('save-system-prompt', newPrompt)
  resetSystemPromptForm()
}

// 更新系统提示词
const updateSystemPrompt = () => {
  if (!systemPromptForm.title.trim()) {
    ElMessage.warning('请输入提示词标题')
    return
  }

  if (!systemPromptForm.prompt.trim()) {
    ElMessage.warning('请输入提示词内容')
    return
  }

  // 查找并更新提示词
  const index = props.systemPromptList.findIndex(p =>
      p.title === systemPromptForm.title && p.prompt !== systemPromptForm.prompt)

  if (index !== -1) {
    const updatedPrompt = {
      ...props.systemPromptList[index],
      prompt: systemPromptForm.prompt,
      title: systemPromptForm.title
    }
    
    emit('update-system-prompt', updatedPrompt)
    resetSystemPromptForm()
  }
}
</script>

<style scoped>
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