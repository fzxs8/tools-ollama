import { defineStore } from 'pinia';
import i18n from '../i18n';

interface I18nState {
  locale: string;
}

// Initialize locale from localStorage
const savedLocale = localStorage.getItem('language') || 'en';
i18n.global.locale.value = savedLocale;

export const useI18nStore = defineStore('i18n', {
  state: (): I18nState => ({
    locale: savedLocale,
  }),

  actions: {
    setLocale(locale: string) {
      this.locale = locale;
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