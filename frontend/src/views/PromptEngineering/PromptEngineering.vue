<!-- FORCE UPDATE: 2025-09-10 20:10 PM -->
<template>
  <div class="prompt-engineering">
    <div class="main-content">
      <!-- 想法输入区 -->
      <div class="idea-input-section">
        <el-input
          v-model="userIdea"
          type="textarea"
          :placeholder="t('promptEngineering.typeMessage')"
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
            {{ t('promptEngineering.sendMessage') }}
          </el-button>
        </div>
      </div>

      <!-- Prompt展示区 -->
      <div class="prompt-display-section">
        <div class="section-title">{{ t('promptEngineering.systemPrompt') }}</div>
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
                  <span>{{ t('common.loading') }}</span>
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
                      {{ t('common.copy') }}
                    </el-button>
                    <el-button
                      type="primary"
                      link
                      @click="regenerateSinglePrompt(model)"
                      size="small"
                      :disabled="isGenerating"
                    >
                      {{ t('chatManager.regenerate') }}
                    </el-button>
                  </div>
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
          <div v-if="selectedModels.length === 0 && !isGenerating" class="empty-state">
            {{ t('promptEngineering.selectModel') }}
          </div>
        </div>

        <div class="action-buttons-footer">
          <el-button
            @click="openOptimizeDrawer"
            :disabled="isGenerating || selectedModels.length === 0 || !activePromptTab || !renderedPrompt[activePromptTab]"
          >
            {{ t('promptEngineering.optimizePrompt') }}
          </el-button>
          <el-button
            type="success"
            @click="openSaveDrawer(null)"
            :disabled="isGenerating || selectedModels.length === 0 || !activePromptTab || !renderedPrompt[activePromptTab]"
          >
            {{ t('common.save') }}
          </el-button>
          <el-button @click="showSavedPrompts = true" :disabled="isGenerating">
            {{ t('promptEngineering.myPrompts') }}
          </el-button>
        </div>
      </div>
    </div>

    <!-- 抽屉区域 -->
    <el-drawer v-model="showOptimizeDrawer" :title="t('promptEngineering.optimizePrompt')" direction="rtl" size="40%">
      <!-- ...抽屉内容... -->
    </el-drawer>

    <el-drawer v-model="showSaveDrawer" :title="drawerTitle" direction="rtl" size="40%">
      <div class="save-drawer-content">
        <el-form :model="promptToSave" label-position="top" ref="saveFormRef">
          <el-form-item :label="t('promptEngineering.promptTitle')" prop="name" :rules="[{ required: true, message: t('promptEngineering.titleRequired'), trigger: 'blur' }]">
            <el-input v-model="promptToSave.name" :placeholder="t('promptEngineering.enterTitle')"></el-input>
          </el-form-item>
          <el-form-item :label="t('promptEngineering.promptDescription')">
            <el-input v-model="promptToSave.description" type="textarea" :rows="3" :placeholder="t('promptEngineering.enterDescription')"></el-input>
          </el-form-item>
          <el-form-item :label="t('promptEngineering.tags')">
            <el-select v-model="promptToSave.tags" multiple filterable allow-create default-first-option :placeholder="t('promptEngineering.selectTags')" style="width: 100%;">
              <el-option v-for="tag in existingTags" :key="tag" :label="tag" :value="tag"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="t('promptEngineering.associatedModels')" class="left-aligned-item">
            <div>
              <el-tag v-for="model in promptToSave.models" :key="model" type="info" style="margin-right: 5px;">{{ model }}</el-tag>
            </div>
          </el-form-item>
          <el-form-item :label="t('promptEngineering.promptContent')" class="left-aligned-item content-preview-item">
            <el-input v-model="promptToSave.content" type="textarea" :rows="8" readonly></el-input>
            <el-button class="content-copy-btn" type="primary" link @click="copySaveDrawerContent">{{ t('common.copy') }}</el-button>
          </el-form-item>
        </el-form>
        <div class="save-drawer-footer">
          <el-button @click="showSaveDrawer = false">{{ t('common.cancel') }}</el-button>
          <el-button type="primary" @click="executeSaveFromDrawer" :loading="isSaving">{{ t('common.confirm') }}</el-button>
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
import { useI18n } from 'vue-i18n'
import ModelSelector from './components/ModelSelector.vue'
import CommonPromptDrawer from '../../components/commons/PromptListDrawer.vue'
import {EventsOn} from "../../../wailsjs/runtime";
import {DeletePrompt, GeneratePromptStream, ListPrompts, SavePrompt} from "../../../wailsjs/go/main/App";
import {types} from "../../../wailsjs/go/models";
import Prompt = types.Prompt;

const { t } = useI18n();

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

const drawerTitle = computed(() => {
  return promptToSave.value.id 
    ? `${t('common.edit')} "${promptToSave.value.name}"`
    : t('promptEngineering.savePrompt');
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
      ElMessage.error(`${t('chatManager.modelGenerationFailed')} ${data.model}: ${data.error}`);
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
    ElMessage.error(`${t('chatManager.modelCallFailed')} ${model}: ${error.message || error}`);
    generatingModels.value[model] = false;
    checkAllModelsFinished();
  }
}

const generatePrompt = async () => {
  if (isGenerating.value) return;
  if (!userIdea.value.trim() || selectedModels.value.length === 0 || !selectedServerId.value) {
    ElMessage.warning(t('promptEngineering.enterIdeaAndSelectModel'));
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
    ElMessage.warning(t('promptEngineering.waitForGeneration'));
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
      ElMessage.success(t('promptEngineering.promptCopied'))
    }).catch(err => {
      ElMessage.error(t('promptEngineering.copyFailed') + ': ' + err)
    })
  }
}

const copySaveDrawerContent = () => {
  if (promptToSave.value.content) {
    navigator.clipboard.writeText(promptToSave.value.content).then(() => {
      ElMessage.success(t('promptEngineering.contentCopied'))
    }).catch(err => {
      ElMessage.error(t('promptEngineering.copyFailed') + ': ' + err)
    })
  }
}

const openOptimizeDrawer = () => {
  showOptimizeDrawer.value = true
}

const openSaveDrawer = (existingPrompt: Prompt | null) => {
  if (existingPrompt) {
    promptToSave.value = JSON.parse(JSON.stringify(existingPrompt));
  } else {
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
    ElMessage.error(t('promptEngineering.fetchPromptsFailed') + ': ' + error.message);
  }
};

const executeSaveFromDrawer = async () => {
  if (!saveFormRef.value) return;
  await saveFormRef.value.validate(async (valid) => {
    if (valid) {
      isSaving.value = true;
      try {
        await SavePrompt(promptToSave.value as Prompt);
        ElMessage.success(promptToSave.value.id ? t('promptEngineering.promptUpdated') : t('promptEngineering.promptSaved'));
        await fetchPrompts();
        showSaveDrawer.value = false;
      } catch (error: any) {
        ElMessage.error(t('promptEngineering.savePromptFailed') + ': ' + error.message);
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
        t('promptEngineering.deleteConfirm', { name: prompt.name }),
        t('promptEngineering.warning'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning',
        }
    );
    await DeletePrompt(prompt.id);
    ElMessage.success(t('promptEngineering.promptDeleted'));
    await fetchPrompts();
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('promptEngineering.deletePromptFailed') + ': ' + String(error));
    }
  }
};

const handlePreviewPrompt = (prompt: Prompt) => {
  if (activePromptTab.value) {
    renderedPrompt.value[activePromptTab.value] = prompt.content;
    ElMessage.info(t('promptEngineering.promptLoaded', { name: prompt.name }));
  } else {
    ElMessage.warning(t('promptEngineering.selectModelFirst'));
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
.prompt-engineering {
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