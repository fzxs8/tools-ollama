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
            :disabled="isGenerating || selectedModels.length === 0 || !userIdea.trim() || !selectedServerId"
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
              class="prompt-tab-pane"
            >
              <div class="prompt-content">
                <div v-if="generatingModels[model]" class="generating-indicator">
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
                      :disabled="generatingModels[model] || !renderedPrompt[model]"
                    >
                      复制
                    </el-button>
                    <el-button
                      type="primary"
                      link
                      @click="regenerateSinglePrompt(model)"
                      size="small"
                      :disabled="isGenerating"
                    >
                      重新生成
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
            @click="openOptimizeDrawer"
            :disabled="isGenerating || selectedModels.length === 0 || !activePromptTab || !renderedPrompt[activePromptTab]"
          >
            优化
          </el-button>
          <el-button
            type="success"
            @click="openSaveDrawer"
            :disabled="isGenerating || selectedModels.length === 0 || !activePromptTab || !renderedPrompt[activePromptTab]"
          >
            保存
          </el-button>
          <el-button @click="showSavedPrompts = true" :disabled="isGenerating">
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

    <SavedPromptsDrawer
      v-model:visible="showSavedPrompts"
      :prompts="savedPrompts"
      @delete="handleDeletePrompt"
      @save="handleSavePrompt"
      @preview="handlePreviewPrompt"
    />

  </div>
</template>

<script setup lang="ts">
import {onMounted, ref, watch} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import ModelSelector from './components/ModelSelector.vue'
import SavedPromptsDrawer from './components/SavedPromptsDrawer.vue'
import {EventsOn} from "../../../wailsjs/runtime";
import {main} from "../../../wailsjs/go/models";
import {DeletePrompt, GeneratePromptStream, ListPrompts, SavePrompt} from "../../../wailsjs/go/main/App";

// 类型定义
type Prompt = main.Prompt;

// 响应式数据
const selectedServerId = ref('')
const selectedModels = ref<string[]>([])
const userIdea = ref('')
const isGenerating = ref(false)
const generatingModels = ref<Record<string, boolean>>({})
const activePromptTab = ref('')
const renderedPrompt = ref<Record<string, string>>({})

const showOptimizeDrawer = ref(false)
const showSaveDrawer = ref(false)
const showSavedPrompts = ref(false)

const savedPrompts = ref<Prompt[]>([])

// --- 生命周期钩子 ---
onMounted(() => {
  fetchPrompts();

  // 监听后端的流式数据事件
  EventsOn('prompt_pilot_stream', (data: { model: string; chunk: string }) => {
    if (data && data.model && data.chunk) {
      renderedPrompt.value[data.model] += data.chunk;
    }
  });

  // 监听后端的错误事件
  EventsOn('prompt_pilot_stream_error', (data: { model: string; error: string }) => {
    if (data && data.model) {
      ElMessage.error(`模型 ${data.model} 生成失败: ${data.error}`);
      generatingModels.value[data.model] = false;
      checkAllModelsFinished();
    }
  });

  // 监听后端的完成事件
  EventsOn('prompt_pilot_stream_done', (data: { model: string }) => {
    if (data && data.model) {
      generatingModels.value[data.model] = false;
      checkAllModelsFinished();
    }
  });
});

// --- 逻辑方法 ---

// 检查是否所有模型都已生成完毕
const checkAllModelsFinished = () => {
  const allDone = Object.values(generatingModels.value).every(status => !status);
  if (allDone) {
    isGenerating.value = false;
  }
};

// 监视 selectedModels 的变化，以更新 activePromptTab
watch(selectedModels, (newModels) => {
  if (newModels.length > 0 && !newModels.includes(activePromptTab.value)) {
    activePromptTab.value = newModels[0]
  } else if (newModels.length === 0) {
    activePromptTab.value = ''
  }
});

// 触发后端开始生成
const performStreamGeneration = async (model: string) => {
  generatingModels.value[model] = true;
  renderedPrompt.value[model] = '';
  try {
    await GeneratePromptStream(userIdea.value, model, selectedServerId.value);
  } catch (error: any) {
    ElMessage.error(`调用模型 ${model} 失败: ${error.message || error}`);
    generatingModels.value[model] = false;
    checkAllModelsFinished();
  }
}

// 点击主“生成”按钮的函数
const generatePrompt = async () => {
  if (isGenerating.value) return;
  if (!userIdea.value.trim() || selectedModels.value.length === 0 || !selectedServerId.value) {
    ElMessage.warning('请输入想法，选择服务和至少一个模型');
    return;
  }

  isGenerating.value = true;
  // 重置状态
  generatingModels.value = {};
  selectedModels.value.forEach(model => {
    generatingModels.value[model] = true;
  });

  // 并行触发所有模型的生成
  selectedModels.value.forEach(model => {
    performStreamGeneration(model);
  });
}

// 点击“重新生成”按钮的函数
const regenerateSinglePrompt = async (model: string) => {
  if (isGenerating.value) {
    ElMessage.warning('正在等待所有模型生成完成，请稍后再试。');
    return;
  }
  isGenerating.value = true;
  generatingModels.value = { [model]: true };
  await performStreamGeneration(model);
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

const openOptimizeDrawer = () => {
  showOptimizeDrawer.value = true
}

const openSaveDrawer = () => {
  // 这里可以预填一些保存表单的数据
  showSaveDrawer.value = true
}

// --- 已保存的Prompt相关逻辑 ---

const fetchPrompts = async () => {
  try {
    savedPrompts.value = await ListPrompts();
  } catch (error: any) {
    ElMessage.error('获取已保存的Prompt列表失败: ' + error.message);
  }
};

const handleSavePrompt = async (promptToSave: Prompt) => {
  try {
    await SavePrompt(promptToSave);
    ElMessage.success('Prompt保存成功');
    await fetchPrompts(); // 刷新列表
  } catch (error: any) {
    ElMessage.error('保存Prompt失败: ' + error.message);
  }
};

const handleDeletePrompt = async (id: string) => {
  try {
    await ElMessageBox.confirm('确定要删除这个Prompt吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    });
    await DeletePrompt(id);
    ElMessage.success('Prompt删除成功');
    await fetchPrompts(); // 刷新列表
  } catch (error) {
    // 如果用户点击取消，会进入catch，但我们不需要显示错误信息
    if (error !== 'cancel') {
      ElMessage.error('删除Prompt失败: ' + error);
    }
  }
};

const handlePreviewPrompt = (prompt: Prompt) => {
  // 将预览的Prompt内容填充到当前激活的tab中
  if (activePromptTab.value) {
    renderedPrompt.value[activePromptTab.value] = prompt.content;
    ElMessage.info(`已将Prompt“${prompt.name}”的内容加载到当前视图`);
  } else {
    ElMessage.warning('请先选择一个模型以加载Prompt内容');
  }
  showSavedPrompts.value = false; // 关闭抽屉
};

</script>

<style>
/* Global style for the dropdown */
.left-aligned-dropdown .el-select-dropdown__item {
  justify-content: flex-start !important;
}
</style>

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
  min-height: 0;
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
  min-height: 0;
}
.prompt-tabs > :deep(.el-tabs__content) {
  flex: 1;
  padding: 15px;
  min-height: 0;
}
.prompt-tab-pane {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.prompt-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  line-height: 1.6;
  min-height: 0;
}
.prompt-raw-content-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  position: relative;
  min-height: 0;
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
  flex: 1;
  overflow-y: auto;
  text-align: left;
}
.prompt-actions {
  position: absolute;
  top: 5px;
  right: 5px;
  display: flex;
  gap: 8px;
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