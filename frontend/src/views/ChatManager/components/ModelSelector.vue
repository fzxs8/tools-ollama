<template>
  <el-card class="model-selector">
    <template #header>
      <div class="card-header">
        <span>{{ t('chatManager.modelSelection') }}</span>
      </div>
    </template>
    <el-select
        :model-value="props.selectedServer"
        :placeholder="t('modelManager.selectServer')"
        style="width: 100%; margin-bottom: 10px;"
        @update:modelValue="val => emit('update:selectedServer', val)"
        @change="() => emit('server-change')">
      <el-option
          v-for="server in props.availableServers"
          :key="server.id"
          :label="server.name"
          :value="server.id"
      />
    </el-select>
    <el-select
        :model-value="props.selectedModel"
        :placeholder="t('modelManager.selectModel')"
        style="width: 100%"
        @update:modelValue="val => emit('update:selectedModel', val)">
      <el-option
          v-for="model in props.localModels"
          :key="model.name"
          :label="model.name"
          :value="model.name"
      />
    </el-select>

    <div style="margin-top: 20px">
      <el-button @click="() => emit('load-model')" style="width: 100%">{{ t('chatManager.loadModel') }}</el-button>
    </div>

    <div style="margin-top: 20px">
      <h4>{{ t('modelManager.modelParameters') }}</h4>
      <el-form :model="props.modelParams" label-width="80px" size="small">
        <el-form-item :label="t('modelManager.temperature')">
          <el-slider
              :model-value="props.modelParams.temperature"
              :min="0"
              :max="1"
              :step="0.1"
              @update:modelValue="val => emit('update:modelParams', { ...props.modelParams, temperature: val })"/>
        </el-form-item>
        <el-form-item label="Top P">
          <el-slider
              :model-value="props.modelParams.topP"
              :min="0"
              :max="1"
              :step="0.1"
              @update:modelValue="val => emit('update:modelParams', { ...props.modelParams, topP: val })"/>
        </el-form-item>
        <el-form-item :label="t('modelManager.context')">
          <el-input-number
              :model-value="props.modelParams.context"
              :min="1"
              :max="32768"
              @update:modelValue="val => emit('update:modelParams', { ...props.modelParams, context: val })"/>
        </el-form-item>
        <el-form-item :label="t('chatManager.maxTokens')">
          <el-input-number
              :model-value="props.modelParams.numPredict"
              :min="1"
              :max="4096"
              @update:modelValue="val => emit('update:modelParams', { ...props.modelParams, numPredict: val })"/>
        </el-form-item>
        <el-form-item label="Top K">
          <el-input-number
              :model-value="props.modelParams.topK"
              :min="1"
              :max="100"
              @update:modelValue="val => emit('update:modelParams', { ...props.modelParams, topK: val })"/>
        </el-form-item>
        <el-form-item :label="t('chatManager.repeatPenalty')">
          <el-input-number
              :model-value="props.modelParams.repeatPenalty"
              :min="0.1"
              :max="2"
              :step="0.1"
              @update:modelValue="val => emit('update:modelParams', { ...props.modelParams, repeatPenalty: val })"/>
        </el-form-item>
        <el-form-item :label="t('chatManager.outputMode')">
          <el-select
              :model-value="props.modelParams.outputMode"
              :placeholder="t('chatManager.selectOutputMode')"
              @update:modelValue="val => emit('update:modelParams', { ...props.modelParams, outputMode: val })">
            <el-option :label="t('chatManager.streamOutput')" value="stream"/>
            <el-option :label="t('chatManager.blockingOutput')" value="blocking"/>
          </el-select>
        </el-form-item>
      </el-form>
      <div style="margin-top: 10px">
        <el-button @click="() => emit('save-model-params')" type="primary" size="small">{{ t('modelManager.saveParameters') }}</el-button>
        <el-button @click="() => emit('reset-model-params')" size="small">{{ t('modelManager.resetParameters') }}</el-button>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import {types} from "../../../../wailsjs/go/models";
import Model = types.Model;
import OllamaServerConfig = types.OllamaServerConfig;
import {ModelParams} from "../../../classes/types";

const { t } = useI18n();

const props = defineProps<{
  localModels: Model[];
  selectedModel: string;
  availableServers: OllamaServerConfig[];
  selectedServer: string;
  modelParams: ModelParams; // Use imported ModelParams
}>();

const emit = defineEmits<{
  (e: 'update:selectedModel', value: string): void;
  (e: 'update:selectedServer', value: string): void;
  (e: 'update:modelParams', value: ModelParams): void; // Use imported ModelParams
  (e: 'load-model'): void;
  (e: 'save-model-params'): void;
  (e: 'reset-model-params'): void;
  (e: 'server-change'): void;
}>();

</script>

<style scoped>
.model-selector {
  height: 100%;
}

.model-selector :deep(.el-card__body) {
  height: calc(100% - 60px);
  overflow-y: auto;
}
</style>
