import React from 'react';
import { ConfigProvider, Layout, theme } from 'antd';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
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
  );
};

export default App;