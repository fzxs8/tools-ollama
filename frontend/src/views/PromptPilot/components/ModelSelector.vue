<template>
  <div class="model-selector-horizontal">
    <el-select
        v-model="props.selectedServer"
        :placeholder="t('promptPilot.selectService')"
        class="selector-item"
        popper-class="left-aligned-dropdown"
        @change="onServerChange"
    >
      <el-option
          v-for="server in availableServers"
          :key="server.id"
          :label="server.name"
          :value="server.id"
      />
    </el-select>
    <el-select
        v-model="props.selectedModels"
        multiple
        :multiple-limit="3"
        collapse-tags
        :placeholder="t('promptPilot.selectModel')"
        class="selector-item"
        popper-class="left-aligned-dropdown"
        @update:modelValue="updateSelectedModels"
    >
      <el-option
          v-for="model in availableModels"
          :key="model.name"
          :label="model.name"
          :value="model.name"
      />
    </el-select>
  </div>
</template>

<script setup lang="ts">
import {onMounted, ref, watch} from 'vue'
import {ElMessage} from 'element-plus'
import { useI18n } from 'vue-i18n'
import {GetActiveServer, GetServers, ListModelsByServer, SetActiveServer} from '../../../../wailsjs/go/main/App'
import {types} from "../../../../wailsjs/go/models";
import Model = types.Model;
import OllamaServerConfig = types.OllamaServerConfig;

const { t } = useI18n();

// Define Props and Emits
const props = defineProps<{
  selectedServer: string
  selectedModels: string[]
}>()

const emit = defineEmits<{
  (e: 'update:selectedServer', value: string): void
  (e: 'update:selectedModels', value: string[]): void
}>()

// Define reactive state
const availableServers = ref<OllamaServerConfig[]>([])
const availableModels = ref<Model[]>([])

// Update selected server
const updateSelectedServer = (value: string) => {
  emit('update:selectedServer', value)
}

// Update selected models
const updateSelectedModels = (value: string[]) => {
  emit('update:selectedModels', value)
}

// Load available servers
const loadAvailableServers = async () => {
  try {
    availableServers.value = await GetServers();

    if (availableServers.value.length === 0) {
      ElMessage.warning('没有配置任何Ollama服务。请在“服务设置”页面添加一个。');
      availableModels.value = [];
      updateSelectedServer('');
      return;
    }

    const activeServer = await GetActiveServer();
    const activeServerExists = activeServer && availableServers.value.some(s => s.id === activeServer.id);

    let serverToSelect = '';
    if (activeServerExists) {
      serverToSelect = activeServer.id;
    } else {
      // Default to the first server in the list if no active one is found
      serverToSelect = availableServers.value[0].id;
      await SetActiveServer(serverToSelect);
    }
    updateSelectedServer(serverToSelect);

  } catch (error) {
    console.error('加载服务配置失败:', error);
    ElMessage.error('加载服务列表失败: ' + (error as Error).message);
    availableServers.value = [];
    updateSelectedServer('');
  }
}

// Load models for server
const loadModelsForServer = async (serverId: string) => {
  if (!serverId) return
  try {
    availableModels.value = await ListModelsByServer(serverId)
  } catch (error: any) {
    console.error(`Failed to get model list (Service ID: ${serverId}):`, error)
    ElMessage.error('Failed to get model list: ' + (error.message || error))
    availableModels.value = []
  }
}

// When server selection changes
const onServerChange = (serverId: string) => {
  updateSelectedServer(serverId)
  // Clear selected models
  updateSelectedModels([])
  loadModelsForServer(serverId)
}

// Watch for selectedServer changes
watch(() => props.selectedServer, (newServerId, oldServerId) => {
  if (newServerId && newServerId !== oldServerId) {
    loadModelsForServer(newServerId)
  }
})

// Load on component mount
onMounted(async () => {
  await loadAvailableServers()
  // If there's an initial selectedServer, load models
  if (props.selectedServer) {
    await loadModelsForServer(props.selectedServer)
  }
})
</script>

<!-- Global style for the dropdown, placed in the component file for better organization -->
<style>
.left-aligned-dropdown .el-select-dropdown__item {
  justify-content: flex-start !important;
}
</style>

<style scoped>
.model-selector-horizontal {
  display: flex;
  gap: 10px;
}

.selector-item {
  width: 220px; /* Set fixed width */
}
</style>
