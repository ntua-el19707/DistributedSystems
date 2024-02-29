import { Injectable } from '@angular/core';
import { TransactionsMsgClientModule } from './transactions-msg-client.module';
import { TransactionMsgClientService } from './transaction-msg-client.service';
import { BehaviorSubject, Observable} from 'rxjs';
import { map } from 'rxjs/operators';
import { TransactionMsgRow, transactionMsgResponse } from '../../sharable';

@Injectable({
  providedIn: TransactionsMsgClientModule
})
export class TransactionMsgBehaviorService {
  #inboxBehaviorSubject:BehaviorSubject<transactionMsgResponse> = new BehaviorSubject<transactionMsgResponse>({nodeDetails:{indexId:0 , nodeId:"" , uri:"" ,uriPublic:""},transactions:[]})
  #inboxBehaviorTransactionsSubject:BehaviorSubject<Array<TransactionMsgRow>> = new BehaviorSubject<Array<TransactionMsgRow>>([])
  #sendBehaviorSubject:BehaviorSubject<transactionMsgResponse> = new BehaviorSubject<transactionMsgResponse>({nodeDetails:{indexId:0 , nodeId:"" , uri:"" ,uriPublic:""},transactions:[]})
  #sendBehaviorTransactionsSubject:BehaviorSubject<Array<TransactionMsgRow>> = new BehaviorSubject<Array<TransactionMsgRow>>([])
  #allBehaviorSubject:BehaviorSubject<transactionMsgResponse> = new BehaviorSubject<transactionMsgResponse>({nodeDetails:{indexId:0 , nodeId:"" , uri:"", uriPublic:""},transactions:[]})
  #allBehaviorTransactionsSubject:BehaviorSubject<Array<TransactionMsgRow>> = new BehaviorSubject<Array<TransactionMsgRow>>([])
  constructor(private transactionMsgClientService :TransactionMsgClientService) { }
  fetchSend() {
    this.transactionMsgClientService.getSend().subscribe(r =>{
    
      this.#sendBehaviorSubject.next(r)
      this.#sendBehaviorTransactionsSubject.next(r.transactions)
    } , err=>{
      console.log(err)
    } , ()=>{})
  }
 getSend():BehaviorSubject<transactionMsgResponse> {
    return this.#sendBehaviorSubject
  }
  filterSend(msg:string){
    const  observable =  this.#sendBehaviorSubject.pipe(map((response:transactionMsgResponse )=>{
      const filteredTransactions = response.transactions.filter(transaction => transaction.Msg.includes(msg));
        return { ...response, transactions: filteredTransactions };
    }))
    observable.subscribe((r:transactionMsgResponse)=>{
     this.#sendBehaviorTransactionsSubject.next(r.transactions)
    })
  }
  getSendTransactions():BehaviorSubject<Array<TransactionMsgRow>>{
    return this.#sendBehaviorTransactionsSubject
  }
   fetchAll() {
    this.transactionMsgClientService.getAll().subscribe(r =>{
    
      this.#allBehaviorSubject.next(r)
      this.#allBehaviorTransactionsSubject.next(r.transactions)
    } , err=>{
      console.log(err)
    } , ()=>{})
  }
 getAll():BehaviorSubject<transactionMsgResponse> {
    return this.#allBehaviorSubject
  }
  filterAll(msg:string){
    const  observable =  this.#allBehaviorSubject.pipe(map((response:transactionMsgResponse )=>{
      const filteredTransactions = response.transactions.filter(transaction => transaction.Msg.includes(msg));
        return { ...response, transactions: filteredTransactions };
    }))
    observable.subscribe((r:transactionMsgResponse)=>{
     this.#allBehaviorTransactionsSubject.next(r.transactions)
    })
  }
  getAllTransactions():BehaviorSubject<Array<TransactionMsgRow>>{
    return this.#allBehaviorTransactionsSubject
  }
  fetchInbox() {
    this.transactionMsgClientService.getINBOX().subscribe(r =>{
     
      this.#inboxBehaviorSubject.next(r)
      this.#inboxBehaviorTransactionsSubject.next(r.transactions)
    } , err=>{
      console.log(err)
    } , ()=>{})
  }
 getInbox():BehaviorSubject<transactionMsgResponse> {
    return this.#inboxBehaviorSubject
  }
  filter(msg:string){
    const  observable =  this.#inboxBehaviorSubject.pipe(map((response:transactionMsgResponse )=>{
      const filteredTransactions = response.transactions.filter(transaction => transaction.Msg.includes(msg));
        return { ...response, transactions: filteredTransactions };
    }))
    observable.subscribe((r:transactionMsgResponse)=>{
     this.#inboxBehaviorTransactionsSubject.next(r.transactions)
    })
  }
  getInboxTransactions():BehaviorSubject<Array<TransactionMsgRow>>{
    return this.#inboxBehaviorTransactionsSubject
  }
}