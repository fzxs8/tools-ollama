<template>
  <div class="language-switcher">
    <button @click="toggleLanguage" class="lang-btn" :title="buttonTitle">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
        <path d="M12.87 15.07L10.33 12.56L10.36 12.53C12.1 10.59 13.34 8.36 14.07 6H17V4H10V2H8V4H1V6H12.17C11.5 7.92 10.44 9.75 9 11.35C8.07 10.32 7.3 9.19 6.69 8H4.69C5.42 9.63 6.42 11.17 7.67 12.56L2.58 17.58L4 19L9 14L12.11 17.11L12.87 15.07ZM18.5 10H16.5L12 22H14L15.12 19H19.87L21 22H23L18.5 10ZM15.88 17L17.5 12.67L19.12 17H15.88Z" stroke="currentColor" stroke-width="0.5" fill="currentColor"/>
      </svg>
      <span>{{ currentLanguageLabel }}</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useI18nStore } from '../stores/i18n';

const i18nStore = useI18nStore();

const currentLanguageLabel = computed(() => {
  return i18nStore.isEnglish ? 'EN' : '中';
});

const buttonTitle = computed(() => {
  return i18nStore.isEnglish ? 'Switch to Chinese' : '切换到英文';
});

const toggleLanguage = () => {
  i18nStore.toggleLocale();
};
</script>

<style scoped>
.language-switcher {
  display: flex;
  align-items: center;
}

.lang-btn {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.5rem;
  background: transparent;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 6px;
  color: #718096;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 0.75rem;
  font-weight: 400;
}

.lang-btn:hover {
  background: rgba(0, 0, 0, 0.05);
  border-color: rgba(0, 0, 0, 0.2);
  color: #4a5568;
}

.lang-btn svg {
  opacity: 0.6;
  width: 12px;
  height: 12px;
}

.lang-btn span {
  min-width: 16px;
  text-align: center;
  font-size: 0.7rem;
}
</style>