export interface Block {
    index: number;
    timestamp: number;
    data: string;
    previousHash: string;
    hash: string;
    nonce: number;
    difficulty: number;
  }
  
  export interface Transaction {
    id: string;
    from: string;
    to: string;
    amount: number;
    timestamp: number;
    signature?: string;
  }
  
  export interface MiningStats {
    hashRate: number;
    difficulty: number;
    blocksFound: number;
    totalReward: number;
  }