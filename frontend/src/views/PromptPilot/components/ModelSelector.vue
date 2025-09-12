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
import {GetActiveServer, GetServers, ListModelsByServer} from '../../../../wailsjs/go/main/App'

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
    // const localBaseUrl = await GetOllamaServerConfig()
    // const localServer: Server = {
    //   id: 'local',
    //   name: '本地服务',
    //   base_url: localBaseUrl,
    //   api_key: '',
    //   is_active: false, // 稍后会根据 GetActiveServer 更新
    //   test_status: '',
    //   type: 'local'
    // }
    //
    let remoteServers: Server[] = []
    try {
      const remoteList: any[] = await GetServers()
      if (remoteList) {
        remoteServers = remoteList.map(server => ({
          id: server.id || server.ID,
          name: server.name || server.Name,
          base_url: server.baseUrl || server.base_url || server.BaseURL,
          api_key: server.apiKey || server.api_key || server.APIKey,
          is_active: server.isActive !== undefined ? server.isActive : (server.is_active !== undefined ? server.is_active : server.IsActive),
          test_status: server.testStatus || server.test_status || server.TestStatus || '',
          type: server.type || server.Type || 'remote'
        }))
      }
    } catch (remoteError) {
      console.error('获取远程服务器列表失败:', remoteError)
    }

    const allServers = [...remoteServers]
    availableServers.value = allServers

    try {
      const activeServer = await GetActiveServer()
      if (activeServer && activeServer.id && allServers.some(s => s.id === activeServer.id)) {
        updateSelectedServer(activeServer.id)
      } else {
        updateSelectedServer('local')
      }
    } catch (e) {
      updateSelectedServer('local')
    }

  } catch (error) {
    console.error('加载服务配置失败:', error)
    const localServer: Server = {
      id: 'local',
      name: '本地服务',
      base_url: '',
      api_key: '',
      is_active: true,
      test_status: '',
      type: 'local'
    };
    availableServers.value = [localServer]
    updateSelectedServer('local')
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
