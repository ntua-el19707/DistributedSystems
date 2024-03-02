import { Injectable } from '@angular/core';
import { TransactionsMsgClientModule } from './transactions-msg-client.module';
import { TransactionMsgClientService } from './transaction-msg-client.service';
import { BehaviorSubject, Observable} from 'rxjs';
import { map } from 'rxjs/operators';
import { TransactionMsgRow, transactionMsgResponse } from '../../sharable';
import { FilterTransactionMsg } from '../filter-transaction-msg-toa-behavir-subject.service';

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
  filterSend(filter:FilterTransactionMsg){
    const  observable = this.FilterBehaviorSubject(filter , this.#sendBehaviorSubject)
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
  filterAll(filter:FilterTransactionMsg){
    const  observable =  this.FilterBehaviorSubject(filter , this.#allBehaviorSubject)
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
  filter(filter:FilterTransactionMsg){
    const  observable =  this.FilterBehaviorSubject(filter , this.#inboxBehaviorSubject)
    observable.subscribe((r:transactionMsgResponse)=>{
     this.#inboxBehaviorTransactionsSubject.next(r.transactions)
    })
  }
  getInboxTransactions():BehaviorSubject<Array<TransactionMsgRow>>{
    return this.#inboxBehaviorTransactionsSubject
  }
  FilterBehaviorSubject(filter:FilterTransactionMsg , subject: BehaviorSubject<transactionMsgResponse>):Observable<transactionMsgResponse>{
      const  observable =  subject.pipe(map((response:transactionMsgResponse )=>{
      const filteredTransactions = response.transactions.filter(transaction => {
      return (!filter.To || transaction.To <= filter.To) &&
             (!filter.From || transaction.From >= filter.From) &&
             (!filter.Message || transaction.Msg.includes(filter.Message))&& 
             (!filter.SendTimeLess || transaction.SendTime >= filter.SendTimeLess) &&
             (!filter.SendTimeMore || transaction.SendTime < filter.SendTimeMore);
    })
        return { ...response, transactions: filteredTransactions };
    }))
    return observable
  }
}