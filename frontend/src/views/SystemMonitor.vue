<template>
  <div class="system-monitor">
    <el-card class="monitor-card">
      <template #header>
        <div class="card-header">
          <span>系统监控</span>
          <el-button @click="refreshData">刷新</el-button>
        </div>
      </template>
      
      <div class="monitor-content">
        <el-row :gutter="20" style="margin-bottom: 20px;">
          <el-col :span="8">
            <el-card class="stat-card">
              <div class="stat-title">CPU 使用率</div>
              <div class="stat-value">{{ cpuUsage }}%</div>
              <el-progress :percentage="cpuUsage" :stroke-width="10" />
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card class="stat-card">
              <div class="stat-title">内存使用率</div>
              <div class="stat-value">{{ memoryUsage }}%</div>
              <el-progress :percentage="memoryUsage" :stroke-width="10" status="success" />
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card class="stat-card">
              <div class="stat-title">磁盘使用率</div>
              <div class="stat-value">{{ diskUsage }}%</div>
              <el-progress :percentage="diskUsage" :stroke-width="10" status="warning" />
            </el-card>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="24">
            <el-card>
              <template #header>
                <div class="card-header">
                  <span>Ollama 服务状态</span>
                </div>
              </template>
              <el-table :data="ollamaProcesses" style="width: 100%" empty-text="暂无数据">
                <el-table-column prop="pid" label="进程ID" />
                <el-table-column prop="name" label="进程名称" />
                <el-table-column prop="cpu_percent" label="CPU%" />
                <el-table-column prop="memory_percent" label="内存%" />
                <el-table-column prop="status" label="状态" />
              </el-table>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import {onMounted, ref} from 'vue'

const cpuUsage = ref(0)
const memoryUsage = ref(0)
const diskUsage = ref(0)
const ollamaProcesses = ref<any[]>([])


const refreshData = () => {
  // 模拟刷新数据
  cpuUsage.value = Math.floor(Math.random() * 100)
  memoryUsage.value = Math.floor(Math.random() * 100)
  diskUsage.value = Math.floor(Math.random() * 100)
  
  // 模拟进程数据
  ollamaProcesses.value = [
    { pid: 1234, name: 'ollama', cpu_percent: '25.6%', memory_percent: '15.2%', status: '运行中' },
    { pid: 5678, name: 'ollama-run', cpu_percent: '10.2%', memory_percent: '8.7%', status: '运行中' }
  ]
}

onMounted(() => {
  refreshData()
})
</script>

<style scoped>
.system-monitor {
  padding: 20px;
  height: 100%;
  box-sizing: border-box;
}

.monitor-card {
  height: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.monitor-card :deep(.el-card__body) {
  height: calc(100% - 60px);
  overflow-y: auto;
}

.monitor-content {
  height: 100%;
}

.stat-card {
  text-align: center;
}

.stat-title {
  font-size: 16px;
  color: #606266;
  margin-bottom: 10px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 10px;
}

</style>