<template>
  <div class="prompt-pilot">
    <div class="main-content">
      <!-- 想法输入区 -->
      <div class="idea-input-section">
        <el-input
          v-model="userIdea"
          type="textarea"
          placeholder="请输入您的想法或需求描述，例如：我想生成一个专业的周报，内容包括本周完成的工作、遇到的问题和下周计划。"
          :rows="5"
          style="width: 100%;"
        />
        <div class="controls-container">
          <ModelSelector
            v-model:selectedServer="selectedServerId"
            v-model:selectedModels="selectedModels"
          />
          <el-button
            type="primary"
            @click="generatePrompt"
            :loading="isGenerating"
            :disabled="selectedModels.length === 0 || !userIdea.trim() || !selectedServerId"
          >
            生成
          </el-button>
        </div>
      </div>

      <!-- Prompt展示区 -->
      <div class="prompt-display-section">
        <div class="section-title">提示词内容</div>
        <div class="prompt-content-container">
          <el-tabs v-model="activePromptTab" class="prompt-tabs">
            <el-tab-pane
              v-for="model in selectedModels"
              :key="model"
              :label="model"
              :name="model"
            >
              <div class="prompt-content">
                <div v-if="isGenerating && activePromptTab === model" class="generating-indicator">
                  <span>正在生成</span>
                  <div class="dot"></div>
                  <div class="dot"></div>
                  <div class="dot"></div>
                </div>
                <div v-else class="prompt-raw-content-wrapper">
                  <pre class="prompt-raw-content">{{ renderedPrompt[model] }}</pre>
                  <div class="prompt-actions">
                    <el-button
                      type="primary"
                      link
                      @click="copyPrompt(model)"
                      size="small"
                    >
                      复制
                    </el-button>
                  </div>
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
          <div v-if="selectedModels.length === 0 && !isGenerating" class="empty-state">
            请先选择服务和模型，然后输入您的想法，点击“生成”开始。
          </div>
        </div>

        <div class="action-buttons-footer">
          <el-button
            @click="optimizePrompt"
            :disabled="isGenerating || selectedModels.length === 0 || !activePromptTab || !renderedPrompt[activePromptTab]"
          >
            优化
          </el-button>
          <el-button
            type="success"
            @click="savePrompt"
            :disabled="isGenerating || selectedModels.length === 0 || !activePromptTab || !renderedPrompt[activePromptTab]"
          >
            保存
          </el-button>
          <el-button @click="showSavedPrompts = true">
            已保存
          </el-button>
        </div>
      </div>
    </div>

    <!-- 抽屉区域 -->
    <el-drawer v-model="showOptimizeDrawer" title="优化Prompt" direction="rtl" size="40%">
      <!-- ...抽屉内容... -->
    </el-drawer>

    <el-drawer v-model="showSaveDrawer" title="保存Prompt" direction="rtl" size="40%">
      <!-- ...抽屉内容... -->
    </el-drawer>

    <SavedPromptsDrawer v-model:visible="showSavedPrompts" />

  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import ModelSelector from './components/ModelSelector.vue'
import SavedPromptsDrawer from './components/SavedPromptsDrawer.vue'

// 响应式数据
const selectedServerId = ref('')
const selectedModels = ref<string[]>([])
const userIdea = ref('')
const isGenerating = ref(false)
const activePromptTab = ref('')
const renderedPrompt = ref<Record<string, string>>({})

const showOptimizeDrawer = ref(false)
const showSaveDrawer = ref(false)
const showSavedPrompts = ref(false)

// 监视 selectedModels 的变化，以更新 activePromptTab
watch(selectedModels, (newModels, oldModels) => {
  if (newModels.length > 0 && !newModels.includes(activePromptTab.value)) {
    activePromptTab.value = newModels[0]
  } else if (newModels.length === 0) {
    activePromptTab.value = ''
  }
});

const generatePrompt = async () => {
  if (!userIdea.value.trim() || selectedModels.value.length === 0 || !selectedServerId.value) {
    ElMessage.warning('请输入想法，选择服务和至少一个模型')
    return
  }

  isGenerating.value = true
  renderedPrompt.value = {} // 清空之前的内容

  // 模拟并行生成
  const generationPromises = selectedModels.value.map(async (model) => {
    renderedPrompt.value[model] = '' // 初始化为空
    const samplePrompt = `作为 expert in prompt engineering, 请根据以下要求，为大语言模型（LLM）设计一个清晰、具体、结构化的 Prompt。\n\n## 原始需求\n${userIdea.value}\n\n## 设计要求\n1.  **角色（Role）**：明确定义 LLM 需要扮演的专家角色。\n2.  **任务（Task）**：清晰、分步骤地描述需要完成的核心任务。\n3.  **格式（Format）**：指定输出的格式，如 Markdown、JSON 等，并提供示例。\n4.  **约束（Constraint）**：提出明确的限制和要求，如内容风格、长度、禁止项等。\n5.  **示例（Example）**：提供一或两个输入/输出示例，帮助 LLM 理解期望。\n\n请只输出设计好的 Prompt 内容，不要包含任何额外的解释或对话。`
    
    // 模拟流式生成效果
    for (let i = 0; i < samplePrompt.length; i++) {
      renderedPrompt.value[model] += samplePrompt.charAt(i)
      await new Promise(resolve => setTimeout(resolve, 5)) // 模拟打字效果
    }
  });

  try {
    await Promise.all(generationPromises);
    ElMessage.success('Prompt生成完成')
  } catch (error) {
    ElMessage.error('生成过程中出现错误')
    console.error(error)
  } finally {
    isGenerating.value = false
  }
}

const copyPrompt = (model: string) => {
  const content = renderedPrompt.value[model]
  if (content) {
    navigator.clipboard.writeText(content).then(() => {
      ElMessage.success('Prompt已复制到剪贴板')
    }).catch(err => {
      ElMessage.error('复制失败: ' + err)
    })
  }
}

const optimizePrompt = () => {
  showOptimizeDrawer.value = true
}

const savePrompt = () => {
  showSaveDrawer.value = true
}

</script>

<style scoped>
.prompt-pilot {
  height: 100%;
  padding: 20px;
  box-sizing: border-box;
  background-color: #f4f5f7;
  display: flex;
  flex-direction: column;
}
.main-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 20px;
  flex: 1;
  min-height: 0;
}
.idea-input-section {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.05);
  flex-shrink: 0;
}
.controls-container {
  display: flex;
  align-items: center;
  justify-content: space-between; /* 实现左右对齐 */
  margin-top: 15px;
}
.prompt-display-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.05);
  min-height: 0; /* flexbox 溢出修复 */
}
.section-title {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 15px;
  color: #333;
}
.prompt-content-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  min-height: 0;
}
.prompt-tabs {
  flex: 1;
  display: flex;
  flex-direction: column;
}
.prompt-tabs > :deep(.el-tabs__content) {
  flex: 1;
  overflow-y: auto;
  padding: 15px;
}
.prompt-content {
  height: 100%;
  line-height: 1.6;
}
.prompt-raw-content-wrapper {
  position: relative;
}
pre.prompt-raw-content {
  white-space: pre-wrap;
  word-break: break-word;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  background-color: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 4px;
  padding: 15px;
  margin: 0;
  max-height: calc(100vh - 500px); /* 需要根据实际情况微调 */
  overflow-y: auto;
  text-align: left;
}
.prompt-actions {
  position: absolute;
  top: 5px;
  right: 5px;
}
.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #909399;
  font-size: 14px;
}
.action-buttons-footer {
  margin-top: 15px;
  text-align: right;
  flex-shrink: 0;
}
.generating-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100px;
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
.generating-indicator .dot:nth-child(2) { animation-delay: 0.2s; }
.generating-indicator .dot:nth-child(3) { animation-delay: 0.4s; }
@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-5px); }
}
</style>