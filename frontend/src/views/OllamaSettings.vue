<template>
  <div class="ollama-settings">
    <el-card class="settings-card">
      <template #header>
        <div class="card-header">
          <span>Ollama 服务配置</span>
          <el-button type="primary" @click="openServiceDrawer(null)">添加新服务</el-button>
        </div>
      </template>

      <el-table :data="allServers" style="width: 100%" empty-text="暂无服务配置">
        <el-table-column prop="name" label="服务名称" width="180"/>
        <el-table-column prop="baseUrl" label="服务地址"/>
        <el-table-column prop="testStatus" label="连接状态" width="120">
          <template #default="scope">
            <el-tag v-if="scope.row.testStatus === 'success'" type="success">连接成功</el-tag>
            <el-tag v-else-if="scope.row.testStatus === 'failed'" type="danger">连接失败</el-tag>
            <el-tag v-else type="info">未知</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button size="small" @click="testConnection(scope.row)" :loading="scope.row.isTesting">测试</el-button>
            <el-button size="small" @click="openServiceDrawer(scope.row)">编辑</el-button>
            <el-button size="small" type="danger" @click="deleteService(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑服务抽屉 -->
    <el-drawer v-model="serviceDrawerVisible" :title="serviceForm.id ? '编辑服务' : '添加新服务'" direction="rtl"
               size="40%">
      <div class="drawer-content">
        <el-form :model="serviceForm" label-width="100px" ref="serviceFormRef">
          <el-form-item label="服务名称" prop="name"
                        :rules="[{ required: true, message: '服务名称不能为空', trigger: 'blur' }]">
            <el-input v-model="serviceForm.name" placeholder="请输入服务名称"></el-input>
          </el-form-item>
          <el-form-item label="服务地址" prop="baseUrl"
                        :rules="[{ required: true, message: '服务地址不能为空', trigger: 'blur' }]">
            <el-input v-model="serviceForm.baseUrl" placeholder="例如: http://localhost:11434"></el-input>
          </el-form-item>
          <el-form-item v-if="serviceForm.id !== 'local'" label="API Key">
            <el-input v-model="serviceForm.apiKey" placeholder="请输入 API Key（可选）" type="password"
                      show-password></el-input>
          </el-form-item>
        </el-form>
        <div class="drawer-footer">
          <el-button @click="serviceDrawerVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSaveService" :loading="isSaving">保存</el-button>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import {onMounted, reactive, ref} from 'vue'
import {ElMessage, ElMessageBox, FormInstance} from 'element-plus'
import {
  AddServer,
  DeleteServer,
  GetServers,
  TestOllamaServer,
  UpdateServer,
  UpdateServerTestStatus,
} from '../../wailsjs/go/main/App'
import {types} from "../../wailsjs/go/models";
import OllamaServerConfig = types.OllamaServerConfig;

// 响应式数据
const allServers = ref<(OllamaServerConfig & { isTesting?: boolean })[]>([])
const serviceDrawerVisible = ref(false)
const isSaving = ref(false)
const serviceFormRef = ref<FormInstance>()

const serviceForm = reactive<Partial<OllamaServerConfig>>({
  id: '',
  name: '',
  baseUrl: '',
  apiKey: '',
  testStatus: 'unknown'
});

// --- 数据加载与处理 ---
const loadAllServers = async () => {
  try {
    const servers = await GetServers();
    allServers.value = servers.map(s => ({...s, isTesting: false}));
  } catch (error) {
    ElMessage.error('加载服务列表失败: ' + (error as Error).message);
  }
};

onMounted(() => {
  loadAllServers();
});

// --- 抽屉与表单逻辑 ---
const openServiceDrawer = (server: OllamaServerConfig | null) => {
  if (server) {
    // 编辑
    Object.assign(serviceForm, server);
  } else {
    // 新建 (默认为远程)
    Object.assign(serviceForm, {
      id: '', // ID为空表示是新建远程服务
      name: '',
      baseUrl: 'http://localhost:11434',
      apiKey: '',
      testStatus: 'unknown'
    });
  }
  serviceDrawerVisible.value = true;
};

const handleSaveService = async () => {
  if (!serviceFormRef.value) return;
  await serviceFormRef.value.validate(async (valid) => {
    if (valid) {
      isSaving.value = true;
      try {
        const configToSave = {...serviceForm} as OllamaServerConfig;
        if (configToSave.id) {
          await UpdateServer(configToSave);
        } else {
          configToSave.id = Date.now().toString();
          await AddServer(configToSave);
        }
        ElMessage.success('服务配置已保存');
        serviceDrawerVisible.value = false;
        await loadAllServers();
      } catch (error) {
        ElMessage.error('保存失败: ' + (error as Error).message);
      } finally {
        isSaving.value = false;
      }
    }
  });
};

// --- 服务操作 ---
const testConnection = async (server: OllamaServerConfig & { isTesting?: boolean }) => {
  server.isTesting = true;
  let newStatus: 'success' | 'failed' = 'failed';
  try {
    await TestOllamaServer(server.baseUrl);
    newStatus = 'success';
    ElMessage.success(`${server.name} 连接成功`);
  } catch (error) {
    newStatus = 'failed';
    ElMessage.error(`${server.name} 连接失败: ${(error as Error).message}`);
  } finally {
    server.isTesting = false;
    server.testStatus = newStatus;
    try {
      await UpdateServerTestStatus(server.id, newStatus)
    } catch (e) {
      ElMessage.error("更新状态失败: " + (e as Error).message)
      // 如果状态更新失败，重新加载以保持一致性
      await loadAllServers();
    }
  }
};

const deleteService = (server: OllamaServerConfig) => {
  ElMessageBox.confirm(`确定要删除 “${server.name}” 吗？`, '确认删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await DeleteServer(server.id);
      ElMessage.success('服务已删除');
      await loadAllServers();
    } catch (error) {
      ElMessage.error('删除失败: ' + (error as Error).message);
    }
  }).catch(() => {
  });
};

</script>

<style scoped>
.ollama-settings {
  padding: 20px;
  height: 100%;
  box-sizing: border-box;
}

.settings-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.settings-card :deep(.el-card__body) {
  flex: 1;
  overflow-y: auto;
}

.drawer-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.drawer-content .el-form {
  flex: 1;
  padding: 0 20px;
}

.drawer-footer {
  text-align: right;
  padding: 20px;
  border-top: 1px solid #e0e0e0;
}
</style>
