import { nodeDetails } from "./nodeDetails"

export interface  TransactionMsgRow{
    From :number 
    To: number  
    Nonce:number
    Msg: string 
    SendTime :number 
    TransactionId:string
}
export interface transactionMsgRequest {
    to :number 
    msg: string
}
export interface transactionMsgResponse {

    transactions: Array<TransactionMsgRow>
    nodeDetails:  nodeDetails
} 