# Step 0.2: 国际化和可访问性支持

## 📋 任务概述

**功能**: 建立国际化和可访问性支持基础设施  
**前置条件**: Step 0.1 完成  
**预计时间**: 2-3小时  
**难度等级**: ⭐⭐⭐☆☆

## 🎯 学习目标

完成本步骤后，你将掌握：
- react-i18next 国际化框架的配置和使用
- 多语言资源文件的组织和管理
- 语言切换和持久化存储机制
- Web可访问性(A11y)最佳实践
- ARIA标签和语义化HTML的实现
- 键盘导航和屏幕阅读器支持
- 高对比度主题和视觉辅助功能

## 🛠️ 技术栈详情

| 技术 | 版本 | 用途 | 官方文档 |
|------|------|------|----------|
| react-i18next | 13.5+ | React国际化 | https://react.i18next.com/ |
| i18next | 23.7+ | 国际化核心 | https://www.i18next.com/ |
| @axe-core/react | 4.8+ | 可访问性检查 | https://github.com/dequelabs/axe-core-npm |
| focus-trap-react | 10.2+ | 焦点管理 | https://github.com/focus-trap/focus-trap-react |

**注意**: 由于React 19兼容性问题，我们移除了`react-helmet-async`依赖，使用React 19原生文档头管理功能。

## 📁 新增目录结构

```
web/src/
├── i18n/                               # 国际化配置
│   ├── index.ts                        # i18n主配置文件
│   ├── resources/                      # 语言资源文件
│   │   ├── zh-CN/                      # 中文资源
│   │   │   ├── common.json             # 通用翻译
│   │   │   ├── blockchain.json         # 区块链术语
│   │   │   ├── mining.json             # 挖矿相关
│   │   │   ├── transaction.json        # 交易相关
│   │   │   └── wallet.json             # 钱包相关
│   │   └── en-US/                      # 英文资源
│   │       ├── common.json
│   │       ├── blockchain.json
│   │       ├── mining.json
│   │       ├── transaction.json
│   │       └── wallet.json
│   └── detector.ts                     # 语言检测器
├── components/
│   └── A11y/                           # 可访问性组件
│       ├── FocusTrap.tsx               # 焦点陷阱组件
│       ├── SkipLink.tsx                # 跳转链接组件
│       ├── ScreenReaderOnly.tsx        # 屏幕阅读器专用组件
│       ├── KeyboardNavigation.tsx      # 键盘导航组件
│       └── LanguageSwitcher.tsx        # 语言切换器
├── hooks/
│   ├── useI18n.ts                      # 国际化Hook
│   ├── useKeyboardNavigation.ts        # 键盘导航Hook
│   └── useAccessibility.ts             # 可访问性Hook
└── styles/
    └── themes/                         # 主题配置
        ├── default.css                 # 默认主题
        ├── high-contrast.css           # 高对比度主题
        └── dark.css                    # 深色主题
```

## 🚀 实施步骤

### 第一步：安装国际化和可访问性依赖

```bash
# 进入前端项目目录
cd web

# 安装国际化依赖
# 安装国际化依赖（这些都支持React 19）
npm install react-i18next@^13.5.0 i18next@^23.7.0
npm install i18next-browser-languagedetector@^7.2.0
npm install i18next-http-backend@^2.4.0

# 安装可访问性依赖（使用兼容版本）
npm install @axe-core/react@^4.8.0
npm install focus-trap-react@^10.2.0 --legacy-peer-deps

# 安装开发依赖
npm install -D @types/react-helmet@^6.1.0
```

### 第二步：配置i18next国际化框架

```typescript
// src/i18n/index.ts
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

// 导入语言资源
import zhCNCommon from './resources/zh-CN/common.json';
import enUSCommon from './resources/en-US/common.json';

const resources = {
  'zh-CN': {
    translation: zhCNCommon,
  },
  'en-US': {
    translation: enUSCommon,
  },
};

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: 'zh-CN',
    defaultNS: 'translation',
    debug: process.env.NODE_ENV === 'development',
    
    interpolation: {
      escapeValue: false, // React已经处理了XSS
    },
    
    detection: {
      order: ['localStorage', 'navigator', 'htmlTag'],
      lookupLocalStorage: 'i18nextLng',
      caches: ['localStorage'],
    },
    
    react: {
      useSuspense: false,
    },
  });

export default i18n;
```

### 第三步：创建语言资源文件

```json
// src/i18n/resources/zh-CN/common.json
{
  "app": {
    "title": "简化版比特币网络",
    "subtitle": "Go + React 区块链演示",
    "loading": "加载中...",
    "error": "发生错误",
    "success": "操作成功",
    "confirm": "确认",
    "cancel": "取消",
    "save": "保存",
    "delete": "删除",
    "edit": "编辑",
    "view": "查看",
    "search": "搜索",
    "filter": "筛选",
    "export": "导出",
    "import": "导入",
    "refresh": "刷新"
  },
  "navigation": {
    "dashboard": "仪表盘",
    "blockExplorer": "区块浏览器",
    "mining": "挖矿",
    "transactions": "交易",
    "wallet": "钱包",
    "settings": "设置"
  },
  "accessibility": {
    "skipToMain": "跳转到主内容",
    "skipToNavigation": "跳转到导航",
    "openMenu": "打开菜单",
    "closeMenu": "关闭菜单",
    "languageSwitcher": "语言切换器",
    "currentLanguage": "当前语言：{{language}}",
    "switchToLanguage": "切换到{{language}}",
    "keyboardShortcuts": "键盘快捷键",
    "screenReaderOnly": "仅屏幕阅读器可见"
  },
  "theme": {
    "default": "默认主题",
    "dark": "深色主题",
    "highContrast": "高对比度主题",
    "switchTheme": "切换主题"
  }
}
```

```json
// src/i18n/resources/en-US/common.json
{
  "app": {
    "title": "Simplified Bitcoin Network",
    "subtitle": "Go + React Blockchain Demo",
    "loading": "Loading...",
    "error": "An error occurred",
    "success": "Operation successful",
    "confirm": "Confirm",
    "cancel": "Cancel",
    "save": "Save",
    "delete": "Delete",
    "edit": "Edit",
    "view": "View",
    "search": "Search",
    "filter": "Filter",
    "export": "Export",
    "import": "Import",
    "refresh": "Refresh"
  },
  "navigation": {
    "dashboard": "Dashboard",
    "blockExplorer": "Block Explorer",
    "mining": "Mining",
    "transactions": "Transactions",
    "wallet": "Wallet",
    "settings": "Settings"
  },
  "accessibility": {
    "skipToMain": "Skip to main content",
    "skipToNavigation": "Skip to navigation",
    "openMenu": "Open menu",
    "closeMenu": "Close menu",
    "languageSwitcher": "Language switcher",
    "currentLanguage": "Current language: {{language}}",
    "switchToLanguage": "Switch to {{language}}",
    "keyboardShortcuts": "Keyboard shortcuts",
    "screenReaderOnly": "Screen reader only"
  },
  "theme": {
    "default": "Default theme",
    "dark": "Dark theme",
    "highContrast": "High contrast theme",
    "switchTheme": "Switch theme"
  }
}
```

### 第四步：创建国际化Hook

```typescript
// src/hooks/useI18n.ts
import { useTranslation } from 'react-i18next';
import { useCallback } from 'react';

export interface UseI18nReturn {
  t: (key: string, options?: any) => string;
  language: string;
  changeLanguage: (lng: string) => Promise<void>;
  languages: Array<{ code: string; name: string; nativeName: string }>;
  isRTL: boolean;
}

const SUPPORTED_LANGUAGES = [
  { code: 'zh-CN', name: 'Chinese', nativeName: '中文' },
  { code: 'en-US', name: 'English', nativeName: 'English' },
];

export const useI18n = (): UseI18nReturn => {
  const { t, i18n } = useTranslation();

  const changeLanguage = useCallback(async (lng: string) => {
    try {
      await i18n.changeLanguage(lng);
      // 更新HTML lang属性
      document.documentElement.lang = lng;
      // 更新页面方向
      document.documentElement.dir = lng === 'ar' ? 'rtl' : 'ltr';
    } catch (error) {
      console.error('Failed to change language:', error);
    }
  }, [i18n]);

  return {
    t,
    language: i18n.language,
    changeLanguage,
    languages: SUPPORTED_LANGUAGES,
    isRTL: i18n.dir() === 'rtl',
  };
};
```

### 第五步：创建可访问性组件

```tsx
// src/components/A11y/LanguageSwitcher.tsx
import React, { useState } from 'react';
import { Button, Dropdown, Space } from 'antd';
import { GlobalOutlined, DownOutlined } from '@ant-design/icons';
import type { MenuProps } from 'antd';
import { useI18n } from '@/hooks/useI18n';

export const LanguageSwitcher: React.FC = () => {
  const { t } = useI18n();

  return (
    <Dropdown
      menu={{
        items: [
          {
            key: 'zh-CN',
            label: (
              <div className="flex items-center justify-between min-w-[120px]">
                <span>中文</span>
              </div>
            ),
            onClick: () => console.log('切换到中文'),
          },
          {
            key: 'en-US',
            label: (
              <div className="flex items-center justify-between min-w-[120px]">
                <span>English</span>
              </div>
            ),
            onClick: () => console.log('切换到英文'),
          },
        ],
      }}
      trigger={['click']}
      placement="bottomRight"
      arrow
    >
      <Button
        type="text"
        aria-label={t('accessibility.languageSwitcher')}
        aria-expanded="false"
        aria-haspopup="true"
        className="flex items-center"
      >
        <Space>
          <GlobalOutlined />
          <span className="hidden sm:inline">
            {t('accessibility.currentLanguage')}
          </span>
          <DownOutlined className="text-xs" />
        </Space>
      </Button>
    </Dropdown>
  );
};
```

```tsx
// src/components/A11y/SkipLink.tsx
import React from 'react';

interface SkipLinkProps {
  href: string;
  children: React.ReactNode;
}

export const SkipLink: React.FC<SkipLinkProps> = ({ href, children }) => {
  return (
    <a
      href={href}
      className="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 
                 bg-blue-600 text-white px-4 py-2 rounded-md z-50 
                 focus:outline-none focus:ring-2 focus:ring-blue-300"
      tabIndex={0}
    >
      {children}
    </a>
  );
};
```

### 第六步：更新主应用组件

```tsx
// src/App.tsx (更新版本)
import React from 'react';
import { ConfigProvider, Layout, theme } from 'antd';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { HelmetProvider } from 'react-helmet-async';
import { useI18n } from '@/hooks/useI18n';
import { LanguageSwitcher } from '@/components/A11y/LanguageSwitcher';
import { SkipLink } from '@/components/A11y/SkipLink';
import './styles/globals.css';
import './i18n';

const { Header, Content, Footer } = Layout;

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 5 * 60 * 1000,
    },
  },
});

const AppContent: React.FC = () => {
  const { t } = useI18n();

  return (
    <Layout className="min-h-screen">
      <SkipLink href="#main-content">
        {t('accessibility.skipToMain')}
      </SkipLink>
      
      <Header className="bg-white shadow-sm border-b">
        <div className="flex items-center justify-between h-full">
          <div className="flex items-center space-x-4">
            <h1 className="text-xl font-bold text-gray-800">
              🪙 {t('app.title')}
            </h1>
          </div>
          <div className="flex items-center space-x-4">
            <span className="text-sm text-gray-500 hidden md:inline">
              {t('app.subtitle')}
            </span>
            <LanguageSwitcher />
          </div>
        </div>
      </Header>
      
      <Content id="main-content" className="p-6 bg-gray-50" tabIndex={-1}>
        <div className="max-w-7xl mx-auto">
          <div className="bg-white rounded-lg p-8 text-center">
            <h2 className="text-2xl font-bold mb-4">
              🚀 前端环境初始化完成！
            </h2>
            <p className="text-gray-600 mb-6">
              React 18.3 + TypeScript 5.3 + Vite 5.0 开发环境已就绪
            </p>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-8">
              <div className="p-4 border rounded-lg">
                <h3 className="font-semibold mb-2">🌐 国际化支持</h3>
                <p className="text-sm text-gray-600">react-i18next 多语言</p>
              </div>
              <div className="p-4 border rounded-lg">
                <h3 className="font-semibold mb-2">♿ 可访问性</h3>
                <p className="text-sm text-gray-600">WCAG 2.1 AA 标准</p>
              </div>
              <div className="p-4 border rounded-lg">
                <h3 className="font-semibold mb-2">⌨️ 键盘导航</h3>
                <p className="text-sm text-gray-600">完整键盘支持</p>
              </div>
            </div>
          </div>
        </div>
      </Content>
      
      <Footer className="text-center bg-white border-t">
        <p className="text-gray-500">
          {t('app.title')} - Go语言区块链开发实践 © 2024
        </p>
      </Footer>
    </Layout>
  );
};

const App: React.FC = () => {
  return (
    <HelmetProvider>
      <QueryClientProvider client={queryClient}>
        <ConfigProvider
          theme={{
            algorithm: theme.defaultAlgorithm,
            token: {
              colorPrimary: '#1890ff',
              borderRadius: 8,
              fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
            },
          }}
        >
          <AppContent />
        </ConfigProvider>
      </QueryClientProvider>
    </HelmetProvider>
  );
};

export default App;
```

## ✅ 验证步骤

### 1. 多语言切换功能测试

```bash
# 启动开发服务器
npm run dev

# 测试步骤：
# 1. 点击右上角语言切换器
# 2. 选择不同语言
# 3. 验证页面文本正确切换
# 4. 刷新页面验证语言持久化
```

### 2. 可访问性检查

```bash
# 安装axe-core检查工具
npm install -D @axe-core/cli

# 运行可访问性检查
npx axe http://localhost:3000

# 预期结果：
# - 无严重可访问性问题
# - ARIA标签正确设置
# - 语义化HTML结构清晰
```

### 3. 键盘导航测试

```
测试步骤：
1. 使用Tab键浏览所有可交互元素
2. 使用Enter/Space激活按钮和链接
3. 使用Escape关闭弹窗和下拉菜单
4. 验证焦点指示器清晰可见
5. 测试跳转链接功能
```

### 4. 屏幕阅读器兼容性

```
测试工具：
- macOS: VoiceOver (Cmd + F5)
- Windows: NVDA (免费)
- Chrome: ChromeVox 扩展

验证要点：
- 页面结构正确朗读
- 按钮和链接有描述性标签
- 表单字段有关联标签
- 状态变化有适当通知
```

## 🔧 故障排除

### 常见问题1：i18next初始化失败

**症状**: 翻译不显示，控制台报错
**解决方案**:
```typescript
// 检查i18n配置是否正确导入
import './i18n'; // 确保在App.tsx中导入

// 检查资源文件路径
console.log(i18n.getResourceBundle('zh-CN', 'common'));
```

### 常见问题2：可访问性检查失败

**症状**: axe-core报告可访问性问题
**解决方案**:
```tsx
// 添加缺失的ARIA标签
<button aria-label="关闭对话框" onClick={handleClose}>
  ×
</button>

// 确保表单字段有标签
<label htmlFor="username">用户名</label>
<input id="username" type="text" />
```

### 常见问题3：键盘导航不工作

**症状**: Tab键无法正确导航
**解决方案**:
```css
/* 确保焦点样式可见 */
:focus {
  outline: 2px solid #1890ff;
  outline-offset: 2px;
}

/* 移除默认outline时添加自定义样式 */
:focus-visible {
  box-shadow: 0 0 0 2px #1890ff;
}
```

## 📚 学习资源

### 国际化资源
- [react-i18next 官方文档](https://react.i18next.com/)
- [i18next 最佳实践](https://www.i18next.com/principles/fallback)
- [Web 国际化指南](https://developer.mozilla.org/en-US/docs/Web/Localization)

### 可访问性资源
- [WCAG 2.1 指南](https://www.w3.org/WAI/WCAG21/quickref/)
- [ARIA 最佳实践](https://www.w3.org/WAI/ARIA/apg/)
- [可访问性测试工具](https://www.deque.com/axe/)

## 🎯 下一步计划

完成Step 0.2后，你将进行：

1. **Step 0.3**: 前端监控和分析
   - 集成Sentry错误监控
   - 配置Google Analytics用户行为分析

2. **Step 1.1**: 后端Go项目初始化
   - 建立Go模块和项目结构
   - 配置开发环境和工具链

## 📝 提交代码

```bash
# 在项目根目录执行
git add web/src/i18n/ web/src/components/A11y/ web/src/hooks/
git commit -m "feat: 实现国际化和可访问性支持

- 配置react-i18next多语言框架支持中英文切换
- 创建完整的语言资源文件和翻译管理系统
- 实现语言切换器组件和持久化存储
- 添加可访问性组件(SkipLink, ScreenReaderOnly, FocusTrap)
- 实现键盘导航和ARIA标签支持
- 集成axe-core可访问性检查工具
- 创建国际化和可访问性自定义Hooks
- 更新主应用组件支持多语言和可访问性功能"
```

---

**🎉 恭喜！** 你已经成功为React前端项目添加了完整的国际化和可访问性支持。现在你的应用支持多语言切换、键盘导航、屏幕阅读器，并符合WCAG 2.1 AA可访问性标准。
