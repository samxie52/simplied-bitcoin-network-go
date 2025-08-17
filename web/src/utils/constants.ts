export const API_BASE_URL = 'http://localhost:8080';
export const WS_URL = 'ws://localhost:8080/ws';

export const ROUTES = {
  DASHBOARD: '/',
  BLOCK_EXPLORER: '/blocks',
  MINING: '/mining',
  TRANSACTIONS: '/transactions',
  WALLET: '/wallet',
} as const;

export const COLORS = {
  PRIMARY: '#1890ff',
  SUCCESS: '#52c41a',
  WARNING: '#faad14',
  ERROR: '#f5222d',
  BITCOIN: '#f7931a',
} as const;