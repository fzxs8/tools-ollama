<template>
  <div class="api-debugger">
    <!-- Request Panel -->
    <div class="request-panel">
      <div class="request-url-bar">
        <select v-model="request.method" class="method-selector">
          <option value="GET">GET</option>
          <option value="POST">POST</option>
          <option value="PUT">PUT</option>
          <option value="DELETE">DELETE</option>
          <option value="PATCH">PATCH</option>
          <option value="HEAD">HEAD</option>
        </select>
        <input 
          type="text" 
          v-model="request.url" 
          :placeholder="computedUrlPlaceholder"
          class="url-input" 
        />
        <button @click="sendRequest" :disabled="isLoading" class="send-btn">
          {{ isLoading ? t('apiDebugger.sending') : t('apiDebugger.send') }}
        </button>
      </div>

      <!-- Request Tabs -->
      <div class="request-tabs">
        <div class="tab-headers">
          <button 
            :class="{ active: activeTab === 'params' }" 
            @click="activeTab = 'params'"
          >
            {{ t('apiDebugger.queryParams') }}
          </button>
          <button 
            :class="{ active: activeTab === 'headers' }" 
            @click="activeTab = 'headers'"
          >
            {{ t('apiDebugger.headers') }}
          </button>
          <button 
            :class="{ active: activeTab === 'body' }" 
            @click="activeTab = 'body'"
          >
            {{ t('apiDebugger.body') }}
          </button>
        </div>
        
        <div class="tab-content">
          <!-- Query Params -->
          <div v-if="activeTab === 'params'" class="key-value-editor">
            <div v-for="(param, index) in request.queryParams" :key="index" class="kv-row">
              <input type="checkbox" v-model="param.enabled" />
              <input type="text" v-model="param.key" :placeholder="t('apiDebugger.key')" />
              <input type="text" v-model="param.value" :placeholder="t('apiDebugger.value')" />
              <button @click="removeQueryParam(index)" class="remove-btn">{{ t('apiDebugger.remove') }}</button>
            </div>
            <button @click="addQueryParam" class="add-btn">{{ t('apiDebugger.addParameter') }}</button>
          </div>

          <!-- Headers -->
          <div v-if="activeTab === 'headers'" class="key-value-editor">
            <div v-for="(header, index) in request.headers" :key="index" class="kv-row">
              <input type="checkbox" v-model="header.enabled" />
              <input type="text" v-model="header.key" :placeholder="t('apiDebugger.key')" />
              <input type="text" v-model="header.value" :placeholder="t('apiDebugger.value')" />
              <button @click="removeHeader(index)" class="remove-btn">{{ t('apiDebugger.remove') }}</button>
            </div>
            <button @click="addHeader" class="add-btn">{{ t('apiDebugger.addHeader') }}</button>
          </div>

          <!-- Body -->
          <div v-if="activeTab === 'body'">
            <div class="body-type-selector">
              <label><input type="radio" v-model="request.body.type" value="none" /> {{ t('apiDebugger.none') }}</label>
              <label><input type="radio" v-model="request.body.type" value="raw" /> {{ t('apiDebugger.raw') }}</label>
              <label><input type="radio" v-model="request.body.type" value="formData" /> {{ t('apiDebugger.formData') }}</label>
            </div>
            
            <div v-if="request.body.type === 'raw'" class="raw-body-editor">
              <div class="raw-body-header">
                <select v-model="request.body.rawContentType" class="content-type-selector">
                  <option value="application/json">JSON</option>
                  <option value="text/plain">Text</option>
                  <option value="text/html">HTML</option>
                  <option value="application/xml">XML</option>
                </select>
                <div class="example-buttons" v-if="request.body.rawContentType === 'application/json'">
                  <button @click="loadChatExample" class="example-btn">{{ t('apiDebugger.chatExample') }}</button>
                  <button @click="loadStreamExample" class="example-btn">{{ t('apiDebugger.streamExample') }}</button>
                </div>
              </div>
              <textarea 
                v-model="request.body.rawContent" 
                :placeholder="t('apiDebugger.enterRawContent')"
                class="body-textarea"
              ></textarea>
            </div>
            
            <div v-if="request.body.type === 'formData'" class="key-value-editor">
              <div v-for="(item, index) in request.body.formData" :key="index" class="kv-row">
                <input type="text" v-model="item.key" :placeholder="t('apiDebugger.key')" />
                <input type="text" v-model="item.value" :placeholder="t('apiDebugger.value')" />
                <button @click="removeFormData(index)" class="remove-btn">{{ t('apiDebugger.remove') }}</button>
              </div>
              <button @click="addFormData" class="add-btn">{{ t('apiDebugger.addFormItem') }}</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Response Panel -->
    <div class="response-panel">
      <div v-if="isLoading" class="response-placeholder">
        <div class="loading-spinner"></div>
        <p>{{ t('apiDebugger.loading') }}</p>
      </div>
      
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
            <button 
              :class="{ active: activeResponseTab === 'body' }" 
              @click="activeResponseTab = 'body'"
            >
              {{ t('apiDebugger.responseBody') }}
            </button>
            <button 
              :class="{ active: activeResponseTab === 'headers' }" 
              @click="activeResponseTab = 'headers'"
            >
              {{ t('apiDebugger.responseHeaders') }}
            </button>
          </div>
          
          <div class="tab-content">
            <div v-if="activeResponseTab === 'body'">
              <pre class="response-body">{{ formattedBody }}</pre>
            </div>
            <div v-if="activeResponseTab === 'headers'" class="key-value-editor">
              <div v-for="(header, index) in response.headers" :key="index" class="kv-row-readonly">
                <input type="text" :value="header.key" readonly />
                <input type="text" :value="header.value" readonly />
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
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

interface KeyValue {
  key: string;
  value: string;
  enabled: boolean;
}

interface FormDataItem {
  key: string;
  value: string;
}

interface RequestBody {
  type: 'none' | 'raw' | 'formData';
  rawContent: string;
  rawContentType: string;
  formData: FormDataItem[];
}

interface ApiRequest {
  method: string;
  url: string;
  queryParams: KeyValue[];
  headers: KeyValue[];
  body: RequestBody;
}

interface Props {
  baseUrl?: string;
  urlPlaceholder?: string;
  onRequest?: (request: ApiRequest) => Promise<any>;
}

const props = withDefaults(defineProps<Props>(), {
  baseUrl: '',
  urlPlaceholder: '',
  onRequest: undefined
});

// 使用国际化的默认占位符
const computedUrlPlaceholder = computed(() => 
  props.urlPlaceholder || t('apiDebugger.enterUrl')
);

const activeTab = ref('params');
const activeResponseTab = ref('body');
const isLoading = ref(false);
const error = ref<string | null>(null);
const response = ref<any>(null);

const request = reactive<ApiRequest>({
  method: 'GET',
  url: '',
  queryParams: [],
  headers: [],
  body: {
    type: 'none',
    rawContent: '',
    rawContentType: 'application/json',
    formData: []
  }
});

// 构建完整的 URL
const fullUrl = computed(() => {
  let url = request.url;
  if (props.baseUrl && !url.startsWith('http')) {
    url = props.baseUrl.replace(/\/$/, '') + '/' + url.replace(/^\//, '');
  }
  
  const enabledParams = request.queryParams.filter(p => p.enabled && p.key);
  if (enabledParams.length > 0) {
    const params = new URLSearchParams();
    enabledParams.forEach(p => params.append(p.key, p.value));
    url += (url.includes('?') ? '&' : '?') + params.toString();
  }
  
  return url;
});

// 格式化响应体
const formattedBody = computed(() => {
  if (response.value?.body) {
    try {
      const parsed = JSON.parse(response.value.body);
      return JSON.stringify(parsed, null, 2);
    } catch (e) {
      return response.value.body;
    }
  }
  return '';
});

// 状态样式
const statusClass = computed(() => {
  if (!response.value) return '';
  const code = response.value.statusCode;
  if (code >= 200 && code < 300) return 'status-success';
  if (code >= 400) return 'status-error';
  if (code >= 300) return 'status-warning';
  return '';
});

// 参数管理
const addQueryParam = () => request.queryParams.push({ key: '', value: '', enabled: true });
const removeQueryParam = (index: number) => request.queryParams.splice(index, 1);

const addHeader = () => request.headers.push({ key: '', value: '', enabled: true });
const removeHeader = (index: number) => request.headers.splice(index, 1);

const addFormData = () => request.body.formData.push({ key: '', value: '' });
const removeFormData = (index: number) => request.body.formData.splice(index, 1);

// 示例加载
const loadChatExample = () => {
  request.body.rawContent = JSON.stringify({
    model: 'llama3',
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
  }, null, 2);
};

const loadStreamExample = () => {
  request.body.rawContent = JSON.stringify({
    model: 'llama3',
    messages: [
      {
        role: 'system',
        content: 'You are a helpful assistant.'
      },
      {
        role: 'user',
        content: 'Tell me a joke.'
      }
    ],
    temperature: 0.8,
    max_tokens: 500,
    stream: true
  }, null, 2);
};

// 发送请求
const sendRequest = async () => {
  if (!request.url) {
    error.value = t('apiDebugger.enterUrl');
    return;
  }

  isLoading.value = true;
  error.value = null;
  response.value = null;

  try {
    let result;
    
    if (props.onRequest) {
      // 使用自定义请求处理器
      result = await props.onRequest(request);
    } else {
      // 使用默认的 fetch 请求
      const headers: Record<string, string> = {};
      request.headers.filter(h => h.enabled && h.key).forEach(h => {
        headers[h.key] = h.value;
      });

      let body: string | FormData | undefined;
      if (request.body.type === 'raw' && request.body.rawContent) {
        body = request.body.rawContent;
        if (request.body.rawContentType && !headers['Content-Type']) {
          headers['Content-Type'] = request.body.rawContentType;
        }
      } else if (request.body.type === 'formData') {
        const formData = new FormData();
        request.body.formData.filter(f => f.key).forEach(f => {
          formData.append(f.key, f.value);
        });
        body = formData;
      }

      const startTime = Date.now();
      const fetchResponse = await fetch(fullUrl.value, {
        method: request.method,
        headers,
        body: ['GET', 'HEAD'].includes(request.method) ? undefined : body
      });
      const endTime = Date.now();

      const responseHeaders = Array.from(fetchResponse.headers.entries()).map(([key, value]) => ({
        key,
        value
      }));

      result = {
        statusCode: fetchResponse.status,
        statusText: fetchResponse.statusText,
        headers: responseHeaders,
        body: await fetchResponse.text(),
        requestDurationMs: endTime - startTime
      };
    }

    if (result.error) {
      error.value = result.error;
    } else {
      response.value = result;
    }
  } catch (e) {
    error.value = `${t('apiDebugger.requestFailed')}: ${e}`;
  } finally {
    isLoading.value = false;
  }
};

// 暴露方法供外部调用
defineExpose({
  setRequest: (newRequest: Partial<ApiRequest>) => {
    Object.assign(request, newRequest);
  },
  sendRequest,
  clearResponse: () => {
    response.value = null;
    error.value = null;
  }
});
</script>

<style scoped>
.api-debugger {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #f8fafc;
  border-radius: 8px;
  overflow: hidden;
}

.request-panel,
.response-panel {
  background-color: #fff;
  border: 1px solid #e2e8f0;
}

.request-panel {
  flex-shrink: 0;
  margin-bottom: 1rem;
}

.response-panel {
  flex-grow: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.request-url-bar {
  display: flex;
  gap: 0.5rem;
  padding: 1rem;
  border-bottom: 1px solid #e2e8f0;
  align-items: center;
}

.method-selector {
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background-color: #f9fafb;
  font-weight: 600;
  min-width: 80px;
}

.url-input {
  flex-grow: 1;
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 0.9rem;
}

.send-btn {
  padding: 0.5rem 1rem;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  transition: background-color 0.2s;
}

.send-btn:hover:not(:disabled) {
  background-color: #2563eb;
}

.send-btn:disabled {
  background-color: #9ca3af;
  cursor: not-allowed;
}

.request-tabs,
.response-tabs {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  min-height: 0;
}

.tab-headers {
  display: flex;
  border-bottom: 1px solid #e2e8f0;
  padding: 0 1rem;
  flex-shrink: 0;
}

.tab-headers button {
  padding: 0.75rem 1rem;
  border: none;
  background-color: transparent;
  cursor: pointer;
  font-size: 0.9rem;
  color: #6b7280;
  border-bottom: 2px solid transparent;
  transition: all 0.2s;
}

.tab-headers button:hover {
  color: #3b82f6;
}

.tab-headers button.active {
  border-bottom-color: #3b82f6;
  color: #3b82f6;
  font-weight: 600;
}

.tab-content {
  padding: 1rem;
  flex-grow: 1;
  overflow-y: auto;
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
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 0.9rem;
}

.kv-row-readonly input[type='text'] {
  background-color: #f9fafb;
}

.remove-btn {
  padding: 0.4rem 0.8rem;
  background-color: #ef4444;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.8rem;
}

.add-btn {
  padding: 0.5rem 1rem;
  background-color: #10b981;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  align-self: flex-start;
}

.body-type-selector {
  margin-bottom: 1rem;
  display: flex;
  gap: 1rem;
}

.body-type-selector label {
  font-size: 0.9rem;
  color: #374151;
  cursor: pointer;
}

.raw-body-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
  gap: 1rem;
}

.content-type-selector {
  flex: 1;
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background-color: #f9fafb;
}

.example-buttons {
  display: flex;
  gap: 0.5rem;
}

.example-btn {
  padding: 0.4rem 0.8rem;
  background-color: #6366f1;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.8rem;
  transition: background-color 0.2s;
}

.example-btn:hover {
  background-color: #4f46e5;
}

.body-textarea {
  width: 100%;
  min-height: 150px;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', monospace;
  font-size: 0.9rem;
  resize: vertical;
  background-color: #fdfdfd;
}

.response-wrapper {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  min-height: 0;
}

.response-status-bar {
  display: flex;
  gap: 1.5rem;
  padding: 1rem;
  border-bottom: 1px solid #e2e8f0;
  font-size: 0.9rem;
  color: #374151;
  flex-shrink: 0;
}

.status-success { color: #059669; }
.status-error { color: #dc2626; }
.status-warning { color: #d97706; }

.response-body {
  background-color: #f8fafc;
  border: 1px solid #e2e8f0;
  padding: 1rem;
  white-space: pre-wrap;
  word-wrap: break-word;
  border-radius: 6px;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', monospace;
  font-size: 0.85rem;
  color: #374151;
  margin: 0;
}

.response-placeholder {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #6b7280;
  padding: 2rem;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #e5e7eb;
  border-top: 3px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-message {
  color: #dc2626;
  background-color: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 8px;
  padding: 1rem;
  margin: 1rem;
}

.error-message h3 {
  margin-top: 0;
  color: #991b1b;
}

.error-message pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', monospace;
  font-size: 0.85rem;
  color: #991b1b;
  margin: 0.5rem 0 0 0;
}
</style>