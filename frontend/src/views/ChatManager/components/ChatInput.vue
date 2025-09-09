<template>
  <div class="input-area">
    <el-input
        v-model="inputValue"
        type="textarea"
        :rows="3"
        placeholder="输入消息..."
        :disabled="disabled"
        @keydown="handleKeydown"
    />
    <div style="margin-top: 10px; text-align: right">
      <el-button type="primary" @click="handleSend" :disabled="disabled">
        发送 (Enter)
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  modelValue: string
  disabled?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'send'): void
  (e: 'keydown', event: KeyboardEvent): void
}>()

const inputValue = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const handleSend = () => {
  emit('send')
}

const handleKeydown = (event: KeyboardEvent) => {
  emit('keydown', event)
}
</script>

<style scoped>
.input-area {
  padding: 24px;
  border-top: 1px solid #e0e0e0;
  background-color: #ffffff;
}
</style>