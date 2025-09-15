<template>
  <div class="api-debugger-container">
    <!-- Left Sidebar: Server Selector and API List -->
    <div class="sidebar">
      <div class="server-selector-wrapper">
        <select v-model="selectedServerId" class="base-url-selector">
          <option value="">{{ t('apiDebugger.selectServer') }}</option>
          <option v-for="server in servers" :key="server.id" :value="server.id">
            {{ server.name }} ({{ server.baseUrl }})
          </option>
        </select>
      </div>
      <div class="api-list">
        <h3>{{ t('apiDebugger.ollamaApiList') }}</h3>
        <ul>
          <li
            v-for="api in apiDefinitions"
            :key="api.name"
            :class="{ active: activeApiName === api.name }"
            @click="selectApi(api)"
          >
            <span :class="['method-badge', api.method.toLowerCase()]">{{ api.method }}</span>
            {{ api.name }}
          </li>
        </ul>
      </div>
    </div>

    <!-- Main Content: API Debugger Component -->
    <div class="main-content">
      <ApiDebugger
        ref="apiDebuggerRef"
        :base-url="selectedServerUrl"
        :url-placeholder="'/api/tags'"
        :on-request="handleRequest"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { GetOllamaServers, SendHttpRequest } from '../../wailsjs/go/main/App';
import ApiDebugger from '../components/ApiDebugger.vue';

const { t } = useI18n();

const servers = ref([]);
const selectedServerId = ref('');
const activeApiName = ref('');
const apiDebuggerRef = ref();

// 选中的服务器 URL
const selectedServerUrl = computed(() => {
  const server = servers.value.find(s => s.id === selectedServerId.value);
  return server ? server.baseUrl : '';
});

// --- API Definitions based on API.md ---
const baseApiDefinitions = [
  {
    key: 'generateCompletion',
    method: 'POST',
    url: '/api/generate',
    body: {
      type: 'raw',
      rawContentType: 'application/json',
      rawContent: JSON.stringify({ model: 'llama3.2', keep_alive: 0 }, null, 2),
    },
    queryParams: [],
    headers: [{ key: 'Content-Type', value: 'application/json', enabled: true }],
  },
  {
    key: 'generateChatCompletion',
    method: 'POST',
    url: '/api/chat',
    body: {
      type: 'raw',
      rawContentType: 'application/json',
      rawContent: JSON.stringify({ model: 'llama3.2', messages: [{ role: 'user', content: 'Hello' }], keep_alive: 60 }, null, 2),
    },
    queryParams: [],
    headers: [{ key: 'Content-Type', value: 'application/json', enabled: true }],
  },
  {
    key: 'createModel',
    method: 'POST',
    url: '/api/create',
    body: {
      type: 'raw',
      rawContentType: 'application/json',
      rawContent: JSON.stringify({ model: 'bert-base-chinese', files: { 'config.json': 'a1b2c3d4e5f6', 'generation_config.json': 'b2c3d4e5f6g7', 'special_tokens_map.json': 'c3d4e5f6g7h8', 'tokenizer.json': 'd4e5f6g7h8i9', 'tokenizer_config.json': 'e5f6g7h8i9j0', 'model.safetensors': 'f6g7h8i9j0k1' } }, null, 2),
    },
    queryParams: [],
    headers: [{ key: 'Content-Type', value: 'application/json', enabled: true }],
  },
  {
    key: 'listLocalModels',
    method: 'GET',
    url: '/api/tags',
    body: { type: 'none', rawContent: '', rawContentType: 'application/json' },
    queryParams: [],
    headers: [],
  },
  {
    key: 'showModelDetails',
    method: 'POST',
    url: '/api/show',
    body: {
      type: 'raw',
      rawContentType: 'application/json',
      rawContent: JSON.stringify({ model: 'llama3.2' }, null, 2),
    },
    queryParams: [],
    headers: [{ key: 'Content-Type', value: 'application/json', enabled: true }],
  },
  {
    key: 'copyModel',
    method: 'POST',
    url: '/api/copy',
    body: {
      type: 'raw',
      rawContentType: 'application/json',
      rawContent: JSON.stringify({ source: 'llama3.2', destination: 'llama3-backup' }, null, 2),
    },
    queryParams: [],
    headers: [{ key: 'Content-Type', value: 'application/json', enabled: true }],
  },
  {
    key: 'deleteModel',
    method: 'DELETE',
    url: '/api/delete',
    body: {
      type: 'raw',
      rawContentType: 'application/json',
      rawContent: JSON.stringify({ model: 'llama3:13b' }, null, 2),
    },
    queryParams: [],
    headers: [{ key: 'Content-Type', value: 'application/json', enabled: true }],
  },
  {
    key: 'pullModel',
    method: 'POST',
    url: '/api/pull',
    body: {
      type: 'raw',
      rawContentType: 'application/json',
      rawContent: JSON.stringify({ model: 'llama3.2' }, null, 2),
    },
    queryParams: [],
    headers: [{ key: 'Content-Type', value: 'application/json', enabled: true }],
  },
  {
    key: 'generateEmbeddingsMultiple',
    method: 'POST',
    url: '/api/embed',
    body: {
      type: 'raw',
      rawContentType: 'application/json',
      rawContent: JSON.stringify({ model: 'gpt-4', input: ['How is the weather today?', 'What is the current temperature in Beijing?'] }, null, 2),
    },
    queryParams: [],
    headers: [{ key: 'Content-Type', value: 'application/json', enabled: true }],
  },
  {
    key: 'listRunningModels',
    method: 'GET',
    url: '/api/ps',
    body: { type: 'none', rawContent: '', rawContentType: 'application/json' },
    queryParams: [],
    headers: [],
  },
  {
    key: 'generateSingleEmbedding',
    method: 'POST',
    url: '/api/embeddings',
    body: {
      type: 'raw',
      rawContentType: 'application/json',
      rawContent: JSON.stringify({ model: 'all-minilm', prompt: 'This is an article about llamas...' }, null, 2),
    },
    queryParams: [],
    headers: [{ key: 'Content-Type', value: 'application/json', enabled: true }],
  },
];

const apiDefinitions = computed(() => {
  return baseApiDefinitions.map(api => ({
    ...api,
    name: t(`apiDebugger.apis.${api.key}`)
  }));
});
// --- End API Definitions ---

onMounted(() => {
  loadServers();
  // Select the first API by default
  if (apiDefinitions.value.length > 0) {
    selectApi(apiDefinitions.value[0]);
  }
});

async function loadServers() {
  try {
    servers.value = await GetOllamaServers();
    // Set a default active server if available
    const activeServer = servers.value.find(s => s.isActive);
    if (activeServer) {
      selectedServerId.value = activeServer.id;
    } else if (servers.value.length > 0) {
      // If no active server, select the first one
      selectedServerId.value = servers.value[0].id;
    }
  } catch (e) {
    console.error(`${t('apiDebugger.failedToLoadServers')}: ${e}`);
  }
}

function selectApi(api) {
  activeApiName.value = api.name;
  
  // 使用 ApiDebugger 组件的 setRequest 方法
  if (apiDebuggerRef.value) {
    apiDebuggerRef.value.setRequest({
      method: api.method,
      url: api.url,
      queryParams: api.queryParams ? JSON.parse(JSON.stringify(api.queryParams)) : [],
      headers: api.headers ? JSON.parse(JSON.stringify(api.headers)) : [],
      body: {
        type: api.body.type,
        rawContent: api.body.rawContent || '',
        rawContentType: api.body.rawContentType || 'application/json',
        formData: api.body.formData ? JSON.parse(JSON.stringify(api.body.formData)) : []
      }
    });
    
    // 清除之前的响应
    apiDebuggerRef.value.clearResponse();
  }
}

// 自定义请求处理器，使用 Wails 的 SendHttpRequest
async function handleRequest(request) {
  if (!selectedServerId.value) {
    throw new Error(t('apiDebugger.selectServer'));
  }

  // 构建后端请求对象
  const backendRequest = {
    method: request.method,
    selectedServerId: selectedServerId.value,
    path: request.url,
    queryParams: request.queryParams.filter(p => p.enabled).map(p => ({ ...p })),
    headers: request.headers.filter(h => h.enabled).map(h => ({ ...h })),
    body: {
      type: request.body.type,
      rawContent: request.body.rawContent,
      rawContentType: request.body.rawContentType,
      formData: request.body.formData.map(f => ({ ...f })),
    },
  };

  const result = await SendHttpRequest(backendRequest);
  if (result.error) {
    throw new Error(result.error);
  }
  
  return result;
}
</script>

<style scoped>
.api-debugger-container {
  display: flex;
  height: 100%;
  background-color: #f0f2f5;
}

.sidebar {
  width: 280px;
  flex-shrink: 0;
  background-color: #fff;
  border-right: 1px solid #e0e0e0;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.server-selector-wrapper {
  margin-bottom: 1rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #eee;
}

.base-url-selector {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  background-color: #f9f9f9;
}

.api-list {
  flex-grow: 1;
  overflow-y: auto;
  min-height: 0;
}

.api-list h3 {
  margin-top: 0;
  margin-bottom: 0.8rem;
  color: #333;
  font-size: 1.1rem;
}

.api-list ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.api-list li {
  padding: 0.6rem 0.5rem;
  cursor: pointer;
  border-radius: 4px;
  margin-bottom: 0.2rem;
  display: flex;
  align-items: center;
  font-size: 0.95rem;
  color: #555;
}

.api-list li:hover {
  background-color: #eef;
}

.api-list li.active {
  background-color: #e0e7ff;
  font-weight: bold;
  color: #333;
}

.method-badge {
  display: inline-block;
  padding: 0.2em 0.5em;
  border-radius: 3px;
  font-size: 0.75em;
  font-weight: bold;
  margin-right: 0.6rem;
  min-width: 45px;
  text-align: center;
  color: #fff;
}

.method-badge.get { background-color: #61affe; }
.method-badge.post { background-color: #49cc90; }
.method-badge.put { background-color: #fca130; }
.method-badge.delete { background-color: #f93e3e; }
.method-badge.patch { background-color: #50e3c2; }
.method-badge.head { background-color: #9012fe; }

.main-content {
  flex-grow: 1;
  padding: 1rem;
  overflow: hidden;
}
</style>
