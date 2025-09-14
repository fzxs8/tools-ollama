<template>
  <div class="adapter-settings-page">
    <div class="setting-card">
      <div class="card-header">
        <h1>OpenAI API 适配器</h1>
        <p>将本地 Ollama 服务模拟为 OpenAI API，兼容各类生态工具。</p>
      </div>

      <div class="card-body">
        <!-- 服务控制 -->
        <div class="control-section">
          <button :class="toggleButtonClass" @click="toggleService" :disabled="isToggling">
            <span v-if="!isToggling" class="button-text">{{ toggleButtonText }}</span>
            <div v-else class="spinner"></div>
          </button>
          <div class="status-indicator">
            <span :class="['dot', status.isRunning ? 'dot-running' : 'dot-stopped']"></span>
            <span>{{ status.isRunning ? '服务运行中' : '服务已停止' }}</span>
            <span v-if="status.error && !status.isRunning" class="error-text"> - {{ status.error }}</span>
          </div>
        </div>

        <!-- 配置区域 -->
        <div class="config-grid">
          <div class="form-group">
            <label for="listen-ip">监听地址</label>
            <input id="listen-ip" type="text" v-model="config.listen_ip" :disabled="status.isRunning" />
          </div>
          <div class="form-group">
            <label for="listen-port">监听端口</label>
            <input id="listen-port" type="number" v-model.number="config.listen_port" :disabled="status.isRunning" />
          </div>
          <div class="form-group">
            <label for="ollama-server">目标 Ollama 服务</label>
            <select id="ollama-server" v-model="config.target_ollama_server_id" :disabled="status.isRunning">
              <option value="" disabled>请选择一个服务</option>
              <option v-for="server in ollamaServers" :key="server.id" :value="server.id">
                {{ server.name }} ({{ server.baseUrl }})
              </option>
            </select>
          </div>
        </div>
      </div>

      <div class="card-footer">
        <button class="btn-secondary" @click="store.fetchApiDocs">查看 API</button>
        <button class="btn-secondary" @click="store.toggleLogDrawer(true)">查看日志</button>
        <button class="btn-primary" @click="store.saveConfig" :disabled="status.isRunning">保存设置</button>
      </div>
    </div>

    <!-- API 文档抽屉 -->
    <el-drawer v-model="isApiDrawerVisible" title="API 调用示例" direction="rtl" size="40%">
      <div class="drawer-content">
        <div v-for="(command, title) in apiDocs" :key="title" class="api-doc-item">
          <h3>{{ title }}</h3>
          <pre><code>{{ command }}</code></pre>
        </div>
      </div>
    </el-drawer>

    <!-- 日志抽屉 -->
    <el-drawer v-model="isLogDrawerVisible" title="适配器实时日志" direction="rtl" size="50%">
      <div class="drawer-content log-drawer">
        <div class="log-actions">
          <button @click="store.clearLogs">清空日志</button>
          <button @click="downloadLogs">下载日志</button>
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
import { useOpenAIAdapterStore } from '../stores/openaiAdapter';
import { EventsOn } from '../../wailsjs/runtime';
import { types } from '../../wailsjs/go/models';
import LogEntry = types.LogEntry;
import OpenAIAdapterStatus = types.OpenAIAdapterStatus;
import { debounce } from 'lodash-es';

const store = useOpenAIAdapterStore();
const { config, status, logs, ollamaServers, isLogDrawerVisible, isApiDrawerVisible, apiDocs } = storeToRefs(store);

const isToggling = ref(false);
const logContainerRef = ref<HTMLElement | null>(null);

// --- Computed Properties for UI --- 
const toggleButtonClass = computed(() => [
  'btn-toggle',
  status.value.isRunning ? 'btn-running' : 'btn-stopped',
]);

const toggleButtonText = computed(() => status.value.isRunning ? '服务运行中' : '启动服务');

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
  } finally {
    isToggling.value = false;
  }
}, 300);

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
  // Fetch initial data
  store.fetchConfig();
  store.fetchOllamaServers();
  store.fetchInitialStatus();

  // Listen for backend events
  EventsOn('openai-adapter-log', (log: LogEntry) => {
    store.addLog(log);
  });

  EventsOn('openai-adapter-status-changed', (newStatus: OpenAIAdapterStatus) => {
    store.setStatus(newStatus);
  });
});

</script>

<style scoped>
:root {
  --card-bg: #2c2c2e;
  --text-primary: #f2f2f7;
  --text-secondary: #aeb0b4;
  --border-color: #48484a;
  --input-bg: #3a3a3c;
  --btn-primary-bg: #0a84ff;
  --btn-secondary-bg: #505052;
  --dot-running: #34c759;
  --dot-stopped: #8e8e93;
  --log-info: #5ac8fa;
  --log-warn: #ffcc00;
  --log-error: #ff453a;
}

.adapter-settings-page {
  padding: 2rem;
  background-color: #1c1c1e;
  color: var(--text-primary);
  height: 100%;
  box-sizing: border-box;
}

.setting-card {
  max-width: 700px;
  margin: 2rem auto;
  background-color: var(--card-bg);
  border-radius: 12px;
  border: 1px solid var(--border-color);
}

.card-header {
  padding: 1.5rem;
  border-bottom: 1px solid var(--border-color);
}
.card-header h1 { margin: 0 0 0.5rem 0; font-size: 1.75rem; }
.card-header p { margin: 0; color: var(--text-secondary); }

.card-body { padding: 1.5rem; }

.control-section {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 2rem;
}

.btn-toggle {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s, transform 0.1s;
  min-width: 140px;
  display: flex;
  justify-content: center;
  align-items: center;
}
.btn-toggle:disabled { cursor: not-allowed; opacity: 0.7; }
.btn-stopped { background-color: var(--btn-primary-bg); color: white; }
.btn-running { background-color: var(--dot-running); color: white; }

.status-indicator {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--text-secondary);
}
.dot { width: 10px; height: 10px; border-radius: 50%; }
.dot-running { background-color: var(--dot-running); }
.dot-stopped { background-color: var(--dot-stopped); }
.error-text { color: var(--log-error); font-size: 0.9rem; }

.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: var(--text-secondary);
}
.form-group input, .form-group select {
  width: 100%;
  padding: 0.75rem;
  background-color: var(--input-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 1rem;
}
.form-group input:disabled, .form-group select:disabled { opacity: 0.6; cursor: not-allowed; }

.card-footer {
  padding: 1.5rem;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
}

.btn-primary, .btn-secondary {
  padding: 0.6rem 1.2rem;
  border: none;
  border-radius: 8px;
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
}
.btn-primary { background-color: var(--btn-primary-bg); color: white; }
.btn-secondary { background-color: var(--btn-secondary-bg); color: white; }
.btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }

.drawer-content { padding: 1rem; }
.log-drawer { display: flex; flex-direction: column; height: 100%; }
.log-actions { flex-shrink: 0; margin-bottom: 1rem; display: flex; gap: 1rem; }
.log-actions button { background-color: var(--btn-secondary-bg); color: white; padding: 0.5rem 1rem; border: none; border-radius: 6px; cursor: pointer; }

.log-container {
  flex-grow: 1;
  background-color: #1c1c1e;
  padding: 1rem;
  border-radius: 8px;
  overflow-y: auto;
  font-family: monospace;
  font-size: 0.9rem;
}

.log-line { display: flex; gap: 1rem; }
.log-time { color: var(--text-secondary); }
.log-level.log-info { color: var(--log-info); }
.log-level.log-warn { color: var(--log-warn); }
.log-level.log-error { color: var(--log-error); }
.log-message { white-space: pre-wrap; word-break: break-all; }

.api-doc-item h3 { border-bottom: 1px solid var(--border-color); padding-bottom: 0.5rem; margin-top: 2rem; }
.api-doc-item pre { background-color: #1c1c1e; padding: 1rem; border-radius: 8px; white-space: pre-wrap; word-break: break-all; }

.spinner { border: 3px solid rgba(255,255,255,0.3); border-top: 3px solid white; border-radius: 50%; width: 20px; height: 20px; animation: spin 1s linear infinite; }
@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
</style>