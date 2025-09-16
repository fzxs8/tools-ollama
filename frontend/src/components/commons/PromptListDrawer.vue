<template>
  <el-drawer
    :model-value="props.visible"
    :title="drawerTitle"
    direction="rtl"
    size="50%"
    :close-on-click-modal="false"
    @update:modelValue="closeDrawer"
  >
    <div class="common-prompts-container">
      <el-empty v-if="!prompts || prompts.length === 0" :description="emptyDescription" />
      
      <div v-else class="prompts-list">
        <div 
          v-for="prompt in prompts" 
          :key="prompt.id"
          class="prompt-item"
          :class="{ 'selectable': props.mode === 'select', 'selected': props.mode === 'select' && props.selectedId === prompt.id }"
          @click="handleItemClick(prompt)"
        >
          <div class="prompt-header">
            <div class="prompt-title">{{ prompt.name }}</div>
            <div class="prompt-actions">
              <!-- Management mode buttons -->
              <template v-if="props.mode === 'manage'">
                <el-button type="primary" link @click.stop="previewPrompt(prompt)" size="small">{{ t('common.preview') }}</el-button>
                <el-button type="primary" link @click.stop="editPrompt(prompt)" size="small">{{ t('common.edit') }}</el-button>
                <el-button type="danger" link @click.stop="deletePrompt(prompt)" size="small">{{ t('common.delete') }}</el-button>
              </template>
              <!-- Selection mode buttons -->
              <template v-if="props.mode === 'select'">
                <el-button v-if="props.selectedId === prompt.id" type="success" link size="small">{{ t('common.selected') }}</el-button>
                <el-button v-else type="primary" link @click.stop="selectPrompt(prompt)" size="small">{{ t('common.select') }}</el-button>
              </template>
            </div>
          </div>
          
          <div class="prompt-meta">
            <div class="prompt-description" v-if="prompt.description">{{ prompt.description }}</div>
            <div class="prompt-tags" v-if="prompt.tags && prompt.tags.length > 0">
              <el-tag v-for="tag in prompt.tags" :key="tag" size="small" type="success" style="margin-right: 5px; margin-bottom: 5px;">{{ tag }}</el-tag>
            </div>
            <div class="prompt-models">
              <el-tag v-for="model in prompt.models" :key="model" type="info" size="small" style="margin-right: 5px; margin-bottom: 5px;">{{ model }}</el-tag>
            </div>
            <div class="prompt-time">{{ t('modelManager.updatedAt') }} {{ formatTime(prompt.updatedAt) }}</div>
          </div>
          
          <div class="prompt-content-preview-wrapper">
            <pre class="prompt-content-preview">{{ truncateContent(prompt.content, 100) }}</pre>
            <el-button class="preview-copy-btn" type="primary" link @click.stop="copyPromptContent(prompt.content)">{{ t('common.copy') }}</el-button>
          </div>
        </div>
      </div>
    </div>
  </el-drawer>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { ElMessage } from "element-plus";
import { useI18n } from 'vue-i18n';
import {types} from "../../../wailsjs/go/models";
import Prompt = types.Prompt;

// Type definitions
const props = defineProps<{
  visible: boolean;
  prompts: Prompt[];
  mode: 'manage' | 'select';
  selectedId?: string;
}>();

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'preview', prompt: Prompt): void;
  (e: 'edit', prompt: Prompt): void;
  (e: 'delete', prompt: Prompt): void;
  (e: 'select', prompt: Prompt): void;
}>();

const { t } = useI18n();

const drawerTitle = computed(() => props.mode === 'manage' ? t('promptEngineering.myPrompts') : t('promptEngineering.selectSystemPrompt'));
const emptyDescription = computed(() => props.mode === 'manage' ? t('promptEngineering.noPrompts') : t('promptEngineering.noSystemPrompts'));

const closeDrawer = () => {
  emit('update:visible', false);
};

const formatTime = (timestamp: number) => {
  const date = new Date(timestamp);
  return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')} ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`;
};

const truncateContent = (content: string, length: number) => {
  if (!content) return '';
  return content.length > length ? content.substring(0, length) + '...' : content;
};

const copyPromptContent = (content: string) => {
  if (content) {
    navigator.clipboard.writeText(content).then(() => {
      ElMessage.success(t('messages.contentCopied'));
    }).catch(err => {
      ElMessage.error(t('messages.copyFailed') + ': ' + err);
    });
  }
};

const handleItemClick = (prompt: Prompt) => {
  if (props.mode === 'select') {
    selectPrompt(prompt);
  }
};

const previewPrompt = (prompt: Prompt) => emit('preview', prompt);
const editPrompt = (prompt: Prompt) => emit('edit', prompt);
const deletePrompt = (prompt: Prompt) => emit('delete', prompt);
const selectPrompt = (prompt: Prompt) => emit('select', prompt);

</script>

<style scoped>
.common-prompts-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.prompts-list {
  flex: 1;
  overflow-y: auto;
  padding-right: 10px;
}

.prompt-item {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 15px;
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  transition: all 0.2s ease-in-out;
}

.prompt-item.selectable {
  cursor: pointer;
}

.prompt-item.selectable:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: #c0c4cc;
}

.prompt-item.selected {
  border-color: #409eff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
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
}

.prompt-description {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
  line-height: 1.4;
  text-align: left;
}

.prompt-tags, .prompt-models {
  margin-bottom: 8px;
  text-align: left;
}

.prompt-time {
  font-size: 12px;
  color: #999;
  text-align: left;
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
  right: 5px;
  padding: 2px;
  font-size: 12px;
}
</style>
