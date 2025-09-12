<template>
  <el-card class="model-selector">
    <template #header>
      <div class="card-header">
        <span>模型选择</span>
      </div>
    </template>
    <el-select
        :model-value="props.selectedServer"
        placeholder="选择服务"
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
        placeholder="选择模型"
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
      <el-button @click="() => emit('load-model')" style="width: 100%">加载模型</el-button>
    </div>

    <div style="margin-top: 20px">
      <h4>参数设置</h4>
      <el-form :model="props.modelParams" label-width="80px" size="small">
        <el-form-item label="温度">
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
        <el-form-item label="上下文">
          <el-input-number
              :model-value="props.modelParams.context"
              :min="1"
              :max="32768"
              @update:modelValue="val => emit('update:modelParams', { ...props.modelParams, context: val })"/>
        </el-form-item>
        <el-form-item label="预测数">
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
        <el-form-item label="重复惩罚">
          <el-input-number
              :model-value="props.modelParams.repeatPenalty"
              :min="0.1"
              :max="2"
              :step="0.1"
              @update:modelValue="val => emit('update:modelParams', { ...props.modelParams, repeatPenalty: val })"/>
        </el-form-item>
        <el-form-item label="输出方式">
          <el-select
              :model-value="props.modelParams.outputMode"
              placeholder="选择输出方式"
              @update:modelValue="val => emit('update:modelParams', { ...props.modelParams, outputMode: val })">
            <el-option label="流式输出" value="stream"/>
            <el-option label="阻塞输出" value="blocking"/>
          </el-select>
        </el-form-item>
      </el-form>
      <div style="margin-top: 10px">
        <el-button @click="() => emit('save-model-params')" type="primary" size="small">保存参数</el-button>
        <el-button @click="() => emit('reset-model-params')" size="small">重置</el-button>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">

interface Model {
  name: string;
  size: number;
  modified_at: string;
}

interface Server {
  id: string;
  name: string;
  baseUrl: string;
  apiKey: string;
  isActive: boolean;
  testStatus: string;
  type: string;
}

interface ModelParams {
  temperature: number;
  topP: number;
  context: number;
  numPredict: number;
  topK: number;
  repeatPenalty: number;
  outputMode: 'stream' | 'blocking';
}

const props = defineProps<{
  localModels: Model[];
  selectedModel: string;
  availableServers: Server[];
  selectedServer: string;
  modelParams: ModelParams;
}>();

const emit = defineEmits<{
  (e: 'update:selectedModel', value: string): void;
  (e: 'update:selectedServer', value: string): void;
  (e: 'update:modelParams', value: ModelParams): void;
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
