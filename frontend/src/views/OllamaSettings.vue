<template>
  <div class="ollama-settings">
    <el-card class="settings-card">
      <template #header>
        <div class="card-header">
          <span>Ollama 服务配置</span>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="本地服务配置" name="local">
          <el-form :model="localConfig" label-width="120px" style="max-width: 600px; margin-top: 20px;">
            <el-form-item label="本地服务地址">
              <el-input v-model="localConfig.baseUrl" placeholder="请输入本地 Ollama 服务地址"></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveLocalConfig">保存配置</el-button>
              <el-button @click="testLocalConnection" :loading="testingLocal">测试连接</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="远程服务管理" name="remote">
          <div style="margin-bottom: 20px;">
            <el-button type="primary" @click="openAddRemoteDialog">添加远程服务</el-button>
          </div>

          <el-table :data="remoteServers" style="width: 100%" empty-text="暂无数据">
            <el-table-column prop="name" label="服务名称"/>
            <el-table-column prop="base_url" label="服务地址"/>
            <el-table-column prop="test_status" label="连接状态">
              <template #default="scope">
                <el-tag v-if="scope.row.test_status === 'success'" type="success">连接成功</el-tag>
                <el-tag v-else-if="scope.row.test_status === 'failed'" type="danger">连接失败</el-tag>
                <el-tag v-else type="info">未知</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <template #default="scope">
                <el-button size="small" @click="testRemoteConnection(scope.row)">测试</el-button>
                <el-button size="small" @click="editRemoteServer(scope.row)">编辑</el-button>
                <el-button size="small" type="danger" @click="deleteRemoteServer(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 添加/编辑远程服务对话框 -->
    <el-dialog v-model="remoteDialogVisible" :title="editingRemoteServer ? '编辑远程服务' : '添加远程服务'"
               width="500px">
      <el-form :model="remoteForm" label-width="100px">
        <el-form-item label="服务名称" required>
          <el-input v-model="remoteForm.name" placeholder="请输入服务名称"></el-input>
        </el-form-item>
        <el-form-item label="服务地址" required>
          <el-input v-model="remoteForm.base_url"
                    placeholder="请输入服务地址，例如：http://192.168.1.100:11434"></el-input>
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="remoteForm.api_key" placeholder="请输入 API Key（可选）" type="password"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="remoteDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveRemoteServer">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {onMounted, reactive, ref} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {
  AddRemoteServer,
  DeleteRemoteServer,
  GetOllamaServerConfig,
  GetRemoteServers,
  SaveOllamaServerConfig,
  TestOllamaServer,
  UpdateRemoteServer
} from '../../wailsjs/go/main/App'

interface RemoteServer {
  id: string
  name: string
  base_url: string
  api_key: string
  is_active: boolean
  test_status: string
}

// 远程服务器列表
const remoteServers = ref<RemoteServer[]>([])

// 活动标签页
const activeTab = ref('local')

// 本地配置
const localConfig = reactive({
  baseUrl: 'http://localhost:11434'
})

// 远程服务器表单
const remoteForm = reactive({
  id: '',
  name: '',
  base_url: '',
  api_key: ''
})

// 是否正在编辑远程服务器
const editingRemoteServer = ref(false)

// 远程服务器对话框可见性
const remoteDialogVisible = ref(false)

// 是否正在测试本地连接
const testingLocal = ref(false)

// 加载本地配置
const loadLocalConfig = async () => {
  try {
    const config = await GetOllamaServerConfig()
    localConfig.baseUrl = config
  } catch (error) {
    console.error('加载本地配置失败:', error)
  }
}

// 保存本地配置
const saveLocalConfig = async () => {
  try {
    await SaveOllamaServerConfig(localConfig.baseUrl)
    ElMessage.success('本地配置已保存')
  } catch (error) {
    ElMessage.error('保存失败: ' + (error as Error).message)
  }
}

// 测试本地连接
const testLocalConnection = async () => {
  testingLocal.value = true
  try {
    await TestOllamaServer(localConfig.baseUrl)
    ElMessage.success('本地服务连接测试成功')
  } catch (error) {
    ElMessage.error('本地服务连接测试失败: ' + (error as Error).message)
  } finally {
    testingLocal.value = false
  }
}

// 打开添加远程服务器对话框
const openAddRemoteDialog = () => {
  editingRemoteServer.value = false
  remoteForm.id = ''
  remoteForm.name = ''
  remoteForm.base_url = ''
  remoteForm.api_key = ''
  remoteDialogVisible.value = true
}

// 编辑远程服务器
const editRemoteServer = (server: RemoteServer) => {
  editingRemoteServer.value = true
  remoteForm.id = server.id
  remoteForm.name = server.name
  remoteForm.base_url = server.base_url
  remoteForm.api_key = server.api_key
  remoteDialogVisible.value = true
}

// 保存远程服务器
const saveRemoteServer = async () => {
  if (!remoteForm.name || !remoteForm.base_url) {
    ElMessage.warning('请填写服务名称和服务地址')
    return
  }

  try {
    if (editingRemoteServer.value) {
      // 更新现有服务器
      await UpdateRemoteServer({
        id: remoteForm.id,
        name: remoteForm.name,
        base_url: remoteForm.base_url,
        api_key: remoteForm.api_key,
        is_active: false,
        test_status: 'unknown'
      })
      ElMessage.success('服务器已更新')
    } else {
      // 添加新服务器
      await AddRemoteServer({
        id: Date.now().toString(),
        name: remoteForm.name,
        base_url: remoteForm.base_url,
        api_key: remoteForm.api_key,
        is_active: false,
        test_status: 'unknown'
      })
      ElMessage.success('服务器已添加')
    }

    remoteDialogVisible.value = false
    loadRemoteServers()
  } catch (error) {
    ElMessage.error((editingRemoteServer.value ? '更新' : '添加') + '失败: ' + (error as Error).message)
  }
}

// 删除远程服务器
const deleteRemoteServer = (server: RemoteServer) => {
  ElMessageBox.confirm(`确定要删除 ${server.name} 吗？`, '确认删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await DeleteRemoteServer(server.id)
      ElMessage.success('服务器已删除')
      loadRemoteServers()
    } catch (error) {
      ElMessage.error('删除失败: ' + (error as Error).message)
    }
  })
}

// 测试远程连接
const testRemoteConnection = async (server: RemoteServer) => {
  try {
    await TestOllamaServer(server.base_url)
    ElMessage.success(`${server.name} 连接测试成功`)

    // 更新测试状态
    const index = remoteServers.value.findIndex(s => s.id === server.id)
    if (index !== -1) {
      remoteServers.value[index].test_status = 'success'
      await UpdateRemoteServer(remoteServers.value[index])
    }
  } catch (error) {
    ElMessage.error(`${server.name} 连接测试失败: ${(error as Error).message}`)

    // 更新测试状态
    const index = remoteServers.value.findIndex(s => s.id === server.id)
    if (index !== -1) {
      remoteServers.value[index].test_status = 'failed'
      await UpdateRemoteServer(remoteServers.value[index])
    }
  }
}

// 加载远程服务器
const loadRemoteServers = async () => {
  try {
    const servers = await GetRemoteServers()
    remoteServers.value = servers
  } catch (error) {
    console.error('加载远程服务器失败:', error)
  }
}

// 初始化数据
const loadData = async () => {
  await loadLocalConfig()
  await loadRemoteServers()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.ollama-settings {
  padding: 20px;
  height: 100%;
  box-sizing: border-box;
}

.settings-card {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.settings-card :deep(.el-card__body) {
  height: calc(100% - 60px);
  overflow-y: auto;
}
</style>
