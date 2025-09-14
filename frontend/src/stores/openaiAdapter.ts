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
      listen_ip: '127.0.0.1',
      listen_port: 11223,
      target_ollama_server_id: '',
    },
    status: {
      is_running: false,
      error: '',
    },
    ollamaServers: [],
    logs: [],
    isLogDrawerVisible: false,
    apiDocs: {},
    isApiDrawerVisible: false,
  }),

  actions: {
    async fetchConfig() {
      try {
        const fetchedConfig = await GetOpenAIAdapterConfig();
        if (fetchedConfig) {
          this.config.listen_ip = fetchedConfig.listen_ip || '127.0.0.1';
          this.config.listen_port = fetchedConfig.listen_port || 11223;
          this.config.target_ollama_server_id = fetchedConfig.target_ollama_server_id;
        }
      } catch (error) {
        ElMessage.error(`获取适配器配置失败: ${error}`);
      }
    },

    async saveConfig() {
      try {
        await SaveOpenAIAdapterConfig(this.config);
        ElMessage.success('配置保存成功！');
      } catch (error) {
        ElMessage.error(`保存配置失败: ${error}`);
        throw error;
      }
    },

    async fetchInitialStatus() {
      try {
        this.status = await GetOpenAIAdapterStatus();
      } catch (error) {
        this.status.is_running = false;
        this.status.error = `无法获取服务状态: ${error}`;
      }
    },

    // Action to be called by the event listener
    setStatus(newStatus: OpenAIAdapterStatus) {
      this.status = newStatus;
    },

    async startServer() {
      try {
        await this.saveConfig();
        await StartAdapterServer();
        ElMessage.success('服务启动指令已发送。');
      } catch (error) {
        ElMessage.error(`服务启动失败: ${error}`);
        // The backend will emit a status update on failure, so no need to fetch here
        throw error;
      }
    },

    async stopServer() {
      try {
        await StopAdapterServer();
        ElMessage.success('服务停止指令已发送。');
      } catch (error) {
        ElMessage.error(`服务停止失败: ${error}`);
        throw error;
      }
    },

    async fetchOllamaServers() {
      try {
        this.ollamaServers = await GetOllamaServers();
      } catch (error) {
        ElMessage.error(`获取 Ollama 服务列表失败: ${error}`);
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
      ElMessage.success('日志已清空');
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
        ElMessage.error(`获取 API 文档失败: ${error}`);
      }
    },
  },
});
