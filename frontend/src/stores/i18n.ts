import { defineStore } from 'pinia';

interface I18nState {
  locale: string;
}

// Initialize locale from localStorage
const savedLocale = localStorage.getItem('language') || 'en';

export const useI18nStore = defineStore('i18n', {
  state: (): I18nState => ({
    locale: savedLocale,
  }),

  actions: {
    async setLocale(locale: string) {
      this.locale = locale;
      // Dynamically import i18n to avoid circular dependency
      const { default: i18n } = await import('../i18n');
      i18n.global.locale.value = locale;
      localStorage.setItem('language', locale);
    },

    toggleLocale() {
      const newLocale = this.locale === 'en' ? 'zh' : 'en';
      this.setLocale(newLocale);
    },
  },

  getters: {
    currentLocale: (state) => state.locale,
    isEnglish: (state) => state.locale === 'en',
    isChinese: (state) => state.locale === 'zh',
  },
});