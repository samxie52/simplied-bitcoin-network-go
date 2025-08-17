import React from 'react';
import { ConfigProvider, Layout, theme } from 'antd';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import './styles/globals.css';

const { Header, Content, Footer } = Layout;

// 创建 QueryClient 实例
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 5 * 60 * 1000, // 5分钟
    },
  },
});

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
        <Layout className="min-h-screen">
          <Header className="bg-white shadow-sm border-b">
            <div className="flex items-center justify-between h-full">
              <div className="flex items-center space-x-4">
                <h1 className="text-xl font-bold text-gray-800">
                  🪙 简化版比特币网络
                </h1>
              </div>
              <div className="text-sm text-gray-500">
                Go + React 区块链演示
              </div>
            </div>
          </Header>
          
          <Content className="p-6 bg-gray-50">
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
                    <h3 className="font-semibold mb-2">⚡ Vite 5.0</h3>
                    <p className="text-sm text-gray-600">极速构建工具</p>
                  </div>
                  <div className="p-4 border rounded-lg">
                    <h3 className="font-semibold mb-2">🔷 TypeScript 5.3</h3>
                    <p className="text-sm text-gray-600">类型安全开发</p>
                  </div>
                  <div className="p-4 border rounded-lg">
                    <h3 className="font-semibold mb-2">🎨 Ant Design 5.0</h3>
                    <p className="text-sm text-gray-600">企业级UI组件</p>
                  </div>
                </div>
              </div>
            </div>
          </Content>
          
          <Footer className="text-center bg-white border-t">
            <p className="text-gray-500">
              简化版比特币网络 - Go语言区块链开发实践 © 2024
            </p>
          </Footer>
        </Layout>
      </ConfigProvider>
    </QueryClientProvider>
  );
};

export default App;