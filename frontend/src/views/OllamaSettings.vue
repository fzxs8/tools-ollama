<template>
  <div class="ollama-settings">
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>Ollama 服务配置</span>
              <div>
                <el-button type="primary" @click="openAddDrawer">添加服务</el-button>
              </div>
            </div>
          </template>
          
          <div class="server-config-content">
            <!-- 服务列表 -->
            <div class="server-list-section">
              <h3>服务列表</h3>
              <el-table :data="allServers" style="width: 100%" :fit="true">
                <el-table-column prop="name" label="名称" :min-width="100" />
                <el-table-column prop="base_url" label="服务地址" :min-width="150" />
                <el-table-column prop="type" label="类型" :width="80">
                  <template #default="scope">
                    <el-tag :type="scope.row.type === 'local' ? 'primary' : 'success'">
                      {{ scope.row.type === 'local' ? '本地' : '远程' }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="test_status" label="测试状态" :width="120">
                  <template #default="scope">
                    <el-tag :type="scope.row.test_status === 'success' ? 'success' : 
                             scope.row.test_status === 'failed' ? 'danger' : 'info'">
                      {{ scope.row.test_status === 'success' ? '连接成功' : 
                         scope.row.test_status === 'failed' ? '连接失败' : '未测试' }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="操作" :width="120">
                  <template #default="scope">
                    <el-dropdown @command="handleCommand">
                      <el-button size="small">
                        操作<i class="el-icon-arrow-down el-icon--right"></i>
                      </el-button>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item :command="{action: 'setActive', server: scope.row}">
                            设置默认
                          </el-dropdown-item>
                          <el-dropdown-item :command="{action: 'edit', server: scope.row}">
                            编辑
                          </el-dropdown-item>
                          <el-dropdown-item :command="{action: 'test', server: scope.row}">
                            测试
                          </el-dropdown-item>
                          <el-dropdown-item 
                            v-if="scope.row.type === 'local'"
                            :command="{action: 'toggle', server: scope.row}">
                            {{ scope.row.is_running ? '停止' : '启动' }}
                          </el-dropdown-item>
                          <el-dropdown-item 
                            v-if="scope.row.type !== 'local'"
                            :command="{action: 'remove', server: scope.row}"
                            :disabled="scope.row.is_active">
                            删除
                          </el-dropdown-item>
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 添加/编辑服务抽屉 -->
    <el-drawer
      v-model="drawerVisible"
      :title="editingServer ? '编辑服务' : '添加服务'"
      direction="rtl"
      size="40%"
    >
      <el-form :model="serverForm" label-width="100px">
        <el-form-item label="服务类型">
          <el-radio-group v-model="serverForm.type" :disabled="editingServer">
            <el-radio label="local">本地服务</el-radio>
            <el-radio label="remote">远程服务</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="名称">
          <el-input v-model="serverForm.name" placeholder="例如：开发环境" />
        </el-form-item>
        
        <el-form-item label="服务地址">
          <el-input v-model="serverForm.base_url" placeholder="http://example.com:11434" />
        </el-form-item>
        
        <el-form-item label="API密钥" v-if="serverForm.type === 'remote'">
          <el-input v-model="serverForm.api_key" type="password" placeholder="可选" show-password />
        </el-form-item>
        
        <el-form-item v-if="serverForm.type === 'local'">
          <el-checkbox v-model="serverForm.autoStart">自动启动</el-checkbox>
        </el-form-item>
        
        <el-form-item>
          <el-button 
            type="primary" 
            @click="saveServer" 
            :loading="isTestingConnection"
          >
            {{ editingServer ? '更新' : '添加' }}
          </el-button>
          <el-button @click="cancelEdit">取消</el-button>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  SaveOllamaServerConfig, 
  GetOllamaServerConfig,
  SaveRemoteServers,
  GetRemoteServers,
  AddRemoteServer,
  UpdateRemoteServer,
  DeleteRemoteServer,
  SetActiveServer,
  TestOllamaServer
} from '../../wailsjs/go/main/App'

// 抽屉可见性
const drawerVisible = ref(false)

// 远程服务器列表
const remoteServers = ref<any[]>([])

// 服务表单
const serverForm = reactive({
  id: '',
  name: '',
  base_url: '',
  api_key: '',
  is_active: false,
  test_status: 'unknown',
  type: 'local',
  autoStart: true
})

// 是否正在编辑服务器
const editingServer = ref(false)

// 是否正在测试连接
const isTestingConnection = ref(false)

// 本地服务状态
const localServerState = reactive({
  test_status: 'unknown' as 'unknown' | 'success' | 'failed', // unknown, success, failed
  is_running: false
})

// 计算属性：所有服务（包括本地服务）
const allServers = computed(() => {
  // 本地服务始终显示在列表中
  const localServer = {
    id: 'local',
    name: '本地服务',
    base_url: 'http://localhost:11434',
    api_key: '',
    is_active: false,
    test_status: localServerState.test_status,
    type: 'local',
    is_running: localServerState.is_running
  }
  
  // 合并本地服务和远程服务
  return [localServer, ...remoteServers.value]
})

// 打开添加抽屉
const openAddDrawer = () => {
  // 重置表单
  serverForm.id = ''
  serverForm.name = ''
  serverForm.base_url = ''
  serverForm.api_key = ''
  serverForm.is_active = false
  serverForm.test_status = 'unknown'
  serverForm.type = 'local'
  serverForm.autoStart = true
  editingServer.value = false
  drawerVisible.value = true
}

// 编辑服务器
const editServer = (server: any) => {
  serverForm.id = server.id
  serverForm.name = server.name
  serverForm.base_url = server.base_url
  serverForm.api_key = server.api_key
  serverForm.is_active = server.is_active
  serverForm.test_status = server.test_status || 'unknown'
  serverForm.type = server.type
  serverForm.autoStart = server.autoStart || true
  editingServer.value = true
  drawerVisible.value = true
}

// 切换本地服务状态
const toggleLocalService = (server: any) => {
  ElMessage.info(server.is_running ? '正在停止本地Ollama服务...' : '正在启动本地Ollama服务...')
  // 模拟状态切换
  localServerState.is_running = !localServerState.is_running
  setTimeout(() => {
    ElMessage.success(localServerState.is_running ? '本地Ollama服务已启动' : '本地Ollama服务已停止')
  }, 1000)
}

// 处理下拉菜单命令
const handleCommand = (command: {action: string, server: any}) => {
  switch (command.action) {
    case 'setActive':
      setActiveServer(command.server)
      break
    case 'edit':
      editServer(command.server)
      break
    case 'test':
      testConnection(command.server)
      break
    case 'toggle':
      toggleLocalService(command.server)
      break
    case 'remove':
      removeRemoteServer(command.server)
      break
  }
}

// 测试连接
const testConnection = async (server: any, isAutoTest = false) => {
  try {
    await TestOllamaServer(server.base_url)
    
    // 只有手动测试时才显示提示信息
    if (!isAutoTest) {
      ElMessage.success(`${server.name} 连接测试成功`)
    }
    
    // 更新测试状态
    if (server.type === 'local') {
      // 更新本地服务状态
      localServerState.test_status = 'success'
    } else {
      // 更新远程服务状态
      const servers = [...remoteServers.value]
      const index = servers.findIndex((s: any) => s.id === server.id)
      if (index !== -1) {
        servers[index].test_status = 'success'
        remoteServers.value = servers
        
        // 同时更新后端存储的状态
        servers[index].test_status = 'success'
        await UpdateRemoteServer(servers[index])
      }
    }
    
    return 'success'
  } catch (error) {
    // 只有手动测试时才显示提示信息
    if (!isAutoTest) {
      ElMessage.error(`${server.name} 连接测试失败: ${(error as Error).message}`)
    }
    
    // 更新测试状态
    if (server.type === 'local') {
      // 更新本地服务状态
      localServerState.test_status = 'failed'
    } else {
      // 更新远程服务状态
      const servers = [...remoteServers.value]
      const index = servers.findIndex((s: any) => s.id === server.id)
      if (index !== -1) {
        servers[index].test_status = 'failed'
        remoteServers.value = servers
        
        // 同时更新后端存储的状态
        servers[index].test_status = 'failed'
        await UpdateRemoteServer(servers[index])
      }
    }
    
    return 'failed'
  }
}

// 保存服务器
const saveServer = async () => {
  if (!serverForm.name || !serverForm.base_url) {
    ElMessage.warning('请填写名称和服务地址')
    return
  }

  try {
    // 在添加或更新服务器之前，先测试连接
    isTestingConnection.value = true
    await TestOllamaServer(serverForm.base_url)
    isTestingConnection.value = false
    
    if (serverForm.type === 'local') {
      // 保存本地配置
      await SaveOllamaServerConfig(serverForm.base_url)
      ElMessage.success('本地配置已保存')
      // 更新本地服务测试状态
      localServerState.test_status = 'success'
    } else {
      // 处理远程服务器
      if (editingServer.value) {
        // 更新现有服务器
        await UpdateRemoteServer({
          id: serverForm.id,
          name: serverForm.name,
          base_url: serverForm.base_url,
          api_key: serverForm.api_key,
          is_active: serverForm.is_active,
          test_status: 'success'
        })
        ElMessage.success('服务器已更新')
      } else {
        // 添加新服务器
        await AddRemoteServer({
          id: Date.now().toString(),
          name: serverForm.name,
          base_url: serverForm.base_url,
          api_key: serverForm.api_key,
          is_active: false,
          test_status: 'success'
        })
        ElMessage.success('服务器已添加')
      }
    }
    
    // 关闭抽屉
    drawerVisible.value = false
    // 刷新服务器列表
    loadRemoteServers()
  } catch (error) {
    isTestingConnection.value = false
    ElMessage.error((editingServer.value ? '更新' : '添加') + '失败: ' + (error as Error).message)
    
    // 更新测试状态为失败
    if (serverForm.type === 'local') {
      localServerState.test_status = 'failed'
    }
  }
}

// 删除远程服务器
const removeRemoteServer = (server: any) => {
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

// 设置活动服务器
const setActiveServer = async (server: any) => {
  try {
    // 如果是本地服务，直接保存本地配置
    if (server.type === 'local') {
      await SaveOllamaServerConfig(server.base_url)
      ElMessage.success('已设置本地服务为默认服务器')
      loadRemoteServers()
      return
    }
    
    await SetActiveServer(server.id)
    ElMessage.success(`已设置 ${server.name} 为默认服务器`)
    loadRemoteServers()
  } catch (error) {
    ElMessage.error('设置失败: ' + (error as Error).message)
  }
}

// 取消编辑
const cancelEdit = () => {
  drawerVisible.value = false
}

// 加载远程服务器列表
const loadRemoteServers = async () => {
  try {
    const servers = await GetRemoteServers()
    remoteServers.value = servers
    
    // 自动测试所有服务的连接
    setTimeout(async () => {
      for (const server of allServers.value) {
        await testConnection(server, true) // 传入true表示是自动测试
      }
    }, 100)
  } catch (error) {
    console.error('加载远程服务器列表失败:', error)
  }
}

onMounted(() => {
  loadRemoteServers()
})
</script>

<style scoped>
.ollama-settings {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.server-config-content {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.local-config-section,
.server-list-section {
  padding: 20px 0;
}

.local-config-section h3,
.server-list-section h3 {
  margin-top: 0;
  margin-bottom: 15px;
  color: #333;
  border-left: 4px solid #409eff;
  padding-left: 10px;
}

.local-config-section {
  border-bottom: 1px solid #ebeef5;
}

/* 表格操作按钮样式优化 */
.el-table .el-button {
  margin-right: 5px;
  margin-bottom: 5px;
}

.el-table .el-button:last-child {
  margin-right: 0;
}
</style>
