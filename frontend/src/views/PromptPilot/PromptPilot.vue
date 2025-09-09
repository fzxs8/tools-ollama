<template>
  <div class="prompt-pilot">
    <el-row :gutter="20" style="height: 100%;">
      <!-- 左侧服务和模型选择区域 -->
      <el-col :span="6" style="height: 100%;">
        <ModelSelector 
          :models="availableModels"
          :servers="availableServers"
          :selected-models="selectedModels"
          :selected-server-id="selectedServerId"
          :is-loading-models="isLoadingModels"
          @update:selected-models="updateSelectedModels"
          @update:selected-server-id="updateSelectedServerId"
          @load-models="loadModelsForServer"
        />
      </el-col>

      <!-- 中间主区域 -->
      <el-col :span="18" style="height: 100%; display: flex; flex-direction: column;">
        <div class="main-content">
          <!-- 想法输入区 -->
          <div class="idea-input-section">
            <el-input
              v-model="userIdea"
              type="textarea"
              placeholder="请输入您的想法或需求描述"
              :rows="4"
              style="width: 100%;"
            />
            <div style="margin-top: 15px; text-align: right;">
              <el-button 
                type="primary" 
                @click="generatePrompt" 
                :loading="isGenerating"
                :disabled="selectedModels.length === 0 || !userIdea.trim() || !selectedServerId">
                生成Prompt
              </el-button>
            </div>
          </div>

          <!-- Prompt展示区 -->
          <div class="prompt-display-section" style="flex: 1; display: flex; flex-direction: column;">
            <div class="section-title">生成的Prompt</div>
            <div class="prompt-content-container" style="flex: 1; overflow-y: auto;">
              <el-tabs v-model="activePromptTab">
                <el-tab-pane 
                  v-for="model in selectedModels" 
                  :key="model" 
                  :label="model" 
                  :name="model">
                  <div class="prompt-content">
                    <div v-if="isGenerating" class="generating-indicator">
                      <span>正在生成</span>
                      <div class="dot"></div>
                      <div class="dot"></div>
                      <div class="dot"></div>
                    </div>
                    <div v-else>
                      <div v-html="renderedPrompt[model]"></div>
                    </div>
                  </div>
                </el-tab-pane>
              </el-tabs>
            </div>
            
            <div style="margin-top: 15px; text-align: right;">
              <el-button 
                @click="optimizePrompt" 
                :disabled="isGenerating || selectedModels.length === 0 || !selectedServerId">
                优化
              </el-button>
              <el-button 
                type="success" 
                @click="savePrompt" 
                :disabled="isGenerating || selectedModels.length === 0 || !selectedServerId"
                style="margin-left: 10px;">
                保存
              </el-button>
              <el-button 
                @click="showSavedPrompts = true"
                style="margin-left: 10px;">
                已保存
              </el-button>
            </div>
          </div>
        </div>
      </el-col>

      <!-- 抽屉区域 -->
      <el-drawer
        v-model="showOptimizeDrawer"
        title="优化Prompt"
        direction="rtl"
        size="400px"
      >
        <div style="padding: 20px;">
          <el-input
            v-model="optimizationFeedback"
            type="textarea"
            placeholder="请提供优化建议或反馈"
            :rows="6"
          />
          <div style="margin-top: 20px; text-align: right;">
            <el-button @click="showOptimizeDrawer = false">取消</el-button>
            <el-button type="primary" @click="performOptimization" style="margin-left: 10px;">优化</el-button>
          </div>
        </div>
      </el-drawer>

      <!-- 保存抽屉 -->
      <el-drawer
        v-model="showSaveDrawer"
        title="保存Prompt"
        direction="rtl"
        size="400px"
      >
        <div style="padding: 20px;">
          <el-form :model="promptToSave" label-width="80px">
            <el-form-item label="名称">
              <el-input v-model="promptToSave.name" placeholder="请输入Prompt名称" />
            </el-form-item>
            <el-form-item label="描述">
              <el-input 
                v-model="promptToSave.description" 
                type="textarea" 
                placeholder="请输入描述（可选）" 
                :rows="3" />
            </el-form-item>
            <el-form-item label="标签">
              <el-select
                v-model="promptToSave.tags"
                multiple
                filterable
                allow-create
                default-first-option
                placeholder="请输入标签，可创建新标签"
                style="width: 100%">
                <el-option
                  v-for="tag in allPromptTags"
                  :key="tag"
                  :label="tag"
                  :value="tag">
                </el-option>
              </el-select>
            </el-form-item>
          </el-form>
          <div style="margin-top: 20px; text-align: right;">
            <el-button @click="showSaveDrawer = false">取消</el-button>
            <el-button type="primary" @click="performSave" style="margin-left: 10px;">保存</el-button>
          </div>
        </div>
      </el-drawer>

      <!-- 已保存的Prompt抽屉 -->
      <SavedPromptsDrawer 
        v-model:visible="showSavedPrompts"
        :prompts="savedPrompts"
        @preview="previewPrompt"
        @edit="handleEditPrompt"
        @delete="deletePrompt"
        @save="handleSaveEditedPrompt"
      />
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import ModelSelector from './components/ModelSelector.vue'
import SavedPromptsDrawer from './components/SavedPromptsDrawer.vue'
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

interface Prompt {
  id: string
  name: string
  content: string
  description: string
  createdAt: number
  updatedAt: number
  models: string[]
  serverId: string
  version: number
  tags: string[]
  createdBy: string
}

// 响应式数据
const availableServers = ref<Server[]>([
  { 
    id: 'local', 
    name: '本地Ollama服务', 
    baseUrl: 'http://localhost:11434', 
    apiKey: '', 
    isActive: true, 
    testStatus: 'success', 
    type: 'local' 
  },
  { 
    id: 'remote1', 
    name: '远程开发环境', 
    baseUrl: 'http://dev.example.com:11434', 
    apiKey: 'sk-xxx', 
    isActive: false, 
    testStatus: 'unknown', 
    type: 'remote' 
  }
])

const availableModels = ref<Model[]>([])
const selectedServerId = ref('')
const isLoadingModels = ref(false)
const selectedModels = ref<string[]>([])
const userIdea = ref('')
const isGenerating = ref(false)
const activePromptTab = ref('')
const showOptimizeDrawer = ref(false)
const showSaveDrawer = ref(false)
const optimizationFeedback = ref('')
const showSavedPrompts = ref(false)

const renderedPrompt = ref<Record<string, string>>({
  // 初始化为空对象
})

const promptToSave = ref({
  name: '',
  description: '',
  tags: [] as string[]
})

const savedPrompts = ref<Prompt[]>([
  {
    id: '1',
    name: '文案创作助手',
    content: '你是一个专业的文案创作助手，请帮我创作一段关于科技产品的宣传文案，要求简洁有力，能够吸引年轻人的注意。',
    description: '用于生成科技产品宣传文案',
    createdAt: Date.now() - 86400000, // 1天前
    updatedAt: Date.now() - 86400000,
    models: ['llama2'],
    serverId: 'local',
    version: 1,
    tags: ['文案', '营销'],
    createdBy: 'user'
  },
  {
    id: '2',
    name: '代码解释器',
    content: '你是一个资深的Python开发者，请解释以下代码的作用和实现原理，并指出可能的优化点。',
    description: '用于解释和分析代码',
    createdAt: Date.now() - 172800000, // 2天前
    updatedAt: Date.now() - 172800000,
    models: ['mistral'],
    serverId: 'local',
    version: 1,
    tags: ['代码', '技术'],
    createdBy: 'user'
  }
])

// 计算所有Prompt标签的集合
const allPromptTags = computed(() => {
  const tagsSet = new Set<string>()
  savedPrompts.value.forEach(prompt => {
    prompt.tags.forEach(tag => tagsSet.add(tag))
  })
  return Array.from(tagsSet)
})

// 方法定义
const updateSelectedModels = (models: string[]) => {
  selectedModels.value = models
  if (models.length > 0 && !models.includes(activePromptTab.value)) {
    activePromptTab.value = models[0]
  } else if (models.length === 0) {
    activePromptTab.value = ''
  }
}

const updateSelectedServerId = (serverId: string) => {
  selectedServerId.value = serverId
}

const loadModelsForServer = async (serverId: string) => {
  // 模拟加载模型
  isLoadingModels.value = true
  
  try {
    // 模拟API调用延迟
    await new Promise(resolve => setTimeout(resolve, 800))
    
    // 根据选中的服务返回不同的模型列表
    if (serverId === 'local') {
      availableModels.value = [
        { name: 'llama2', size: 3790000000, modified_at: '2023-10-10T10:00:00Z' },
        { name: 'mistral', size: 4620000000, modified_at: '2023-11-15T14:30:00Z' },
        { name: 'neural-chat', size: 5200000000, modified_at: '2023-12-01T09:15:00Z' },
        { name: 'codellama', size: 3790000000, modified_at: '2023-09-20T16:45:00Z' },
        { name: 'llama2-uncensored', size: 3790000000, modified_at: '2023-10-12T11:20:00Z' }
      ]
    } else {
      availableModels.value = [
        { name: 'llama2:13b', size: 7600000000, modified_at: '2023-10-10T10:00:00Z' },
        { name: 'mistral:7b', size: 4620000000, modified_at: '2023-11-15T14:30:00Z' },
        { name: 'mixtral:8x7b', size: 18500000000, modified_at: '2023-12-01T09:15:00Z' }
      ]
    }
    
    // 清空已选择的模型
    selectedModels.value = []
    activePromptTab.value = ''
  } catch (error) {
    ElMessage.error('加载模型列表失败')
  } finally {
    isLoadingModels.value = false
  }
}

const generatePrompt = async () => {
  if (!userIdea.value.trim() || selectedModels.value.length === 0 || !selectedServerId.value) {
    ElMessage.warning('请输入想法，选择服务和至少一个模型')
    return
  }

  isGenerating.value = true
  
  // 模拟生成过程
  for (const model of selectedModels.value) {
    renderedPrompt.value[model] = ''
    
    // 模拟流式生成效果
    const samplePrompt = `你是一个${model}模型，请根据以下要求生成内容：\n\n${userIdea.value}\n\n请以专业且易懂的方式回应，确保内容准确且有深度。`
    for (let i = 0; i < samplePrompt.length; i++) {
      renderedPrompt.value[model] += samplePrompt.charAt(i)
      await new Promise(resolve => setTimeout(resolve, 10)) // 模拟打字效果
    }
  }
  
  isGenerating.value = false
  ElMessage.success('Prompt生成完成')
}

const optimizePrompt = () => {
  if (selectedModels.value.length === 0 || !selectedServerId.value) {
    ElMessage.warning('请至少选择一个模型并选择服务')
    return
  }
  
  showOptimizeDrawer.value = true
  optimizationFeedback.value = ''
}

const performOptimization = async () => {
  if (!optimizationFeedback.value.trim()) {
    ElMessage.warning('请输入优化建议')
    return
  }

  isGenerating.value = true
  showOptimizeDrawer.value = false
  
  // 模拟优化过程
  for (const model of selectedModels.value) {
    const currentContent = renderedPrompt.value[model]
    renderedPrompt.value[model] = ''
    
    // 模拟流式优化效果
    const optimizedPrompt = `${currentContent}\n\n优化建议：${optimizationFeedback.value}\n\n请根据以上建议优化Prompt内容。`
    for (let i = 0; i < optimizedPrompt.length; i++) {
      renderedPrompt.value[model] += optimizedPrompt.charAt(i)
      await new Promise(resolve => setTimeout(resolve, 10)) // 模拟打字效果
    }
  }
  
  isGenerating.value = false
  ElMessage.success('Prompt优化完成')
  optimizationFeedback.value = ''
}

const savePrompt = () => {
  if (selectedModels.value.length === 0 || !selectedServerId.value) {
    ElMessage.warning('请至少选择一个模型并选择服务')
    return
  }
  
  promptToSave.value.name = userIdea.value.substring(0, 20) + (userIdea.value.length > 20 ? '...' : '')
  promptToSave.value.tags = []
  showSaveDrawer.value = true
}

const performSave = () => {
  if (!promptToSave.value.name.trim()) {
    ElMessage.warning('请输入Prompt名称')
    return
  }

  const newPrompt: Prompt = {
    id: Date.now().toString(),
    name: promptToSave.value.name,
    content: renderedPrompt.value[activePromptTab.value] || '',
    description: promptToSave.value.description,
    createdAt: Date.now(),
    updatedAt: Date.now(),
    models: [...selectedModels.value],
    serverId: selectedServerId.value,
    version: 1,
    tags: [...promptToSave.value.tags],
    createdBy: 'user'
  }

  savedPrompts.value.push(newPrompt)
  showSaveDrawer.value = false
  ElMessage.success('Prompt保存成功')
  
  // 重置表单
  promptToSave.value.name = ''
  promptToSave.value.description = ''
  promptToSave.value.tags = []
}

const previewPrompt = (prompt: Prompt) => {
  // 在标签页中显示选中的Prompt
  for (const model of prompt.models) {
    if (!renderedPrompt.value[model]) {
      renderedPrompt.value[model] = prompt.content
    }
  }
  
  // 设置选中模型和服务
  selectedModels.value = prompt.models
  selectedServerId.value = prompt.serverId
  activePromptTab.value = prompt.models[0]
  
  // 加载对应服务的模型
  loadModelsForServer(prompt.serverId)
  
  ElMessage.success(`已加载Prompt: ${prompt.name}`)
}

const handleEditPrompt = (prompt: Prompt) => {
  // 这里可以实现编辑功能
  ElMessage.info('编辑功能已通过抽屉实现')
}

const handleSaveEditedPrompt = (prompt: Prompt) => {
  // 更新或添加Prompt
  const index = savedPrompts.value.findIndex(p => p.id === prompt.id)
  if (index !== -1) {
    // 更新现有Prompt
    savedPrompts.value[index] = prompt
  } else {
    // 添加新Prompt
    savedPrompts.value.push(prompt)
  }
  
  ElMessage.success('Prompt保存成功')
}

const deletePrompt = (id: string) => {
  savedPrompts.value = savedPrompts.value.filter(prompt => prompt.id !== id)
  ElMessage.success('Prompt删除成功')
}

const renderMarkdown = (content: string) => {
  return md.render(content)
}

// 渲染Prompt内容
const renderedPromptContent = (model: string) => {
  return renderMarkdown(renderedPrompt.value[model] || '')
}

onMounted(() => {
  // 初始化时显示一个示例Prompt
  renderedPrompt.value['llama2'] = '你是一个AI助手，请根据用户的需求提供专业的帮助和建议。'
})
</script>

<style scoped>
.prompt-pilot {
  height: 100%;
  padding: 20px;
  box-sizing: border-box;
}

.main-content {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.idea-input-section {
  margin-bottom: 20px;
}

.section-title {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 15px;
  color: #333;
}

.prompt-content-container {
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  padding: 15px;
  background-color: #f9f9f9;
}

.prompt-content {
  min-height: 200px;
  white-space: pre-wrap;
  line-height: 1.6;
}

.generating-indicator {
  display: flex;
  align-items: center;
  font-size: 16px;
  color: #4a5568;
}

.generating-indicator .dot {
  width: 8px;
  height: 8px;
  background-color: #a0aec0;
  border-radius: 50%;
  margin: 0 4px;
  animation: bounce 1.5s infinite;
}

.generating-indicator .dot:nth-child(2) {
  animation-delay: 0.2s;
}

.generating-indicator .dot:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-5px); }
}

.sidebar {
  height: 100%;
  border-left: 1px solid #e0e0e0;
  padding: 20px;
  box-sizing: border-box;
}

.saved-prompts-trigger {
  cursor: pointer;
  display: flex;
  align-items: center;
  padding: 10px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.saved-prompts-trigger:hover {
  background-color: #f0f0f0;
}

.saved-prompts-trigger .el-icon {
  margin-right: 10px;
  font-size: 20px;
}
</style>