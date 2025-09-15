# 国际化 (i18n) 功能说明

## 概述
本项目已添加中英文双语支持，用户可以在界面中切换语言。

## 文件结构
```
src/
├── i18n/
│   ├── index.ts              # i18n 配置文件
│   └── locales/
│       ├── en.json           # 英文翻译资源
│       └── zh.json           # 中文翻译资源
├── stores/
│   └── i18n.ts              # 语言状态管理
└── components/
    └── LanguageSwitcher.vue  # 语言切换组件
```

## 使用方法

### 1. 在组件中使用翻译
```vue
<template>
  <div>
    <h1>{{ t('nav.modelManager') }}</h1>
    <button>{{ t('common.save') }}</button>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>
```

### 2. 在 JavaScript 中使用翻译
```typescript
import i18n from '../i18n';

const { t } = i18n.global;
const message = t('messages.configSaved');
```

### 3. 语言切换
语言切换器组件已集成到侧边栏中，用户可以点击切换中英文。

## 翻译资源结构
翻译资源按功能模块组织：
- `common`: 通用词汇（保存、取消、删除等）
- `nav`: 导航菜单
- `promptPilot`: 提示词驾驶舱模块
- `modelManager`: 模型管理模块
- `serverSettings`: 服务设置模块
- `openaiAdapter`: OpenAI 适配器模块
- `messages`: 系统消息和提示

## 添加新翻译
1. 在 `src/i18n/locales/en.json` 中添加英文翻译
2. 在 `src/i18n/locales/zh.json` 中添加对应的中文翻译
3. 在组件中使用 `t('key.path')` 调用翻译

## 语言持久化
用户选择的语言会自动保存到 localStorage，下次打开应用时会记住用户的语言偏好。