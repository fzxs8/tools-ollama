<template>
  <div class="api-debugger-container">
    <!-- Left Sidebar: Server Selector and API List -->
    <div class="sidebar">
      <div class="server-selector-wrapper">
        <select v-model="request.selectedServerId" class="base-url-selector">
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

    <!-- Main Content: Request Panel and Response Panel -->
    <div class="main-content">
      <!-- Request Panel -->
      <div class="request-panel">
        <div class="request-url-bar">
          <span :class="['http-method-display', request.method.toLowerCase()]">{{ request.method }}</span>
          <input type="text" v-model="request.path" placeholder="/api/tags" class="url-input" />
          <button @click="sendRequest" :disabled="isLoading">
            {{ isLoading ? t('apiDebugger.sending') : t('apiDebugger.send') }}
          </button>
        </div>

        <!-- Request Tabs -->
        <div class="request-tabs">
          <div class="tab-headers">
            <button :class="{ active: activeTab === 'params' }" @click="activeTab = 'params'">{{ t('apiDebugger.queryParams') }}</button>
            <button :class="{ active: activeTab === 'headers' }" @click="activeTab = 'headers'">{{ t('apiDebugger.headers') }}</button>
            <button :class="{ active: activeTab === 'body' }" @click="activeTab = 'body'">{{ t('apiDebugger.body') }}</button>
          </div>
          <div class="tab-content-request">
            <!-- Query Params -->
            <div v-if="activeTab === 'params'" class="key-value-editor">
              <div v-for="(param, index) in request.queryParams" :key="index" class="kv-row">
                <input type="checkbox" v-model="param.enabled" />
                <input type="text" v-model="param.key" :placeholder="t('apiDebugger.key')" />
                <input type="text" v-model="param.value" :placeholder="t('apiDebugger.value')" />
                <button @click="removeQueryParam(index)">{{ t('apiDebugger.remove') }}</button>
              </div>
              <button @click="addQueryParam">{{ t('apiDebugger.addParameter') }}</button>
            </div>

            <!-- Headers -->
            <div v-if="activeTab === 'headers'" class="key-value-editor">
              <div v-for="(header, index) in request.headers" :key="index" class="kv-row">
                <input type="checkbox" v-model="header.enabled" />
                <input type="text" v-model="header.key" :placeholder="t('apiDebugger.key')" />
                <input type="text" v-model="header.value" :placeholder="t('apiDebugger.value')" />
                <button @click="removeHeader(index)">{{ t('apiDebugger.remove') }}</button>
              </div>
              <button @click="addHeader">{{ t('apiDebugger.addHeader') }}</button>
            </div>

            <!-- Body -->
            <div v-if="activeTab === 'body'">
              <div class="body-type-selector">
                <label><input type="radio" v-model="request.body.type" value="none" /> {{ t('apiDebugger.none') }}</label>
                <label><input type="radio" v-model="request.body.type" value="raw" /> {{ t('apiDebugger.raw') }}</label>
                <label><input type="radio" v-model="request.body.type" value="formData" /> {{ t('apiDebugger.formData') }}</label>
              </div>
              <div v-if="request.body.type === 'raw'" class="raw-body-editor">
                <select v-model="request.body.rawContentType">
                  <option value="application/json">JSON</option>
                  <option value="text/plain">Text</option>
                  <option value="text/html">HTML</option>
                  <option value="application/xml">XML</option>
                </select>
                <textarea v-model="request.body.rawContent" :placeholder="t('apiDebugger.enterRawContent')"></textarea>
              </div>
              <div v-if="request.body.type === 'formData'" class="key-value-editor">
                <div v-for="(item, index) in request.body.formData" :key="index" class="kv-row">
                  <input type="text" v-model="item.key" :placeholder="t('apiDebugger.key')" />
                  <input type="text" v-model="item.value" :placeholder="t('apiDebugger.value')" />
                  <button @click="removeFormData(index)">{{ t('apiDebugger.remove') }}</button>
                </div>
                <button @click="addFormData">{{ t('apiDebugger.addFormItem') }}</button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Response Panel -->
      <div class="response-panel">
        <div v-if="isLoading" class="response-placeholder">{{ t('apiDebugger.loading') }}</div>
        <div v-else-if="error" class="response-placeholder error-message">
          <h3>{{ t('apiDebugger.error') }}</h3>
          <pre>{{ error }}</pre>
        </div>
        <div v-else-if="response" class="response-wrapper">
          <div class="response-status-bar">
            <span>{{ t('apiDebugger.status') }}: <strong :class="statusClass">{{ response.statusCode }} {{ response.statusText }}</strong></span>
            <span>{{ t('apiDebugger.duration') }}: <strong>{{ response.requestDurationMs }} ms</strong></span>
          </div>
          <div class="response-tabs">
            <div class="tab-headers">
              <button :class="{ active: activeResponseTab === 'body' }" @click="activeResponseTab = 'body'">{{ t('apiDebugger.responseBody') }}</button>
              <button :class="{ active: activeResponseTab === 'headers' }" @click="activeResponseTab = 'headers'">{{ t('apiDebugger.responseHeaders') }}</button>
            </div>
            <div class="tab-content-response">
              <div v-if="activeResponseTab === 'body'">
                <pre class="response-body">{{ formattedBody }}</pre>
              </div>
              <div v-if="activeResponseTab === 'headers'" class="key-value-editor">
                <div v-for="(header, index) in response.headers" :key="index" class="kv-row-readonly">
                  <input type="text" :value="header.key" readonly :placeholder="t('apiDebugger.key')" />
                  <input type="text" :value="header.value" readonly :placeholder="t('apiDebugger.value')" />
                </div>
              </div>
            </div>
          </div>
        </div>
         <div v-else class="response-placeholder">
          <p>{{ t('apiDebugger.clickSendToRequest') }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { GetOllamaServers, SendHttpRequest } from '../../wailsjs/go/main/App';

const { t } = useI18n();

const servers = ref([]);
const activeTab = ref('params');
const activeResponseTab = ref('body');
const isLoading = ref(false);
const error = ref(null);
const activeApiName = ref(''); // To highlight the active API in the list

const request = reactive({
  method: 'GET',
  selectedServerId: '',
  path: '/api/tags',
  queryParams: [],
  headers: [],
  body: {
    type: 'none',
    rawContent: '',
    rawContentType: 'application/json',
    formData: [],
  },
});

const response = ref(null);

// --- API Definitions based on API.md ---
const baseApiDefinitions = [
  {
    key: 'generateCompletion',
    method: 'POST',
    path: '/api/generate',
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
    path: '/api/chat',
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
    path: '/api/create',
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
    path: '/api/tags',
    body: { type: 'none', rawContent: '', rawContentType: 'application/json' },
    queryParams: [],
    headers: [],
  },
  {
    key: 'showModelDetails',
    method: 'POST',
    path: '/api/show',
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
    path: '/api/copy',
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
    path: '/api/delete',
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
    path: '/api/pull',
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
    path: '/api/embed',
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
    path: '/api/ps',
    body: { type: 'none', rawContent: '', rawContentType: 'application/json' },
    queryParams: [],
    headers: [],
  },
  {
    key: 'generateSingleEmbedding',
    method: 'POST',
    path: '/api/embeddings',
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
      request.selectedServerId = activeServer.id;
    } else if (servers.value.length > 0) {
      // If no active server, select the first one
      request.selectedServerId = servers.value[0].id;
    }
  } catch (e) {
    error.value = `${t('apiDebugger.failedToLoadServers')}: ${e}`;
  }
}

function selectApi(api) {
  activeApiName.value = api.name;
  request.method = api.method;
  request.path = api.path;
  request.queryParams = api.queryParams ? JSON.parse(JSON.stringify(api.queryParams)) : [];
  request.headers = api.headers ? JSON.parse(JSON.stringify(api.headers)) : [];
  request.body.type = api.body.type;
  request.body.rawContent = api.body.rawContent || '';
  request.body.rawContentType = api.body.rawContentType || 'application/json';
  request.body.formData = api.body.formData ? JSON.parse(JSON.stringify(api.body.formData)) : [];

  // Automatically switch to body tab if there's a body, otherwise params
  if (api.body.type !== 'none' && api.body.rawContent) {
    activeTab.value = 'body';
  } else if (api.queryParams.length > 0) {
    activeTab.value = 'params';
  } else {
    activeTab.value = 'headers'; // Default to headers if no params or body
  }
  response.value = null; // Clear previous response
  error.value = null; // Clear previous error
}

async function sendRequest() {
  isLoading.value = true;
  error.value = null;
  response.value = null;

  try {
    // Construct the request object for the Go backend
    const backendRequest = {
      method: request.method,
      selectedServerId: request.selectedServerId,
      path: request.path,
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
      error.value = result.error;
    } else {
      response.value = result;
    }
  } catch (e) {
    error.value = `${t('apiDebugger.unexpectedError')}: ${e}`;
  } finally {
    isLoading.value = false;
  }
}

// Helper functions for managing params, headers, and form data
const addQueryParam = () => request.queryParams.push({ key: '', value: '', enabled: true });
const removeQueryParam = (index) => request.queryParams.splice(index, 1);

const addHeader = () => request.headers.push({ key: '', value: '', enabled: true });
const removeHeader = (index) => request.headers.splice(index, 1);

const addFormData = () => request.body.formData.push({ key: '', value: '' });
const removeFormData = (index) => request.body.formData.splice(index, 1);

// Computed property for formatted response body
const formattedBody = computed(() => {
  if (response.value && response.value.body) {
    try {
      // Check if the body is JSON
      const parsed = JSON.parse(response.value.body);
      return JSON.stringify(parsed, null, 2);
    } catch (e) {
      // Not JSON, return as is
      return response.value.body;
    }
  }
  return '';
});

const statusClass = computed(() => {
  if (!response.value) return '';
  if (response.value.statusCode >= 200 && response.value.statusCode < 300) return 'status-success';
  if (response.value.statusCode >= 400) return 'status-error';
  if (response.value.statusCode >= 300) return 'status-warning';
  return '';
});
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

.method-badge.get { background-color: #61affe; } /* Blue */
.method-badge.post { background-color: #49cc90; } /* Green */
.method-badge.put { background-color: #fca130; } /* Orange */
.method-badge.delete { background-color: #f93e3e; } /* Red */
.method-badge.patch { background-color: #50e3c2; } /* Teal */
.method-badge.head { background-color: #9012fe; } /* Purple */

.main-content {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  padding: 1rem;
  gap: 1rem;
  overflow: hidden; /* Prevent content from pushing the layout */
}

.request-panel,
.response-panel {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.request-panel {
  flex-shrink: 0; /* Prevent request panel from shrinking */
}

.response-panel {
  flex-grow: 1;
  min-height: 0; /* Crucial for flex-grow in a flex column */
  display: flex;
  flex-direction: column;
}

.response-wrapper {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  min-height: 0;
}

.request-url-bar {
  display: flex;
  gap: 0.5rem;
  padding: 1rem;
  border-bottom: 1px solid #eee;
  align-items: center;
}

.http-method-display {
  padding: 0.5rem 0.8rem;
  border-radius: 4px;
  font-weight: bold;
  color: #fff;
  min-width: 70px;
  text-align: center;
  font-size: 0.9rem;
}

.http-method-display.get { background-color: #61affe; }
.http-method-display.post { background-color: #49cc90; }
.http-method-display.put { background-color: #fca130; }
.http-method-display.delete { background-color: #f93e3e; }
.http-method-display.patch { background-color: #50e3c2; }
.http-method-display.head { background-color: #9012fe; }

.url-input {
  flex-grow: 1;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 0.9rem;
}

.request-url-bar button {
  padding: 0.5rem 1rem;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
}

.request-url-bar button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.response-tabs {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  min-height: 0;
}

.tab-headers {
  display: flex;
  border-bottom: 1px solid #eee;
  padding: 0 1rem;
  flex-shrink: 0;
}

.tab-headers button {
  padding: 0.8rem 1.2rem;
  border: none;
  background-color: transparent;
  cursor: pointer;
  font-size: 0.9rem;
  color: #555;
  border-bottom: 2px solid transparent;
  transition: all 0.2s ease;
}

.tab-headers button:hover {
  color: #007bff;
}

.tab-headers button.active {
  border-bottom: 2px solid #007bff;
  font-weight: bold;
  color: #007bff;
}

.tab-content-request {
  padding: 1rem;
}

.tab-content-response {
  flex-grow: 1;
  overflow-y: auto;
  padding: 1rem;
}

.key-value-editor {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.kv-row,
.kv-row-readonly {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.kv-row input[type='text'],
.kv-row-readonly input[type='text'] {
  flex-grow: 1;
  padding: 0.4rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 0.9rem;
}

.kv-row button {
  padding: 0.4rem 0.8rem;
  background-color: #dc3545;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
}

.body-type-selector {
  margin-bottom: 1rem;
  display: flex;
  gap: 1rem;
}

.body-type-selector label {
  font-size: 0.9rem;
  color: #333;
}

.raw-body-editor select {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  margin-bottom: 0.5rem;
  background-color: #f9f9f9;
}

.raw-body-editor textarea {
  width: 100%;
  min-height: 150px;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-family: 'Fira Code', 'Cascadia Code', monospace;
  font-size: 0.9rem;
  resize: vertical;
  background-color: #fdfdfd;
}

.response-status-bar {
  display: flex;
  gap: 1.5rem;
  padding: 1rem;
  border-bottom: 1px solid #eee;
  font-size: 0.9rem;
  color: #333;
  flex-shrink: 0;
}

.response-status-bar strong {
  font-weight: bold;
}

.status-success { color: #28a745; } 
.status-error { color: #dc3545; } 
.status-warning { color: #ffc107; } 

.response-body {
  background-color: #f8f8f8;
  border: 1px solid #eee;
  padding: 0.8rem;
  white-space: pre-wrap;
  word-wrap: break-word;
  border-radius: 4px;
  font-family: 'Fira Code', 'Cascadia Code', monospace;
  font-size: 0.9rem;
  color: #333;
  text-align: left;
}

.response-placeholder {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #888;
}

.error-message {
  color: #dc3545;
  padding: 1rem;
  background-color: #ffebeb;
  border-radius: 8px;
  border: 1px solid #f5c6cb;
}

.error-message h3 {
  margin-top: 0;
  color: #721c24;
}

.error-message pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: 'Fira Code', 'Cascadia Code', monospace;
  font-size: 0.85rem;
  color: #721c24;
}
</style>
