<template>
  <div class="api-debugger-widget">
    <div class="api-header">
      <span :class="['method-badge', method.toLowerCase()]">{{ method }}</span>
      <span class="api-path">{{ path }}</span>
      <button class="test-btn" @click="sendRequest" :disabled="isLoading">
        {{ isLoading ? t('apiDebugger.sending') : t('common.test') }}
      </button>
    </div>
    
    <div class="request-body" v-if="requestBody">
      <h4>{{ t('apiDebugger.body') }}</h4>
      <pre class="json-content">{{ requestBody }}</pre>
    </div>
    
    <div class="response-section" v-if="response || error">
      <div v-if="error" class="error-message">
        <h4>{{ t('apiDebugger.error') }}</h4>
        <pre>{{ error }}</pre>
      </div>
      <div v-else-if="response" class="response-content">
        <div class="response-status">
          <span :class="statusClass">{{ response.statusCode }} {{ response.statusText }}</span>
          <span class="duration">{{ response.requestDurationMs }}ms</span>
        </div>
        <div class="response-body">
          <h4>{{ t('apiDebugger.responseBody') }}</h4>
          <pre class="json-content">{{ formattedResponseBody }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { SendHttpRequest } from '../../wailsjs/go/main/App';
import { types } from '../../wailsjs/go/models';

const { t } = useI18n();

interface Props {
  method: string;
  path: string;
  requestBody?: string;
  serverId: string;
}

const props = defineProps<Props>();

const isLoading = ref(false);
const response = ref<any>(null);
const error = ref<string | null>(null);

const formattedResponseBody = computed(() => {
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

const statusClass = computed(() => {
  if (!response.value) return '';
  const code = response.value.statusCode;
  if (code >= 200 && code < 300) return 'status-success';
  if (code >= 400) return 'status-error';
  if (code >= 300) return 'status-warning';
  return '';
});

const sendRequest = async () => {
  if (!props.serverId) {
    error.value = t('apiDebugger.configureServerFirst');
    return;
  }

  isLoading.value = true;
  error.value = null;
  response.value = null;

  try {
    const request: types.ApiRequest = {
      method: props.method,
      selectedServerId: props.serverId,
      path: props.path,
      queryParams: [],
      headers: props.requestBody ? [
        { key: 'Content-Type', value: 'application/json', enabled: true }
      ] : [],
      body: {
        type: props.requestBody ? 'raw' : 'none',
        rawContent: props.requestBody || '',
        rawContentType: 'application/json',
        formData: []
      }
    };

    const result = await SendHttpRequest(request);
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
</script>

<style scoped>
.api-debugger-widget {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  margin-bottom: 1rem;
  background: white;
}

.api-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
}

.method-badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: bold;
  color: white;
  min-width: 50px;
  text-align: center;
}

.method-badge.get { background-color: #61affe; }
.method-badge.post { background-color: #49cc90; }
.method-badge.put { background-color: #fca130; }
.method-badge.delete { background-color: #f93e3e; }

.api-path {
  flex-grow: 1;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', monospace;
  font-size: 0.9rem;
  color: #374151;
}

.test-btn {
  padding: 0.5rem 1rem;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.2s ease;
}

.test-btn:hover:not(:disabled) {
  background: #5a67d8;
}

.test-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.request-body, .response-section {
  padding: 1rem;
}

.request-body h4, .response-section h4 {
  margin: 0 0 0.5rem 0;
  font-size: 0.9rem;
  color: #4a5568;
}

.json-content {
  background: #f7fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 1rem;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', monospace;
  font-size: 0.85rem;
  color: #2d3748;
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
}

.response-status {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  font-size: 0.9rem;
}

.status-success { color: #38a169; }
.status-error { color: #e53e3e; }
.status-warning { color: #d69e2e; }

.duration {
  color: #718096;
  font-weight: 500;
}

.error-message {
  background: #fed7d7;
  border: 1px solid #feb2b2;
  border-radius: 6px;
  padding: 1rem;
}

.error-message h4 {
  color: #c53030;
  margin: 0 0 0.5rem 0;
}

.error-message pre {
  color: #c53030;
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>