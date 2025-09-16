<template>
  <div class="system-monitor">
    <!-- Page Header -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-icon">
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2" stroke="currentColor" stroke-width="2"/>
            <circle cx="9" cy="9" r="2" stroke="currentColor" stroke-width="2"/>
            <path d="M21 15.5A3.5 3.5 0 0 0 17.5 12H9A3.5 3.5 0 0 0 5.5 15.5A3.5 3.5 0 0 0 9 19Z" stroke="currentColor" stroke-width="2"/>
          </svg>
        </div>
        <div class="header-text">
          <h1>System Monitor</h1>
          <p>Monitor system performance and Ollama service status</p>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="main-content">
      <!-- Control Panel -->
      <div class="control-panel">
        <div class="panel-title">
          <h3>System Overview</h3>
          <p>Real-time system performance metrics</p>
        </div>
        <button class="btn primary" @click="refreshData">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
            <polyline points="23,4 23,10 17,10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <polyline points="1,20 1,14 7,14" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M20.49 9A9 9 0 0 0 5.64 5.64L1 10M22.88 14.36A9 9 0 0 1 8.51 18.36L4 14" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          Refresh Data
        </button>
      </div>

      <!-- System Stats -->
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-header">
            <div class="stat-icon cpu">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <rect x="4" y="4" width="16" height="16" rx="2" stroke="currentColor" stroke-width="2"/>
                <rect x="9" y="9" width="6" height="6" stroke="currentColor" stroke-width="2"/>
                <line x1="9" y1="1" x2="9" y2="4" stroke="currentColor" stroke-width="2"/>
                <line x1="15" y1="1" x2="15" y2="4" stroke="currentColor" stroke-width="2"/>
                <line x1="9" y1="20" x2="9" y2="23" stroke="currentColor" stroke-width="2"/>
                <line x1="15" y1="20" x2="15" y2="23" stroke="currentColor" stroke-width="2"/>
                <line x1="20" y1="9" x2="23" y2="9" stroke="currentColor" stroke-width="2"/>
                <line x1="20" y1="14" x2="23" y2="14" stroke="currentColor" stroke-width="2"/>
                <line x1="1" y1="9" x2="4" y2="9" stroke="currentColor" stroke-width="2"/>
                <line x1="1" y1="14" x2="4" y2="14" stroke="currentColor" stroke-width="2"/>
              </svg>
            </div>
            <div class="stat-info">
              <h4>CPU Usage</h4>
              <div class="stat-value">{{ cpuUsage }}%</div>
            </div>
          </div>
          <div class="progress-container">
            <div class="progress-bar">
              <div class="progress-fill cpu" :style="{ width: cpuUsage + '%' }"></div>
            </div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-header">
            <div class="stat-icon memory">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <path d="M2 3H6A4 4 0 0 1 10 7V17A3 3 0 0 0 13 20A3 3 0 0 0 16 17V7A4 4 0 0 1 20 3H22" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </div>
            <div class="stat-info">
              <h4>Memory Usage</h4>
              <div class="stat-value">{{ memoryUsage }}%</div>
            </div>
          </div>
          <div class="progress-container">
            <div class="progress-bar">
              <div class="progress-fill memory" :style="{ width: memoryUsage + '%' }"></div>
            </div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-header">
            <div class="stat-icon disk">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
                <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
              </svg>
            </div>
            <div class="stat-info">
              <h4>Disk Usage</h4>
              <div class="stat-value">{{ diskUsage }}%</div>
            </div>
          </div>
          <div class="progress-container">
            <div class="progress-bar">
              <div class="progress-fill disk" :style="{ width: diskUsage + '%' }"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Process Table -->
      <div class="process-table-container">
        <div class="table-header">
          <h3>Ollama Service Status</h3>
          <p>Active Ollama processes and their resource usage</p>
        </div>
        
        <div class="table-wrapper">
          <table class="process-table">
            <thead>
              <tr>
                <th>Process ID</th>
                <th>Process Name</th>
                <th>CPU %</th>
                <th>Memory %</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="ollamaProcesses.length === 0">
                <td colspan="5" class="empty-state">
                  <div class="empty-content">
                    <svg width="48" height="48" viewBox="0 0 24 24" fill="none">
                      <rect x="3" y="3" width="18" height="18" rx="2" ry="2" stroke="currentColor" stroke-width="2"/>
                      <circle cx="9" cy="9" r="2" stroke="currentColor" stroke-width="2"/>
                      <path d="M21 15.5A3.5 3.5 0 0 0 17.5 12H9A3.5 3.5 0 0 0 5.5 15.5A3.5 3.5 0 0 0 9 19Z" stroke="currentColor" stroke-width="2"/>
                    </svg>
                    <p>No Ollama processes found</p>
                    <span>Ollama service may not be running</span>
                  </div>
                </td>
              </tr>
              <tr v-for="process in ollamaProcesses" :key="process.pid" class="process-row">
                <td class="process-pid">{{ process.pid }}</td>
                <td class="process-name">
                  <div class="name-info">
                    <span class="name">{{ process.name }}</span>
                  </div>
                </td>
                <td>{{ process.cpu_percent }}</td>
                <td>{{ process.memory_percent }}</td>
                <td>
                  <span class="status-badge running">
                    <span class="status-dot"></span>
                    {{ process.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {onMounted, ref} from 'vue'

const cpuUsage = ref(0)
const memoryUsage = ref(0)
const diskUsage = ref(0)
const ollamaProcesses = ref<any[]>([])


const refreshData = () => {
  // Simulate refreshing data
  cpuUsage.value = Math.floor(Math.random() * 100)
  memoryUsage.value = Math.floor(Math.random() * 100)
  diskUsage.value = Math.floor(Math.random() * 100)
  
  // Simulate process data
  ollamaProcesses.value = [
    { pid: 1234, name: 'ollama', cpu_percent: '25.6%', memory_percent: '15.2%', status: 'Running' },
    { pid: 5678, name: 'ollama-run', cpu_percent: '10.2%', memory_percent: '8.7%', status: 'Running' }
  ]
}

onMounted(() => {
  refreshData()
})
</script>

<style scoped>
.system-monitor {
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 2rem;
  box-sizing: border-box;
  overflow-y: auto;
}

.page-header {
  margin-bottom: 2rem;
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
  border-radius: 20px;
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
  border-radius: 12px;
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

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.stat-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  padding: 1.5rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.stat-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-icon.cpu {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.memory {
  background: linear-gradient(135deg, #48bb78 0%, #38a169 100%);
}

.stat-icon.disk {
  background: linear-gradient(135deg, #ed8936 0%, #dd6b20 100%);
}

.stat-info h4 {
  margin: 0 0 0.25rem 0;
  font-size: 0.9rem;
  font-weight: 500;
  color: #718096;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: #2d3748;
  margin: 0;
}

.progress-container {
  margin-top: 1rem;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: #e2e8f0;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  transition: width 0.6s ease;
  border-radius: 4px;
}

.progress-fill.cpu {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.progress-fill.memory {
  background: linear-gradient(135deg, #48bb78 0%, #38a169 100%);
}

.progress-fill.disk {
  background: linear-gradient(135deg, #ed8936 0%, #dd6b20 100%);
}

.process-table-container {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
}

.table-header {
  padding: 1.5rem;
  border-bottom: 1px solid #e2e8f0;
}

.table-header h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #2d3748;
}

.table-header p {
  margin: 0;
  color: #718096;
  font-size: 0.95rem;
}

.table-wrapper {
  overflow-x: auto;
}

.process-table {
  width: 100%;
  border-collapse: collapse;
}

.process-table th {
  background: #f8fafc;
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: #4a5568;
  border-bottom: 1px solid #e2e8f0;
  font-size: 0.9rem;
}

.process-table td {
  padding: 1rem;
  border-bottom: 1px solid #f1f5f9;
  color: #2d3748;
}

.process-row:hover {
  background: #f8fafc;
}

.process-pid {
  font-family: monospace;
  font-weight: 600;
  color: #667eea;
}

.name-info .name {
  font-weight: 600;
  color: #2d3748;
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

.status-badge.running {
  background: rgba(72, 187, 120, 0.1);
  color: #38a169;
  border: 1px solid rgba(72, 187, 120, 0.2);
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

@media (max-width: 768px) {
  .system-monitor {
    padding: 1rem;
  }
  
  .control-panel {
    flex-direction: column;
    align-items: stretch;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .process-table {
    font-size: 0.875rem;
  }
  
  .process-table th,
  .process-table td {
    padding: 0.75rem 0.5rem;
  }
}
</style>