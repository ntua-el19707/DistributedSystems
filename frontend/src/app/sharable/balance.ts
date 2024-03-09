import { NodeInfoGraphQL, TransactionCoinsRowGraphQL } from "./graphQL"
import { nodeDetails } from "./nodeDetails"

export interface BalanceRsp {
    availableBalance :number
}
export interface  TransactionCoinsRow{
    From :number 
    To: number  
    Coins: number 
    Nonce:number
    Reason: string 
    SendTime :number 
    TransactionId:string
}
export interface transactionCoinRequest {
     to:number,
     amount:number
}
export interface transactionCoinResponse {

    transactions: Array<TransactionCoinsRow>
    nodeDetails:  nodeDetails
} 
export interface TransactionCoinsAndNodeDetails {
  transactions: Array<TransactionCoinsRowGraphQL>;
  nodeDetails: NodeInfoGraphQL;
} 