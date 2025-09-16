<template>
  <div class="adapter-settings-page">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-icon">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L13.09 8.26L22 9L13.09 9.74L12 16L10.91 9.74L2 9L10.91 8.26L12 2Z" stroke="currentColor" stroke-width="2" stroke-linejoin="round"/>
            <path d="M12 12C14.21 12 16 10.21 16 8S14.21 4 12 4S8 5.79 8 8S9.79 12 12 12Z" stroke="currentColor" stroke-width="2"/>
            <path d="M8 21V19C8 16.79 9.79 15 12 15S16 16.79 16 19V21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <div class="header-text">
          <h1>{{ t('openaiAdapter.title') }}</h1>
          <p>{{ t('openaiAdapter.description') }}</p>
        </div>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="main-content">
      <!-- Top Row: Configuration and Tools -->
      <div class="top-row">
        <!-- Configuration Card -->
        <div class="config-card">
          <div class="card-header">
            <div>
              <h3>{{ t('openaiAdapter.serviceConfiguration') || 'Service Configuration' }}</h3>
              <p>{{ t('openaiAdapter.configDescription') || 'Configure OpenAI adapter settings' }}</p>
            </div>
            <div class="status-badge" :class="{ 'running': status.isRunning, 'stopped': !status.isRunning }">
              <span class="status-dot"></span>
              {{ status.isRunning ? t('common.running') || 'Running' : t('common.stopped') || 'Stopped' }}
            </div>
          </div>
          
          <div class="config-form">
            <div class="form-row">
              <div class="form-field">
                <label for="listen-ip">{{ t('openaiAdapter.listenAddress') || 'Listen Address' }}</label>
                <div class="input-wrapper">
                  <input id="listen-ip" type="text" v-model="store.config.listenIp" :disabled="status.isRunning" placeholder="0.0.0.0" />
                </div>
              </div>
              
              <div class="form-field">
                <label for="listen-port">{{ t('openaiAdapter.listenPort') || 'Listen Port' }}</label>
                <div class="input-wrapper">
                  <input id="listen-port" type="number" v-model.number="store.config.listenPort" :disabled="status.isRunning" placeholder="11434" />
                </div>
              </div>
            </div>
            
            <div class="form-field full-width">
              <label for="ollama-server">{{ t('openaiAdapter.targetService') || 'Target Service' }}</label>
              <div class="select-wrapper">
                <select id="ollama-server" v-model="store.config.targetOllamaServerId" :disabled="status.isRunning" @change="onServerChange">
                  <option value="" disabled>{{ t('openaiAdapter.selectService') || 'Select Service' }}</option>
                  <option v-for="server in ollamaServers" :key="server.id" :value="server.id">
                    {{ server.name }} ({{ server.baseUrl }})
                  </option>
                </select>
                <svg class="select-arrow" width="12" height="12" viewBox="0 0 24 24" fill="none">
                  <path d="M6 9L12 15L18 9" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </div>
            </div>
            
            <div class="form-field full-width" v-if="store.config.targetOllamaServerId">
              <label for="target-model">{{ t('openaiAdapter.targetModel') || 'Target Model' }}</label>
              <div class="select-wrapper">
                <select id="target-model" v-model="selectedModel" :disabled="status.isRunning || !availableModels.length">
                  <option value="" disabled>{{ t('openaiAdapter.selectModel') || 'Select Model' }}</option>
                  <option v-for="model in availableModels" :key="model.name" :value="model.name">
                    {{ model.name }}
                  </option>
                </select>
                <svg class="select-arrow" width="12" height="12" viewBox="0 0 24 24" fill="none">
                  <path d="M6 9L12 15L18 9" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </div>
            </div>
          </div>
          
          <!-- Control Buttons -->
          <div class="control-buttons">
            <button :class="['toggle-btn', { 'running': status.isRunning }]" @click="toggleService" :disabled="isToggling">
              <span v-if="!isToggling" class="btn-content">
                <svg v-if="!status.isRunning" width="16" height="16" viewBox="0 0 24 24" fill="none">
                  <polygon points="5,3 19,12 5,21" fill="currentColor"/>
                </svg>
                <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none">
                  <rect x="6" y="4" width="4" height="16" fill="currentColor"/>
                  <rect x="14" y="4" width="4" height="16" fill="currentColor"/>
                </svg>
                {{ status.isRunning ? t('openaiAdapter.stopService') || 'Stop Service' : t('openaiAdapter.startService') || 'Start Service' }}
              </span>
              <div v-else class="loading-spinner"></div>
            </button>
            
            <button class="action-btn primary" @click="store.saveConfig" :disabled="status.isRunning">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                <path d="M19 21H5A2 2 0 0 1 3 19V5A2 2 0 0 1 5 3H16L21 8V19A2 2 0 0 1 19 21Z" stroke="currentColor" stroke-width="2"/>
                <polyline points="17,21 17,13 7,13 7,21" stroke="currentColor" stroke-width="2"/>
                <polyline points="7,3 7,8 15,8" stroke="currentColor" stroke-width="2"/>
              </svg>
              {{ t('openaiAdapter.saveSettings') || 'Save Settings' }}
            </button>
          </div>
          
          <div v-if="status.error && !status.isRunning" class="error-message">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
              <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
              <line x1="15" y1="9" x2="9" y2="15" stroke="currentColor" stroke-width="2"/>
              <line x1="9" y1="9" x2="15" y2="15" stroke="currentColor" stroke-width="2"/>
            </svg>
            {{ status.error }}
          </div>
        </div>

        <!-- Tools Card -->
        <div class="tools-card">
          <div class="card-header">
            <div>
              <h3>{{ t('openaiAdapter.tools') || '工具' }}</h3>
              <p>{{ t('openaiAdapter.toolsDescription') || '调试和日志工具' }}</p>
            </div>
          </div>
          
          <div class="tools-buttons">
            <button class="tool-btn" @click="store.fetchApiDocs" :disabled="!status.isRunning">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                <path d="M9 12L11 14L15 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M21 12C21 16.97 16.97 21 12 21S3 16.97 3 12S7.03 3 12 3S21 7.03 21 12Z" stroke="currentColor" stroke-width="2"/>
              </svg>
              {{ t('openaiAdapter.apiDebugger') || 'API Debugger' }}
            </button>
            
            <button class="tool-btn" @click="store.toggleLogDrawer(true)">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                <path d="M14 2H6A2 2 0 0 0 4 4V20A2 2 0 0 0 6 22H18A2 2 0 0 0 20 20V8Z" stroke="currentColor" stroke-width="2"/>
                <polyline points="14,2 14,8 20,8" stroke="currentColor" stroke-width="2"/>
                <line x1="9" y1="15" x2="15" y2="15" stroke="currentColor" stroke-width="2"/>
              </svg>
              {{ t('openaiAdapter.viewLogs') || 'View Logs' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- API Documentation Drawer -->
    <el-drawer v-model="isApiDrawerVisible" :title="t('openaiAdapter.apiUsageExamples')" direction="rtl" size="50%" :close-on-click-modal="false">
      <div class="drawer-content">
        <div class="api-docs-tabs">
          <div class="tab-headers">
            <button 
              :class="{ active: activeApiTab === 'examples' }" 
              @click="activeApiTab = 'examples'"
            >
              {{ t('openaiAdapter.curlExamples') }}
            </button>
            <button 
              :class="{ active: activeApiTab === 'debugger' }" 
              @click="activeApiTab = 'debugger'"
            >
              {{ t('openaiAdapter.apiDebugger') }}
            </button>
          </div>
          
          <div class="tab-content">
            <!-- cURL 示例 -->
            <div v-if="activeApiTab === 'examples'" class="examples-content">
              <div v-for="(command, title) in apiDocs" :key="title" class="api-doc-item">
                <h3>{{ title }}</h3>
                <pre><code>{{ command }}</code></pre>
              </div>
            </div>
            
            <!-- API 调试器 -->
            <div v-if="activeApiTab === 'debugger'" class="debugger-content">
              <div class="debugger-header">
                <div class="header-left">
                  <h3>{{ t('openaiAdapter.apiDebuggerTitle') }}</h3>
                  <p>{{ t('openaiAdapter.apiDebuggerDescription') }}</p>
                </div>
                <div class="header-right">
                  <div class="stream-toggle">
                    <label class="toggle-label">
                      <input 
                        type="checkbox" 
                        v-model="isStreamMode" 
                        @change="onStreamModeChange"
                        class="toggle-checkbox"
                      />
                      <span class="toggle-slider"></span>
                      <span class="toggle-text">{{ t('openaiAdapter.streamingMode') }}</span>
                    </label>
                  </div>
                </div>
              </div>
              
              <ApiDebugger 
                ref="apiDebuggerRef"
                :base-url="adapterBaseUrl"
                url-placeholder="/v1/chat/completions"
                :on-request="handleApiRequest"
                :selected-model="selectedModel"
                default-url="/v1/chat/completions"
              />
            </div>
          </div>
        </div>
      </div>
    </el-drawer>

    <!-- Log Drawer -->
    <el-drawer v-model="isLogDrawerVisible" :title="t('openaiAdapter.realtimeLogs')"
               direction="rtl" size="50%" :close-on-click-modal="false">
      <div class="drawer-content log-drawer">
        <div class="log-actions">
          <button @click="store.clearLogs">{{ t('openaiAdapter.clearLogs') }}</button>
          <button @click="downloadLogs">{{ t('openaiAdapter.downloadLogs') }}</button>
        </div>
        <div class="log-container" ref="logContainerRef">
          <div v-for="(log, index) in logs" :key="index" :class="['log-line', `log-${log.level.toLowerCase()}`]">
            <span class="log-time">{{ new Date(log.timestamp).toLocaleTimeString() }}</span>
            <span class="log-level">[{{ log.level }}]</span>
            <span class="log-message">{{ log.message }}</span>
          </div>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, nextTick } from 'vue';
import { storeToRefs } from 'pinia';
import { useI18n } from 'vue-i18n';
import { useOpenAIAdapterStore } from '../stores/openaiAdapter';
import { EventsOn } from '../../wailsjs/runtime';
import { types } from '../../wailsjs/go/models';
import LogEntry = types.LogEntry;
import OpenAIAdapterStatus = types.OpenAIAdapterStatus;
import Model = types.Model;
import { debounce } from 'lodash-es';
import { ListModelsByServer } from '../../wailsjs/go/main/App';
import ApiDebugger from '../components/ApiDebugger.vue';

const { t } = useI18n();
const store = useOpenAIAdapterStore();
const { config, status, logs, ollamaServers, isLogDrawerVisible, isApiDrawerVisible, apiDocs } = storeToRefs(store);

const isToggling = ref(false);
const logContainerRef = ref<HTMLElement | null>(null);
const activeApiTab = ref('examples');
const apiDebuggerRef = ref();
const selectedModel = ref('');
const availableModels = ref<Model[]>([]);
const isStreamMode = ref(true);

// --- Computed Properties for UI --- 
const toggleButtonClass = computed(() => [
  'btn-toggle',
  status.value.isRunning ? 'btn-running' : 'btn-stopped',
]);

const toggleButtonText = computed(() => status.value.isRunning ? t('common.running') : t('openaiAdapter.startService'));

// 适配器基础 URL
const adapterBaseUrl = computed(() => {
  if (!status.value.isRunning) return '';
  const ip = config.value.listenIp === '0.0.0.0' ? '127.0.0.1' : config.value.listenIp;
  return `http://${ip}:${config.value.listenPort}`;
});

// 检查服务连接状态
const checkServiceConnection = async () => {
  if (!status.value.isRunning || !adapterBaseUrl.value) {
    return false;
  }
  
  try {
    const response = await fetch(`${adapterBaseUrl.value}/v1/models`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' }
    });
    return response.ok;
  } catch (error) {
    console.warn('Service connection check failed:', error);
    return false;
  }
};

// API 模板
const apiTemplates = computed(() => {
  const currentModel = selectedModel.value || 'llama3';
  
  return [
    {
      name: t('openaiAdapter.nonStreamingChat'),
      method: 'POST',
      url: '/v1/chat/completions',
      headers: [
        { key: 'Content-Type', value: 'application/json', enabled: true }
      ],
      body: {
        type: 'raw',
        rawContentType: 'application/json',
        rawContent: JSON.stringify({
          model: currentModel,
          messages: [
            {
              role: 'system',
              content: 'You are a helpful assistant.'
            },
            {
              role: 'user',
              content: 'Hello! Please introduce yourself.'
            }
          ],
          temperature: 0.7,
          max_tokens: 1000,
          stream: false
        }, null, 2)
      }
    },
    {
      name: t('openaiAdapter.streamingChat'),
      method: 'POST',
      url: '/v1/chat/completions',
      headers: [
        { key: 'Content-Type', value: 'application/json', enabled: true }
      ],
      body: {
        type: 'raw',
        rawContentType: 'application/json',
        rawContent: JSON.stringify({
          model: currentModel,
          messages: [
            {
              role: 'system',
              content: 'You are a helpful assistant that tells jokes.'
            },
            {
              role: 'user',
              content: 'Tell me a programming joke.'
            }
          ],
          temperature: 0.8,
          max_tokens: 500,
          stream: true
        }, null, 2)
      }
    }
  ];
});

// --- Service Control ---
const toggleService = debounce(async () => {
  if (isToggling.value) return;
  isToggling.value = true;

  try {
    if (status.value.isRunning) {
      await store.stopServer();
    } else {
      await store.startServer();
    }
    
    // Debug connection after service toggle
    if (process.env.NODE_ENV === 'development') {
      setTimeout(debugServiceConnection, 2000);
    }
  } finally {
    isToggling.value = false;
  }
}, 300);

// --- Model Management ---
const loadModelsForServer = async (serverId: string) => {
  if (!serverId) {
    availableModels.value = [];
    selectedModel.value = '';
    return;
  }
  
  try {
    availableModels.value = await ListModelsByServer(serverId);
    // Auto-select first model if available
    if (availableModels.value.length > 0 && !selectedModel.value) {
      selectedModel.value = availableModels.value[0].name;
    }
  } catch (error) {
    console.error('Failed to load models:', error);
    availableModels.value = [];
    selectedModel.value = '';
  }
};

const onServerChange = () => {
  selectedModel.value = '';
  if (store.config.targetOllamaServerId) {
    loadModelsForServer(store.config.targetOllamaServerId);
  }
};

// --- API Request Handling ---
const handleApiRequest = async (request) => {
  if (!status.value.isRunning) {
    throw new Error(t('openaiAdapter.serviceNotRunning'));
  }
  
  // 构建完整 URL
  let url = request.url;
  if (adapterBaseUrl.value && !url.startsWith('http')) {
    url = adapterBaseUrl.value.replace(/\/$/, '') + '/' + url.replace(/^\//, '');
  }
  
  // 验证 URL
  try {
    new URL(url);
  } catch (e) {
    throw new Error(`${t('apiDebugger.requestFailed')}: Invalid URL format`);
  }
  
  // 构建请求头
  const headers = {};
  request.headers.filter(h => h.enabled && h.key).forEach(h => {
    headers[h.key] = h.value;
  });
  
  // 构建请求体
  let body;
  if (request.body.type === 'raw' && request.body.rawContent) {
    body = request.body.rawContent;
  } else if (request.body.type === 'formData') {
    const formData = new FormData();
    request.body.formData.filter(f => f.key).forEach(f => {
      formData.append(f.key, f.value);
    });
    body = formData;
  }
  
  // 发送请求
  const startTime = Date.now();
  try {
    const fetchOptions = {
      method: request.method,
      headers
    };
    
    // 只有非 GET/HEAD 请求才添加 body
    if (!['GET', 'HEAD'].includes(request.method) && body) {
      fetchOptions.body = body;
    }
    
    const fetchResponse = await fetch(url, fetchOptions);
    const endTime = Date.now();
    
    const responseHeaders = Array.from(fetchResponse.headers.entries()).map(([key, value]) => ({
      key,
      value
    }));
    
    let responseBody = '';
    const contentType = fetchResponse.headers.get('content-type') || '';
    
    if (contentType.includes('text/event-stream')) {
      // Handle streaming response
      const reader = fetchResponse.body?.getReader();
      const decoder = new TextDecoder();
      
      if (reader) {
        try {
          while (true) {
            const { done, value } = await reader.read();
            if (done) break;
            responseBody += decoder.decode(value, { stream: true });
          }
        } finally {
          reader.releaseLock();
        }
      }
    } else {
      // Handle regular response
      responseBody = await fetchResponse.text();
    }
    
    return {
      statusCode: fetchResponse.status,
      statusText: fetchResponse.statusText,
      headers: responseHeaders,
      body: responseBody,
      requestDurationMs: endTime - startTime
    };
  } catch (error) {
    console.error('Request failed:', error);
    
    // 提供更详细的错误信息
    let errorMessage = 'Network error';
    if (error.name === 'TypeError' && error.message.includes('Failed to fetch')) {
      errorMessage = 'Connection failed - service may be down or unreachable';
    } else if (error.name === 'AbortError') {
      errorMessage = 'Request timeout';
    } else if (error.message) {
      errorMessage = error.message;
    }
    
    throw new Error(`${t('apiDebugger.requestFailed')}: ${errorMessage}`);
  }
};

// --- Stream Mode Handling ---
const onStreamModeChange = () => {
  if (apiDebuggerRef.value) {
    const currentModel = selectedModel.value || 'llama3';
    const template = isStreamMode.value ? apiTemplates.value[1] : apiTemplates.value[0];
    apiDebuggerRef.value.setRequest({
      method: template.method,
      url: template.url,
      queryParams: [],
      headers: template.headers || [],
      body: template.body || {
        type: 'none',
        rawContent: '',
        rawContentType: 'application/json',
        formData: []
      }
    });
  }
};

// --- Debug Functions ---
const debugServiceConnection = async () => {
  console.log('=== Service Connection Debug ===');
  console.log('Service running:', status.value.isRunning);
  console.log('Adapter base URL:', adapterBaseUrl.value);
  console.log('Config:', config.value);
  
  if (status.value.isRunning && adapterBaseUrl.value) {
    const isConnected = await checkServiceConnection();
    console.log('Connection test result:', isConnected);
    
    if (!isConnected) {
      console.warn('Service appears to be running but is not responding to requests');
    }
  }
};

// --- API Template Handling ---
const loadApiTemplate = (template) => {
  if (apiDebuggerRef.value) {
    apiDebuggerRef.value.setRequest({
      method: template.method,
      url: template.url,
      queryParams: [],
      headers: template.headers || [],
      body: template.body || {
        type: 'none',
        rawContent: '',
        rawContentType: 'application/json',
        formData: []
      }
    });
    apiDebuggerRef.value.clearResponse();
  }
};

// 暴露调试函数供开发者控制台使用
if (process.env.NODE_ENV === 'development') {
  (window as any).debugOpenAIAdapter = {
    checkConnection: debugServiceConnection,
    getConfig: () => config.value,
    getStatus: () => status.value,
    getBaseUrl: () => adapterBaseUrl.value
  };
}

// --- Log Handling ---
const downloadLogs = () => {
  const logContent = logs.value.map(log => 
    `[${new Date(log.timestamp).toISOString()}] [${log.level}] ${log.message}`
  ).join('\n');
  const blob = new Blob([logContent], { type: 'text/plain;charset=utf-8' });
  const link = document.createElement('a');
  link.href = URL.createObjectURL(blob);
  link.download = `openai_adapter_${new Date().toISOString().split('T')[0]}.log`;
  link.click();
  URL.revokeObjectURL(link.href);
};

watch(logs, () => {
  nextTick(() => {
    const container = logContainerRef.value;
    if (container) {
      container.scrollTop = container.scrollHeight;
    }
  });
}, { deep: true });

// --- Lifecycle Hooks ---
onMounted(() => {
  // Delay initialization to avoid hot reload issues
  setTimeout(async () => {
    try {
      await store.fetchConfig();
      await store.fetchOllamaServers();
      await store.fetchInitialStatus();
      
      // Load models if server is already selected
      if (store.config.targetOllamaServerId) {
        await loadModelsForServer(store.config.targetOllamaServerId);
      }
      
      // Debug service connection
      if (process.env.NODE_ENV === 'development') {
        await debugServiceConnection();
      }
    } catch (error) {
      console.error('Failed to initialize OpenAI Adapter:', error);
    }
  }, 100);

  // Listen for backend events
  EventsOn('openai-adapter-log', (log: LogEntry) => {
    store.addLog(log);
  });

  EventsOn('openai-adapter-status-changed', (newStatus: OpenAIAdapterStatus) => {
    store.setStatus(newStatus);
    // Debug connection when status changes
    if (process.env.NODE_ENV === 'development') {
      setTimeout(debugServiceConnection, 1000);
    }
  });
});

</script>

<style scoped>
.adapter-settings-page {
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 2rem;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  overflow: hidden;
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
  flex: 1;
  overflow-y: auto;
  min-height: 0;
}

.top-row {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 1.5rem;
}

.control-buttons {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 1px solid #e2e8f0;
  justify-content: center;
}

.control-buttons .toggle-btn,
.control-buttons .action-btn {
  flex: 1;
  max-width: 200px;
}

.tools-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.tools-buttons {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.tool-btn {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 0.75rem;
  padding: 1rem 1.5rem;
  border: none;
  border-radius: 8px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  background: white;
  color: #4a5568;
  border: 2px solid #e2e8f0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  width: 100%;
  text-align: left;
}

.tool-btn:hover {
  background: #f7fafc;
  border-color: #cbd5e0;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.tool-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.status-card, .config-card, .actions-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
}

.card-header h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #2d3748;
}

.card-header p {
  margin: 0;
  color: #718096;
  font-size: 0.95rem;
}

.status-badge {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  border-radius: 50px;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.3s ease;
}

.status-badge.running {
  background: rgba(72, 187, 120, 0.1);
  color: #38a169;
  border: 1px solid rgba(72, 187, 120, 0.2);
}

.status-badge.stopped {
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

.status-body {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.toggle-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 1rem 2rem;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.toggle-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
}

.toggle-btn.running {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  box-shadow: 0 4px 15px rgba(245, 87, 108, 0.4);
}

.toggle-btn.running:hover {
  box-shadow: 0 6px 20px rgba(245, 87, 108, 0.6);
}

.toggle-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-content {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem;
  background: rgba(245, 101, 101, 0.1);
  border: 1px solid rgba(245, 101, 101, 0.2);
  border-radius: 8px;
  color: #e53e3e;
  font-size: 0.9rem;
}

.card-header {
  margin-bottom: 2rem;
}

.card-header h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #2d3748;
}

.card-header p {
  margin: 0;
  color: #718096;
  font-size: 0.95rem;
}

.config-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-field.full-width {
  grid-column: 1 / -1;
}

.form-field label {
  font-weight: 500;
  color: #4a5568;
  font-size: 0.9rem;
}

.input-wrapper, .select-wrapper {
  position: relative;
}

.input-wrapper input, .select-wrapper select {
  width: 100%;
  padding: 0.875rem 1rem;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 1rem;
  background: white;
  color: #2d3748;
  transition: all 0.3s ease;
  box-sizing: border-box;
}

.input-wrapper input:focus, .select-wrapper select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.input-wrapper input:disabled, .select-wrapper select:disabled {
  background: #f7fafc;
  color: #a0aec0;
  cursor: not-allowed;
}

.select-wrapper {
  position: relative;
}

.select-wrapper select {
  appearance: none;
  padding-right: 2.5rem;
}

.select-arrow {
  position: absolute;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
  color: #a0aec0;
  pointer-events: none;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 1rem 1.5rem;
  border: none;
  border-radius: 8px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.action-btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.action-btn.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
}

.action-btn.primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.action-btn.secondary {
  background: white;
  color: #4a5568;
  border: 2px solid #e2e8f0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.action-btn.secondary:hover {
  background: #f7fafc;
  border-color: #cbd5e0;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

@media (max-width: 768px) {
  .adapter-settings-page {
    padding: 1rem;
  }
  
  .top-row {
    grid-template-columns: 1fr;
  }
  
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .control-buttons {
    flex-direction: column;
  }
  
  .header-content {
    flex-direction: column;
    text-align: center;
  }
}

.drawer-content { 
  padding: 1.5rem; 
  background: #f8fafc;
  height: 100%;
  width: 100%;
  box-sizing: border-box;
}

.log-drawer { 
  display: flex; 
  flex-direction: column; 
  height: 100%; 
}

.log-actions { 
  flex-shrink: 0; 
  margin-bottom: 1rem; 
  display: flex; 
  gap: 1rem; 
}

.log-actions button { 
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white; 
  padding: 0.75rem 1.5rem; 
  border: none; 
  border-radius: 8px; 
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
}

.log-actions button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.log-container {
  flex-grow: 1;
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  overflow-y: auto;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  font-size: 0.9rem;
  border: 1px solid #e2e8f0;
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.05);
}

.log-line { 
  display: flex; 
  gap: 1rem; 
  padding: 0.25rem 0;
  border-bottom: 1px solid #f1f5f9;
}

.log-time { 
  color: #64748b;
  font-weight: 500;
}

.log-level {
  font-weight: 600;
  min-width: 60px;
}

.log-info .log-level { color: #0ea5e9; }
.log-warn .log-level { color: #f59e0b; }
.log-error .log-level { color: #ef4444; }

.log-message { 
  white-space: pre-wrap; 
  word-break: break-all;
  color: #374151;
}

.api-doc-item h3 { 
  border-bottom: 2px solid #e2e8f0; 
  padding-bottom: 0.75rem; 
  margin: 2rem 0 1rem 0;
  color: #1f2937;
  font-weight: 600;
}

.api-doc-item pre { 
  background: #f8fafc; 
  padding: 1.5rem; 
  border-radius: 12px; 
  white-space: pre-wrap; 
  word-break: break-all;
  border: 1px solid #e2e8f0;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
  color: #374151;
}

.api-docs-tabs {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.api-docs-tabs .tab-headers {
  display: flex;
  border-bottom: 2px solid #e2e8f0;
  margin-bottom: 1.5rem;
  flex-shrink: 0;
}

.api-docs-tabs .tab-headers button {
  padding: 1rem 1.5rem;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 500;
  color: #6b7280;
  border-bottom: 2px solid transparent;
  transition: all 0.3s ease;
}

.api-docs-tabs .tab-headers button:hover {
  color: #374151;
  background: #f9fafb;
}

.api-docs-tabs .tab-headers button.active {
  color: #667eea;
  border-bottom-color: #667eea;
  font-weight: 600;
}

.api-docs-tabs .tab-content {
  flex-grow: 1;
  overflow-y: auto;
}

.examples-content {
  padding: 0;
}

.debugger-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.debugger-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e2e8f0;
}

.header-left h3 {
  margin: 0 0 0.5rem 0;
  color: #1f2937;
  font-size: 1.25rem;
  font-weight: 600;
}

.header-left p {
  margin: 0;
  color: #6b7280;
  font-size: 0.95rem;
}

.header-right {
  flex-shrink: 0;
  margin-left: 2rem;
}

.stream-toggle {
  display: flex;
  align-items: center;
}

.toggle-label {
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 0.5rem;
}

.toggle-checkbox {
  display: none;
}

.toggle-slider {
  position: relative;
  width: 44px;
  height: 24px;
  background-color: #cbd5e0;
  border-radius: 12px;
  transition: background-color 0.3s;
}

.toggle-slider::before {
  content: '';
  position: absolute;
  top: 2px;
  left: 2px;
  width: 20px;
  height: 20px;
  background-color: white;
  border-radius: 50%;
  transition: transform 0.3s;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.toggle-checkbox:checked + .toggle-slider {
  background-color: #667eea;
}

.toggle-checkbox:checked + .toggle-slider::before {
  transform: translateX(20px);
}

.toggle-text {
  font-size: 0.9rem;
  color: #374151;
  font-weight: 500;
}

.template-buttons {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.template-btn {
  padding: 0.5rem 1rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(102, 126, 234, 0.2);
}

.template-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(102, 126, 234, 0.3);
}
</style>