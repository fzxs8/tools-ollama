<!-- FORCE UPDATE: 2025-09-10 20:10 PM -->
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
            @click="openSaveDrawer(null)"
            :disabled="isGenerating || selectedModels.length === 0 || !activePromptTab || !renderedPrompt[activePromptTab]"
          >
            保存
          </el-button>
          <el-button @click="showSavedPrompts = true" :disabled="isGenerating">
            我的提示词
          </el-button>
        </div>
      </div>
    </div>

    <!-- 抽屉区域 -->
    <el-drawer v-model="showOptimizeDrawer" title="优化Prompt" direction="rtl" size="40%">
      <!-- ...抽屉内容... -->
    </el-drawer>

    <el-drawer v-model="showSaveDrawer" :title="promptToSave.id ? `编辑 “${promptToSave.name}”` : '保存Prompt'" direction="rtl" size="40%">
      <div class="save-drawer-content">
        <el-form :model="promptToSave" label-position="top" ref="saveFormRef">
          <el-form-item label="标题" prop="name" :rules="[{ required: true, message: '标题不能为空', trigger: 'blur' }]">
            <el-input v-model="promptToSave.name" placeholder="请输入Prompt标题"></el-input>
          </el-form-item>
          <el-form-item label="描述">
            <el-input v-model="promptToSave.description" type="textarea" :rows="3" placeholder="请输入详细描述、使用场景等"></el-input>
          </el-form-item>
          <el-form-item label="标签">
            <el-select v-model="promptToSave.tags" multiple filterable allow-create default-first-option placeholder="选择或创建标签" style="width: 100%;">
              <el-option v-for="tag in existingTags" :key="tag" :label="tag" :value="tag"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="关联模型" class="left-aligned-item">
            <div>
              <el-tag v-for="model in promptToSave.models" :key="model" type="info" style="margin-right: 5px;">{{ model }}</el-tag>
            </div>
          </el-form-item>
          <el-form-item label="Prompt内容 (只读预览)" class="left-aligned-item content-preview-item">
            <el-input v-model="promptToSave.content" type="textarea" :rows="8" readonly></el-input>
            <el-button class="content-copy-btn" type="primary" link @click="copySaveDrawerContent">复制</el-button>
          </el-form-item>
        </el-form>
        <div class="save-drawer-footer">
          <el-button @click="showSaveDrawer = false">取消</el-button>
          <el-button type="primary" @click="executeSaveFromDrawer" :loading="isSaving">确认保存</el-button>
        </div>
      </div>
    </el-drawer>

    <CommonPromptDrawer
      v-model:visible="showSavedPrompts"
      :prompts="savedPrompts"
      mode="manage"
      @delete="handleDeletePrompt"
      @edit="handleEditPrompt"
      @preview="handlePreviewPrompt"
    />

  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, ref, watch} from 'vue'
import {ElMessage, ElMessageBox, FormInstance} from 'element-plus'
import ModelSelector from './components/ModelSelector.vue'
import CommonPromptDrawer from './components/CommonPromptDrawer.vue' // Changed from SavedPromptsDrawer
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
const isSaving = ref(false)

const savedPrompts = ref<Prompt[]>([])
const saveFormRef = ref<FormInstance>()
const promptToSave = ref<Partial<Prompt>>({
  id: undefined,
  name: '',
  description: '',
  tags: [],
  models: [],
  content: ''
});

// 计算属性
const existingTags = computed(() => {
  const tags = new Set<string>();
  savedPrompts.value.forEach(p => {
    p.tags?.forEach(t => tags.add(t));
  });
  return Array.from(tags);
});

// --- 生命周期钩子 ---
onMounted(() => {
  fetchPrompts();

  EventsOn('prompt_pilot_stream', (data: { model: string; chunk: string }) => {
    if (data && data.model && data.chunk) {
      renderedPrompt.value[data.model] = (renderedPrompt.value[data.model] || '') + data.chunk;
    }
  });

  EventsOn('prompt_pilot_stream_error', (data: { model: string; error: string }) => {
    if (data && data.model) {
      ElMessage.error(`模型 ${data.model} 生成失败: ${data.error}`);
      generatingModels.value[data.model] = false;
      checkAllModelsFinished();
    }
  });

  EventsOn('prompt_pilot_stream_done', (data: { model: string }) => {
    if (data && data.model) {
      generatingModels.value[data.model] = false;
      checkAllModelsFinished();
    }
  });
});

// --- 逻辑方法 ---

const checkAllModelsFinished = () => {
  const allDone = Object.values(generatingModels.value).every(status => !status);
  if (allDone) {
    isGenerating.value = false;
  }
};

watch(selectedModels, (newModels) => {
  if (newModels.length > 0 && !newModels.includes(activePromptTab.value)) {
    activePromptTab.value = newModels[0]
  } else if (newModels.length === 0) {
    activePromptTab.value = ''
  }
});

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

const generatePrompt = async () => {
  if (isGenerating.value) return;
  if (!userIdea.value.trim() || selectedModels.value.length === 0 || !selectedServerId.value) {
    ElMessage.warning('请输入想法，选择服务和至少一个模型');
    return;
  }

  isGenerating.value = true;
  generatingModels.value = {};
  selectedModels.value.forEach(model => {
    generatingModels.value[model] = true;
  });

  selectedModels.value.forEach(model => {
    performStreamGeneration(model);
  });
}

const regenerateSinglePrompt = async (model: string) => {
  if (isGenerating.value) {
    ElMessage.warning('正在等待所有模型生成完成，请稍后再试。');
    return;
  }
  isGenerating.value = true;
  generatingModels.value = {[model]: true};
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

const copySaveDrawerContent = () => {
  if (promptToSave.value.content) {
    navigator.clipboard.writeText(promptToSave.value.content).then(() => {
      ElMessage.success('内容已复制到剪贴板')
    }).catch(err => {
      ElMessage.error('复制失败: ' + err)
    })
  }
}

const openOptimizeDrawer = () => {
  showOptimizeDrawer.value = true
}

const openSaveDrawer = (existingPrompt: Prompt | null) => {
  if (existingPrompt) {
    // 编辑模式: 深拷贝以避免意外修改原始列表中的数据
    promptToSave.value = JSON.parse(JSON.stringify(existingPrompt));
  } else {
    // 新建模式
    promptToSave.value = {
      id: undefined,
      name: userIdea.value.substring(0, 20) || '',
      description: '',
      tags: [],
      models: [...selectedModels.value],
      content: renderedPrompt.value[activePromptTab.value] || ''
    };
  }
  showSaveDrawer.value = true;
}

const fetchPrompts = async () => {
  try {
    savedPrompts.value = await ListPrompts();
  } catch (error: any) {
    ElMessage.error('获取已保存的Prompt列表失败: ' + error.message);
  }
};

const executeSaveFromDrawer = async () => {
  if (!saveFormRef.value) return;
  await saveFormRef.value.validate(async (valid) => {
    if (valid) {
      isSaving.value = true;
      try {
        await SavePrompt(promptToSave.value as Prompt);
        ElMessage.success(promptToSave.value.id ? 'Prompt更新成功' : 'Prompt保存成功');
        await fetchPrompts();
        showSaveDrawer.value = false;
      } catch (error: any) {
        ElMessage.error('保存Prompt失败: ' + error.message);
      } finally {
        isSaving.value = false;
      }
    }
  });
};

const handleEditPrompt = (prompt: Prompt) => {
  openSaveDrawer(prompt);
}

const handleDeletePrompt = async (prompt: Prompt) => {
  try {
    await ElMessageBox.confirm(
        `您确定要删除 “${prompt.name}” 吗？`,
        '警告',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        }
    );
    await DeletePrompt(prompt.id);
    ElMessage.success('Prompt删除成功');
    await fetchPrompts();
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除Prompt失败: ' + String(error));
    }
  }
};

const handlePreviewPrompt = (prompt: Prompt) => {
  if (activePromptTab.value) {
    renderedPrompt.value[activePromptTab.value] = prompt.content;
    ElMessage.info(`已将Prompt“${prompt.name}”的内容加载到当前视图`);
  } else {
    ElMessage.warning('请先选择一个模型以加载Prompt内容');
  }
  showSavedPrompts.value = false;
};

</script>

<style>
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
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
  min-height: 0;
}

.idea-input-section {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  flex-shrink: 0;
}

.controls-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 15px;
}

.prompt-display-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
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
  display: flex;
  flex-direction: column;
}

.prompt-tab-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
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
  min-height: 0;
}

.prompt-actions {
  position: absolute;
  top: 5px;
  right: 5px;
  display: flex;
  gap: 8px;
}

.empty-state, .generating-indicator {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #909399;
  font-size: 14px;
}

.action-buttons-footer, .save-drawer-footer {
  text-align: right;
  margin-top: 15px;
  flex-shrink: 0;
}

.save-drawer-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.save-drawer-content .el-form {
  flex: 1;
  overflow-y: auto;
  padding: 0 10px;
}

.left-aligned-item :deep(.el-form-item__content) {
  justify-content: flex-start;
}

.content-preview-item {
  position: relative;
}

.content-copy-btn {
  position: absolute;
  top: 0px;
  right: 5px;
}

</style>