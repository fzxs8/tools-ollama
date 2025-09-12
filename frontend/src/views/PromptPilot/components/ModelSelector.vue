<template>
  <div class="model-selector-horizontal">
    <el-select
        v-model="props.selectedServer"
        placeholder="选择服务"
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
        placeholder="选择模型 (最多3个)"
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
import {GetActiveServer, GetServers, ListModelsByServer, SetActiveServer} from '../../../../wailsjs/go/main/App'

// 定义接口
interface Model {
  name: string
  size: number
  modified_at: string
}

interface Server {
  id: string
  name: string
  base_url: string
  api_key: string
  is_active: boolean
  test_status: string
  type: string
}

// 定义 Props 和 Emits
const props = defineProps<{
  selectedServer: string
  selectedModels: string[]
}>()

const emit = defineEmits<{
  (e: 'update:selectedServer', value: string): void
  (e: 'update:selectedModels', value: string[]): void
}>()

// 定义响应式状态
const availableServers = ref<Server[]>([])
const availableModels = ref<Model[]>([])

// 更新选中的服务
const updateSelectedServer = (value: string) => {
  emit('update:selectedServer', value)
}

// 更新选中的模型
const updateSelectedModels = (value: string[]) => {
  emit('update:selectedModels', value)
}

// 加载可用服务
const loadAvailableServers = async () => {
  try {
    const allServers = await GetServers();
    availableServers.value = allServers as Server[];

    if (allServers.length === 0) {
      ElMessage.warning('没有配置任何Ollama服务。请在“服务设置”页面添加一个。');
      availableModels.value = [];
      updateSelectedServer('');
      return;
    }

    const activeServer = await GetActiveServer();
    const activeServerExists = activeServer && allServers.some(s => s.id === activeServer.id);

    let serverToSelect = '';
    if (activeServerExists) {
      serverToSelect = activeServer.id;
    } else {
      // Default to the first server in the list if no active one is found
      serverToSelect = allServers[0].id;
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

// 根据服务加载模型
const loadModelsForServer = async (serverId: string) => {
  if (!serverId) return
  try {
    availableModels.value = await ListModelsByServer(serverId)
  } catch (error: any) {
    console.error(`获取模型列表失败 (服务ID: ${serverId}):`, error)
    ElMessage.error('获取模型列表失败: ' + (error.message || error))
    availableModels.value = []
  }
}

// 当服务选择变化时
const onServerChange = (serverId: string) => {
  updateSelectedServer(serverId)
  // 清空已选模型
  updateSelectedModels([])
  loadModelsForServer(serverId)
}

// 监视 selectedServer 的变化
watch(() => props.selectedServer, (newServerId, oldServerId) => {
  if (newServerId && newServerId !== oldServerId) {
    loadModelsForServer(newServerId)
  }
})

// 组件挂载时加载
onMounted(async () => {
  await loadAvailableServers()
  // 如果初始有 selectedServer，则加载模型
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
  width: 220px; /* 设置固定宽度 */
}
</style>
