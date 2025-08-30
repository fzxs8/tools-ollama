<template>
  <div class="system-monitor">
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>Ollama服务状态</span>
            </div>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="服务状态">
              <el-tag type="success">运行中</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="版本">0.1.30</el-descriptions-item>
            <el-descriptions-item label="运行时间">2小时 15分钟</el-descriptions-item>
            <el-descriptions-item label="端口">11434</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>系统资源</span>
            </div>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="CPU使用率">15%</el-descriptions-item>
            <el-descriptions-item label="内存使用率">42%</el-descriptions-item>
            <el-descriptions-item label="磁盘使用率">68%</el-descriptions-item>
            <el-descriptions-item label="网络">上行: 1.2MB/s 下行: 0.8MB/s</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>模型加载状态</span>
            </div>
          </template>
          <el-table :data="loadedModels" style="width: 100%">
            <el-table-column prop="name" label="模型名称" />
            <el-table-column prop="status" label="状态">
              <template #default="scope">
                <el-tag :type="scope.row.status === '已加载' ? 'success' : 'info'">
                  {{ scope.row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="memory" label="内存占用" />
            <el-table-column prop="lastUsed" label="最后使用" />
            <el-table-column label="操作">
              <template #default="scope">
                <el-button 
                  size="small" 
                  @click="unloadModel(scope.row)"
                  :disabled="scope.row.status !== '已加载'"
                >
                  卸载
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>最近活动</span>
            </div>
          </template>
          <el-timeline>
            <el-timeline-item
              v-for="(activity, index) in activities"
              :key="index"
              :timestamp="activity.timestamp"
              :type="activity.type"
            >
              {{ activity.content }}
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface LoadedModel {
  name: string
  status: string
  memory: string
  lastUsed: string
}

interface Activity {
  content: string
  timestamp: string
  type?: '' | 'primary' | 'success' | 'warning' | 'danger' | 'info'
}

const loadedModels = ref<LoadedModel[]>([
  { name: 'llama3:8b', status: '已加载', memory: '4.2 GB', lastUsed: '2分钟前' },
  { name: 'mistral:7b', status: '未加载', memory: '0 GB', lastUsed: '1小时前' }
])

const activities = ref<Activity[]>([
  { content: '启动Ollama服务', timestamp: '2023-08-30 10:00:00', type: 'success' },
  { content: '下载模型 llama3:8b', timestamp: '2023-08-30 10:05:00', type: 'primary' },
  { content: '加载模型 llama3:8b', timestamp: '2023-08-30 10:10:00', type: 'primary' },
  { content: '用户开始对话', timestamp: '2023-08-30 10:15:00', type: 'info' },
  { content: '用户结束对话', timestamp: '2023-08-30 10:30:00', type: 'info' }
])

const unloadModel = (model: LoadedModel) => {
  model.status = '未加载'
  model.memory = '0 GB'
}
</script>

<style scoped>
.system-monitor {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>