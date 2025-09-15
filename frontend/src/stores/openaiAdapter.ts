import { defineStore } from 'pinia';
import { ElMessage } from 'element-plus';
import i18n from '../i18n';

const { t } = i18n.global;
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
    async fetchConfig() {
      try {
        const fetchedConfig = await GetOpenAIAdapterConfig();
        if (fetchedConfig) {
          this.config.listenIp = fetchedConfig.listenIp || '127.0.0.1';
          this.config.listenPort = fetchedConfig.listenPort || 11223;
          this.config.targetOllamaServerId = fetchedConfig.targetOllamaServerId;
        }
      } catch (error) {
        ElMessage.error(`${t('messages.fetchConfigFailed')}: ${error}`);
      }
    },

    async saveConfig() {
      try {
        await SaveOpenAIAdapterConfig(this.config);
        ElMessage.success(t('messages.configSaved'));
      } catch (error) {
        ElMessage.error(`${t('messages.configSaveFailed')}: ${error}`);
        throw error;
      }
    },

    async fetchInitialStatus() {
      try {
        this.status = await GetOpenAIAdapterStatus();
      } catch (error) {
        this.status.isRunning = false;
        this.status.error = `${t('messages.getServiceStatusFailed')}: ${error}`;
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
        ElMessage.success(t('messages.serviceStarted'));
      } catch (error) {
        ElMessage.error(`${t('messages.serviceStartFailed')}: ${error}`);
        // The backend will emit a status update on failure, so no need to fetch here
        throw error;
      }
    },

    async stopServer() {
      try {
        await StopAdapterServer();
        ElMessage.success(t('messages.serviceStopped'));
      } catch (error) {
        ElMessage.error(`${t('messages.serviceStopFailed')}: ${error}`);
        throw error;
      }
    },

    async fetchOllamaServers() {
      try {
        this.ollamaServers = await GetOllamaServers();
      } catch (error) {
        ElMessage.error(`${t('messages.fetchServersFailed')}: ${error}`);
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
      ElMessage.success(t('messages.logsClearedSuccess'));
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
        ElMessage.error(`${t('messages.fetchApiDocsFailed')}: ${error}`);
      }
    },
  },
});
