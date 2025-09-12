<template>
  <el-drawer
    :model-value="props.visible"
    title="选择系统提示词"
    direction="rtl"
    size="40%"
    @update:modelValue="closeDrawer"
  >
    <div class="system-prompt-selector-container">
      <el-select
        v-model="selectedPromptId"
        placeholder="请选择一个已保存的提示词作为系统提示词"
        style="width: 100%;"
        @change="onPromptSelect"
        filterable
        clearable
      >
        <el-option
          v-for="prompt in props.prompts"
          :key="prompt.id"
          :label="prompt.name"
          :value="prompt.id"
        />
      </el-select>

      <div v-if="selectedPromptContent" class="prompt-preview-area">
        <div class="preview-title">内容预览</div>
        <pre class="prompt-content-preview">{{ selectedPromptContent }}</pre>
      </div>

      <div class="drawer-footer">
        <el-button @click="closeDrawer">取消</el-button>
        <el-button type="primary" @click="applySystemPrompt" :disabled="!selectedPromptId">应用</el-button>
      </div>
    </div>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import {types} from "../../../../wailsjs/go/models";
import Prompt = types.Prompt;


const props = defineProps<{
  visible: boolean;
  prompts: Prompt[];
  currentSystemPromptId?: string;
}>();

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'apply', promptId: string | undefined): void;
}>();

const selectedPromptId = ref<string | undefined>(props.currentSystemPromptId);

// 监听外部传入的ID变化
watch(() => props.currentSystemPromptId, (newId) => {
  selectedPromptId.value = newId;
});

const selectedPromptContent = computed(() => {
  if (!selectedPromptId.value) return '';
  const selected = props.prompts.find(p => p.id === selectedPromptId.value);
  return selected ? selected.content : '';
});

const closeDrawer = () => {
  emit('update:visible', false);
};

const onPromptSelect = (promptId: string | undefined) => {
  selectedPromptId.value = promptId;
};

const applySystemPrompt = () => {
  emit('apply', selectedPromptId.value);
  closeDrawer();
};

</script>

<style scoped>
.system-prompt-selector-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 0 20px;
}

.prompt-preview-area {
  margin-top: 20px;
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.preview-title {
  font-size: 14px;
  color: #606266;
  margin-bottom: 10px;
  font-weight: 500;
}

pre.prompt-content-preview {
  flex: 1;
  font-size: 13px;
  color: #555;
  background-color: #f8f9fa;
  padding: 15px;
  border-radius: 6px;
  white-space: pre-wrap;
  line-height: 1.6;
  border: 1px solid #eee;
  margin: 0;
  overflow-y: auto;
}

.drawer-footer {
  text-align: right;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #e0e0e0;
}
</style>
