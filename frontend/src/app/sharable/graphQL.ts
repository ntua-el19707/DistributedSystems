export interface NodeInfoGraphQL {
  nodeId?: string;
  indexId?: number;
  uri?: string;
  uriPublic?: string;
}
export  interface blockDto{
    index?: number; // Block Index in chain
    created_at?: number; // The time that block is created
    validator?: number; // Validator of block in index
    capacity?: number; // Capacity of block
    current_hash?: string; // Current hash of block
    parrent_hash?: string; // Parent hash of block

}
export  interface blockCoinDto{
  block?:blockDto 
  transactions?:Array<TransactionCoinsRowGraphQL>
}
export  interface blockMsgDto{
  block?:blockDto 
  transactions?:Array<TransactionMsgRowGraphQL>
}
export interface NodeInfoSelfGraphQL {
  client?: NodeInfoGraphQL;
  total?: number;
}

export interface TransactionMsgRowGraphQL {
  From?: number;
  To?: number;
  Nonce?: number;
  Msg?: string;
  Time?: number;
  TransactionId?: string;
}

export interface TransactionMsgListGraphQL {
  transactions: TransactionMsgRowGraphQL[];

}
export interface balance {
  availableBalance?:number
}
export interface TransactionMsgList {
  transactions: TransactionMsgRowGraphQL[];
  nodeDetails:NodeInfoGraphQL
}

export interface TransactionCoinsRowGraphQL {
  From?: number;
  To?: number;
  Coins?: number;
  Nonce?: number;
  Reason?: string;
  Time?: number;
  TransactionId?: string;
}

export interface TransactionCoinListGraphQL {
  Transactions?: TransactionCoinsRowGraphQL[];
}

export interface GraphQLData {
  clients?: NodeInfoGraphQL[];
  self?: NodeInfoSelfGraphQL;
  allNodes?: NodeInfoGraphQL[];
  inbox?: TransactionMsgListGraphQL;
  allTransactionMsg?: TransactionMsgListGraphQL;
  send?: TransactionMsgListGraphQL;
  getTransactionsCoins?: TransactionCoinListGraphQL;
  nodeTransactions?:TransactionCoinListGraphQL
  blockChainCoins?: Array<blockCoinDto>;
  blockChainMsg?: Array<blockMsgDto>;
  balance?:balance
}

export interface GraphQLResponse {
  data?: GraphQLData;
}
