<template>
  <div class="ollama-settings">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-icon">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 20H21L19 16H14L12 20Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M12 20V4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M6 8L12 4L18 8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M6 16L12 20L18 16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <circle cx="12" cy="12" r="2" stroke="currentColor" stroke-width="2"/>
          </svg>
        </div>
        <div class="header-text">
          <h1>{{ t('ollamaSettings.title') }}</h1>
          <p>{{ t('ollamaSettings.description') }}</p>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="main-content">
      <!-- Control Panel -->
      <div class="control-panel">
        <div class="panel-title">
          <h3>{{ t('ollamaSettings.ollamaServers') }}</h3>
          <p>{{ t('ollamaSettings.manageServers') }}</p>
        </div>
        <button class="btn primary" @click="openServiceDrawer(null)">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
            <line x1="12" y1="5" x2="12" y2="19" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <line x1="5" y1="12" x2="19" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          {{ t('ollamaSettings.addNewServer') }}
        </button>
      </div>

      <!-- Servers Table -->
      <div class="servers-table-container">
        <div class="table-wrapper">
          <table class="servers-table">
            <thead>
              <tr>
                <th>{{ t('ollamaSettings.serverName') }}</th>
                <th>{{ t('ollamaSettings.serverAddress') }}</th>
                <th>{{ t('ollamaSettings.connectionStatus') }}</th>
                <th>{{ t('ollamaSettings.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="allServers.length === 0">
                <td colspan="4" class="empty-state">
                  <div class="empty-content">
                    <svg width="48" height="48" viewBox="0 0 24 24" fill="none">
                      <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
                      <path d="M19.4 15A1.65 1.65 0 0 0 21 13.09A1.65 1.65 0 0 0 19.4 9A1.65 1.65 0 0 0 21 6.91A1.65 1.65 0 0 0 19.4 3" stroke="currentColor" stroke-width="2"/>
                      <path d="M4.6 9A1.65 1.65 0 0 0 3 10.91A1.65 1.65 0 0 0 4.6 15A1.65 1.65 0 0 0 3 17.09A1.65 1.65 0 0 0 4.6 21" stroke="currentColor" stroke-width="2"/>
                    </svg>
                    <p>{{ t('ollamaSettings.noServersConfigured') }}</p>
                    <span>{{ t('ollamaSettings.addServerToStart') }}</span>
                  </div>
                </td>
              </tr>
              <tr v-for="server in allServers" :key="server.id" class="server-row">
                <td class="server-name">
                  <div class="server-info">
                    <span class="name">{{ server.name }}</span>
                    <span v-if="server.id === 'local'" class="local-badge">{{ t('ollamaSettings.local') }}</span>
                  </div>
                </td>
                <td>{{ server.baseUrl }}</td>
                <td>
                  <span class="status-badge" :class="{ 
                    success: server.testStatus === 'success', 
                    failed: server.testStatus === 'failed',
                    unknown: server.testStatus === 'unknown' || !server.testStatus
                  }">
                    <span class="status-dot"></span>
                    {{ server.testStatus === 'success' ? t('ollamaSettings.connected') : 
                       server.testStatus === 'failed' ? t('ollamaSettings.failed') : t('ollamaSettings.unknown') }}
                  </span>
                </td>
                <td>
                  <div class="action-buttons">
                    <button class="btn-small secondary" @click="testConnection(server)" :disabled="server.isTesting">
                      <span v-if="!server.isTesting">
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
                          <polyline points="22,12 18,12 15,21 9,3 6,12 2,12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                        </svg>
                        {{ t('ollamaSettings.test') }}
                      </span>
                      <div v-else class="loading-spinner"></div>
                    </button>
                    
                    <button class="btn-small secondary" @click="openServiceDrawer(server)">
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
                        <path d="M11 4H4A2 2 0 0 0 2 6V18A2 2 0 0 0 4 20H16A2 2 0 0 0 18 18V13" stroke="currentColor" stroke-width="2"/>
                        <path d="M18.5 2.5A2.12 2.12 0 0 1 21 5L12 14L8 15L9 11L18.5 2.5Z" stroke="currentColor" stroke-width="2"/>
                      </svg>
                      {{ t('ollamaSettings.edit') }}
                    </button>
                    
                    <button class="btn-small danger" @click="deleteService(server)" v-if="server.id !== 'local'">
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none">
                        <polyline points="3,6 5,6 21,6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                        <path d="M19 6V20A2 2 0 0 1 17 22H7A2 2 0 0 1 5 20V6M8 6V4A2 2 0 0 1 10 2H14A2 2 0 0 1 16 4V6" stroke="currentColor" stroke-width="2"/>
                      </svg>
                      {{ t('ollamaSettings.delete') }}
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Add/Edit Server Drawer -->
    <div v-if="serviceDrawerVisible" class="drawer-overlay">
      <div class="drawer-content" @click.stop>
        <div class="drawer-header">
          <h3>{{ serviceForm.id ? t('ollamaSettings.editServer') : t('ollamaSettings.addNewServerTitle') }}</h3>
          <button class="close-btn" @click="serviceDrawerVisible = false">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
              <line x1="18" y1="6" x2="6" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <line x1="6" y1="6" x2="18" y2="18" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </button>
        </div>
        
        <div class="drawer-body">
          <el-form ref="serviceFormRef" :model="serviceForm" :rules="formRules" label-position="top">
            <el-form-item :label="t('ollamaSettings.serverNameRequired')" prop="name">
              <el-input v-model="serviceForm.name" :placeholder="t('ollamaSettings.enterServerName')" />
            </el-form-item>
            
            <el-form-item :label="t('ollamaSettings.serverAddressRequired')" prop="baseUrl">
              <el-input v-model="serviceForm.baseUrl" :placeholder="t('ollamaSettings.serverAddressPlaceholder')" />
            </el-form-item>
            
            <el-form-item v-if="serviceForm.id !== 'local'" :label="t('ollamaSettings.apiKeyOptional')" prop="apiKey">
              <el-input v-model="serviceForm.apiKey" :placeholder="t('ollamaSettings.enterApiKey')" type="password" show-password />
            </el-form-item>
          </el-form>
        </div>
        
        <div class="drawer-footer">
          <button class="btn secondary" @click="serviceDrawerVisible = false">{{ t('ollamaSettings.cancel') }}</button>
          <button class="btn primary" @click="handleSaveService" :disabled="isSaving">
            <span v-if="!isSaving">{{ t('ollamaSettings.saveServer') }}</span>
            <div v-else class="loading-spinner"></div>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {onMounted, reactive, ref} from 'vue'
import {ElMessage, ElMessageBox, FormInstance, FormRules} from 'element-plus'
import { useI18n } from 'vue-i18n'
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

const { t } = useI18n();

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

const formRules = reactive<FormRules>({
  name: [
    { required: true, message: t('ollamaSettings.serverNameRequired'), trigger: 'blur' }
  ],
  baseUrl: [
    { required: true, message: t('ollamaSettings.serverAddressRequired'), trigger: 'blur' }
  ]
});

// --- Data Loading and Processing ---
const loadAllServers = async () => {
  try {
    const servers = await GetServers();
    allServers.value = servers.map(s => ({...s, isTesting: false}));
  } catch (error) {
    ElMessage.error(t('ollamaSettings.loadServerListFailed') + ': ' + (error as Error).message);
  }
};

onMounted(() => {
  loadAllServers();
});

// --- Drawer and Form Logic ---
const openServiceDrawer = (server: OllamaServerConfig | null) => {
  if (server) {
    // Edit
    Object.assign(serviceForm, server);
  } else {
    // New (default to remote)
    Object.assign(serviceForm, {
      id: '', // Empty ID indicates new remote service
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
  
  try {
    const valid = await serviceFormRef.value.validate();
    if (!valid) return;
    
    isSaving.value = true;
    const configToSave = {...serviceForm} as OllamaServerConfig;
    
    if (configToSave.id) {
      await UpdateServer(configToSave);
    } else {
      configToSave.id = Date.now().toString();
      await AddServer(configToSave);
    }
    
    ElMessage.success(t('ollamaSettings.serverConfigSaved'));
    serviceDrawerVisible.value = false;
    await loadAllServers();
  } catch (error) {
    ElMessage.error(t('ollamaSettings.saveFailed') + ': ' + (error as Error).message);
  } finally {
    isSaving.value = false;
  }
};

// --- Service Operations ---
const testConnection = async (server: OllamaServerConfig & { isTesting?: boolean }) => {
  server.isTesting = true;
  let newStatus: 'success' | 'failed' = 'failed';
  try {
    await TestOllamaServer(server.baseUrl);
    newStatus = 'success';
    ElMessage.success(`${server.name} ${t('ollamaSettings.connectionSuccess')}`);
  } catch (error) {
    newStatus = 'failed';
    let errorMessage = t('ollamaSettings.connectionFailed');
    
    if (error && typeof error === 'object' && 'message' in error) {
      const msg = (error as Error).message;
      if (msg && msg !== 'undefined' && msg.trim() !== '') {
        errorMessage = msg;
      } else {
        errorMessage = t('ollamaSettings.serverNotReachable');
      }
    } else {
      errorMessage = t('ollamaSettings.networkError');
    }
    
    ElMessage.error(`${server.name} ${errorMessage}`);
  } finally {
    server.isTesting = false;
    server.testStatus = newStatus;
    try {
      await UpdateServerTestStatus(server.id, newStatus)
    } catch (e) {
      ElMessage.error(t('ollamaSettings.updateStatusFailed') + ': ' + (e as Error).message)
      // Reload to maintain consistency if status update fails
      await loadAllServers();
    }
  }
};

const deleteService = (server: OllamaServerConfig) => {
  ElMessageBox.confirm(t('ollamaSettings.deleteConfirm', { name: server.name }), t('ollamaSettings.confirmDelete'), {
    confirmButtonText: t('ollamaSettings.confirm'),
    cancelButtonText: t('ollamaSettings.cancel'),
    type: 'warning'
  }).then(async () => {
    try {
      await DeleteServer(server.id);
      ElMessage.success(t('ollamaSettings.serverDeleted'));
      await loadAllServers();
    } catch (error) {
      ElMessage.error(t('ollamaSettings.deleteFailed') + ': ' + (error as Error).message);
    }
  }).catch(() => {
  });
};

</script>

<style scoped>
.ollama-settings {
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 2rem;
  box-sizing: border-box;
  overflow-y: auto;
}

.page-header {
  margin-bottom: 2rem;
  flex-shrink: 0;
  position: sticky;
  top: 0;
  z-index: 10;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.header-content {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-icon {
  width: 60px;
  height: 60px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  backdrop-filter: blur(10px);
}

.header-text h1 {
  margin: 0;
  font-size: 2rem;
  font-weight: 700;
  color: white;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-text p {
  margin: 0.5rem 0 0 0;
  color: rgba(255, 255, 255, 0.8);
  font-size: 1.1rem;
}

.main-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.control-panel {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
}

.panel-title h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #2d3748;
}

.panel-title p {
  margin: 0;
  color: #718096;
  font-size: 0.95rem;
}

.btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.25rem;
  border: none;
  border-radius: 8px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  text-decoration: none;
}

.btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.btn.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
}

.btn.secondary {
  background: white;
  color: #4a5568;
  border: 2px solid #e2e8f0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.btn.secondary:hover {
  background: #f7fafc;
  border-color: #cbd5e0;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.servers-table-container {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
}

.table-wrapper {
  overflow-x: auto;
}

.servers-table {
  width: 100%;
  border-collapse: collapse;
}

.servers-table th {
  background: #f8fafc;
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: #4a5568;
  border-bottom: 1px solid #e2e8f0;
  font-size: 0.9rem;
}

.servers-table td {
  padding: 1rem;
  border-bottom: 1px solid #f1f5f9;
  color: #2d3748;
}

.server-row:hover {
  background: #f8fafc;
}

.server-name {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.server-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.server-info .name {
  font-weight: 600;
  color: #2d3748;
}

.local-badge {
  background: rgba(72, 187, 120, 0.1);
  color: #38a169;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 500;
  border: 1px solid rgba(72, 187, 120, 0.2);
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.375rem 0.75rem;
  border-radius: 50px;
  font-size: 0.875rem;
  font-weight: 500;
}

.status-badge.success {
  background: rgba(72, 187, 120, 0.1);
  color: #38a169;
  border: 1px solid rgba(72, 187, 120, 0.2);
}

.status-badge.failed {
  background: rgba(245, 101, 101, 0.1);
  color: #e53e3e;
  border: 1px solid rgba(245, 101, 101, 0.2);
}

.status-badge.unknown {
  background: rgba(160, 174, 192, 0.1);
  color: #718096;
  border: 1px solid rgba(160, 174, 192, 0.2);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.btn-small {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.5rem 0.875rem;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-small.secondary {
  background: white;
  color: #4a5568;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.btn-small.secondary:hover:not(:disabled) {
  background: #f7fafc;
  border-color: #cbd5e0;
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
}

.btn-small.danger {
  background: rgba(245, 101, 101, 0.1);
  color: #e53e3e;
  border: 1px solid rgba(245, 101, 101, 0.2);
}

.btn-small.danger:hover {
  background: rgba(245, 101, 101, 0.2);
  transform: translateY(-1px);
}

.btn-small:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.empty-state {
  text-align: center;
  padding: 3rem;
}

.empty-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  color: #718096;
}

.empty-content svg {
  opacity: 0.5;
}

.empty-content p {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #4a5568;
}

.empty-content span {
  font-size: 0.95rem;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid currentColor;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Drawer Styles */
.drawer-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  justify-content: flex-end;
  z-index: 1000;
}

.drawer-content {
  background: white;
  width: 400px;
  max-width: 90vw;
  height: 100vh;
  box-shadow: -10px 0 30px rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
}

.drawer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
}

.drawer-header h3 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #2d3748;
}

.close-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 4px;
  color: #718096;
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: #f7fafc;
  color: #4a5568;
}

.drawer-body {
  flex: 1;
  padding: 1.5rem;
  overflow-y: auto;
}

.server-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-field label {
  font-weight: 500;
  color: #4a5568;
  font-size: 0.9rem;
}

.text-input {
  padding: 0.875rem 1rem;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 1rem;
  background: white;
  color: #2d3748;
  transition: all 0.3s ease;
}

.text-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1.5rem;
  border-top: 1px solid #e2e8f0;
  background: #f8fafc;
}

@media (max-width: 768px) {
  .ollama-settings {
    padding: 1rem;
  }
  
  .control-panel {
    flex-direction: column;
    align-items: stretch;
  }
  
  .servers-table {
    font-size: 0.875rem;
  }
  
  .servers-table th,
  .servers-table td {
    padding: 0.75rem 0.5rem;
  }
  
  .action-buttons {
    flex-direction: column;
  }
  
  .drawer-content {
    width: 100vw;
  }
}
</style>