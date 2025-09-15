import { defineStore } from 'pinia';
import { ElMessage } from 'element-plus';
import { types } from '../../wailsjs/go/models';
import OpenAIAdapterConfig = types.OpenAIAdapterConfig;
import OpenAIAdapterStatus = types.OpenAIAdapterStatus;
import OllamaServerConfig = types.OllamaServerConfig;
import LogEntry = types.LogEntry;
import {
    GetAdapterAPIDocs,
    GetOllamaServers,
    GetOpenAIAdapterConfig,
    GetOpenAIAdapterStatus,
    SaveOpenAIAdapterConfig,
    StartAdapterServer,
    StopAdapterServer
} from "../../wailsjs/go/main/App";

const MAX_LOGS = 500;

interface OpenAIAdapterState {
  config: OpenAIAdapterConfig;
  status: OpenAIAdapterStatus;
  ollamaServers: OllamaServerConfig[];
  logs: LogEntry[];
  isLogDrawerVisible: boolean;
  apiDocs: Record<string, string>;
  isApiDrawerVisible: boolean;
}

export const useOpenAIAdapterStore = defineStore('openaiAdapter', {
  state: (): OpenAIAdapterState => ({
    config: {
      listenIp: '127.0.0.1',
      listenPort: 11223,
      targetOllamaServerId: '',
    },
    status: {
      isRunning: false,
      error: '',
    },
    ollamaServers: [],
    logs: [],
    isLogDrawerVisible: false,
    apiDocs: {},
    isApiDrawerVisible: false,
  }),

  actions: {
    // IP地址验证方法
    isValidIP(ip: string): boolean {
      const ipRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
      return ipRegex.test(ip) || ip === 'localhost';
    },

    async fetchConfig() {
      try {
        const fetchedConfig = await GetOpenAIAdapterConfig();
        if (fetchedConfig) {
          this.config.listenIp = fetchedConfig.listenIp || '127.0.0.1';
          this.config.listenPort = fetchedConfig.listenPort || 11223;
          this.config.targetOllamaServerId = fetchedConfig.targetOllamaServerId;
        }
      } catch (error) {
        ElMessage.error(`Failed to fetch adapter configuration: ${error}`);
      }
    },

    async saveConfig() {
      try {
        await SaveOpenAIAdapterConfig(this.config);
        ElMessage.success('Configuration saved successfully!');
      } catch (error) {
        ElMessage.error(`Failed to save configuration: ${error}`);
        throw error;
      }
    },

    async fetchInitialStatus() {
      try {
        this.status = await GetOpenAIAdapterStatus();
      } catch (error) {
        this.status.isRunning = false;
        this.status.error = `Unable to get service status: ${error}`;
      }
    },

    // Action to be called by the event listener
    setStatus(newStatus: OpenAIAdapterStatus) {
      this.status = newStatus;
    },

    async startServer() {
      try {
        // 全面验证配置
        const validationErrors = [];
        
        if (!this.config.listenIp || this.config.listenIp.trim() === '') {
          validationErrors.push('请输入监听地址');
        } else if (!this.isValidIP(this.config.listenIp.trim())) {
          validationErrors.push('请输入有效的IP地址');
        }
        
        if (!this.config.listenPort || this.config.listenPort < 1 || this.config.listenPort > 65535) {
          validationErrors.push('请输入有效的端口号 (1-65535)');
        }
        
        if (!this.config.targetOllamaServerId) {
          validationErrors.push('请选择目标 Ollama 服务器');
        }
        
        if (validationErrors.length > 0) {
          ElMessage.error('配置不完整：' + validationErrors.join('；'));
          return;
        }
        
        await this.saveConfig();
        await StartAdapterServer();
        ElMessage.success('Service start command sent.');
      } catch (error) {
        ElMessage.error(`Failed to start service: ${error}`);
        // The backend will emit a status update on failure, so no need to fetch here
        throw error;
      }
    },

    async stopServer() {
      try {
        await StopAdapterServer();
        ElMessage.success('Service stop command sent.');
      } catch (error) {
        ElMessage.error(`Failed to stop service: ${error}`);
        throw error;
      }
    },

    async fetchOllamaServers() {
      try {
        this.ollamaServers = await GetOllamaServers();
      } catch (error) {
        ElMessage.error(`Failed to fetch Ollama server list: ${error}`);
      }
    },

    addLog(log: LogEntry) {
      this.logs.push(log);
      if (this.logs.length > MAX_LOGS) {
        this.logs.shift();
      }
    },

    clearLogs() {
      this.logs = [];
      ElMessage.success('Logs cleared');
    },

    toggleLogDrawer(visible?: boolean) {
      this.isLogDrawerVisible = visible ?? !this.isLogDrawerVisible;
    },

    toggleApiDrawer(visible?: boolean) {
      this.isApiDrawerVisible = visible ?? !this.isApiDrawerVisible;
    },

    async fetchApiDocs() {
      if (Object.keys(this.apiDocs).length > 0) {
        this.toggleApiDrawer(true);
        return;
      }
      try {
        this.apiDocs = await GetAdapterAPIDocs();
        this.toggleApiDrawer(true);
      } catch (error) {
        ElMessage.error(`Failed to fetch API documentation: ${error}`);
      }
    },
  },
});
