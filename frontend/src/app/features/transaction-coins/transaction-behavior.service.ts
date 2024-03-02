import { Injectable } from '@angular/core';
import { TransactionCoinsModule } from './transaction-coins.module';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, map} from 'rxjs';
import { BalanceRsp, transactionCoinResponse } from '../../sharable';
import { TransactionClientService } from './transaction-client.service';
import { filterTranctionCoinSubject } from '../filter-transaction-coin-toa-behavir-subject.service';

@Injectable({
  providedIn: TransactionCoinsModule
})
export class TransactionBehaviorService {
  #balanceBehaviorSubject : BehaviorSubject<number> = new BehaviorSubject<number>(0.0)
  #myTransactionsBehaviorSubject:BehaviorSubject<transactionCoinResponse> = new BehaviorSubject<transactionCoinResponse>({nodeDetails:{indexId:0 , nodeId:"" , uri:"" ,uriPublic:""},transactions:[]})
  #myTransactionsfilterBehaviorSubject:BehaviorSubject<transactionCoinResponse> = new BehaviorSubject<transactionCoinResponse>({nodeDetails:{indexId:0 , nodeId:"" , uri:"",uriPublic:""},transactions:[]})
 
   #allTransactionsBehaviorSubject:BehaviorSubject<transactionCoinResponse> = new BehaviorSubject<transactionCoinResponse>({nodeDetails:{indexId:0 , nodeId:"" , uri:"",uriPublic:""},transactions:[]})
  #allTransactionsfilterBehaviorSubject:BehaviorSubject<transactionCoinResponse> = new BehaviorSubject<transactionCoinResponse>({nodeDetails:{indexId:0 , nodeId:"" , uri:"",uriPublic:""},transactions:[]})
 
  constructor(private transactionClientService:TransactionClientService){  }
  fetchBalance(){
this.transactionClientService.getBalance().subscribe((r:BalanceRsp)=>{
  console.log(r)    
  this.#balanceBehaviorSubject.next(r.availableBalance)
    } ,err=>{} , ()=>{})
  }
  getBalanceSubject():BehaviorSubject<number>{
    return this.#balanceBehaviorSubject
  }

  fetchMyTransactions() {
    this.transactionClientService.getMyTransactions().subscribe(r =>{
      this.#myTransactionsBehaviorSubject.next(r)
      this.#myTransactionsfilterBehaviorSubject.next(r)
    } , err=>{
      console.log(err)
    } , ()=>{})
  }
getMyFilterTransactions():BehaviorSubject<transactionCoinResponse> {
    return this.#myTransactionsfilterBehaviorSubject
  }
  filterCoins(filter:filterTranctionCoinSubject){
    const  observable =  this.#myTransactionsBehaviorSubject.pipe(map((response:transactionCoinResponse )=>{
      const filteredTransactions = response.transactions.filter(transaction => {
      return (!filter.To || transaction.To <= filter.To) &&
            (!filter.Reason || transaction.Reason === filter.Reason)&& 
            (!filter.From || transaction.From >= filter.From) &&
             (!filter.CoinsMin || transaction.Coins >= filter.CoinsMin) &&
             (!filter.CoinsMax || transaction.Coins <= filter.CoinsMax) &&
             (!filter.SendTimeLess || transaction.SendTime >= filter.SendTimeLess) &&
             (!filter.SendTimeMore || transaction.SendTime < filter.SendTimeMore);
    })
  
        return { ...response, transactions: filteredTransactions };
    }))
    observable.subscribe((r:transactionCoinResponse)=>{
     this.#myTransactionsfilterBehaviorSubject.next(r)
    })
  }
  getMyTransactions():BehaviorSubject<transactionCoinResponse>{
    return this.#myTransactionsBehaviorSubject
  }
getMyAllFilterTransactions():BehaviorSubject<transactionCoinResponse>{
return this.#allTransactionsfilterBehaviorSubject
}
  fetchAllTransactions() {
    this.transactionClientService.getAllTransactions().subscribe(r =>{
      this.#allTransactionsBehaviorSubject.next(r)
      this.#allTransactionsfilterBehaviorSubject.next(r)
    } , err=>{
      console.log(err)
    } , ()=>{})
  }
  filterAllCoins(filter:filterTranctionCoinSubject){
    const  observable =  this.#allTransactionsBehaviorSubject.pipe(map((response:transactionCoinResponse )=>{
      const filteredTransactions = response.transactions.filter(transaction => {
      return (!filter.To || transaction.To <= filter.To) &&
             (!filter.From || transaction.From >= filter.From) &&
             (!filter.Reason || transaction.Reason === filter.Reason)&& 
             (!filter.CoinsMin || transaction.Coins >= filter.CoinsMin) &&
             (!filter.CoinsMax || transaction.Coins <= filter.CoinsMax) &&
             (!filter.SendTimeLess || transaction.SendTime >= filter.SendTimeLess) &&
             (!filter.SendTimeMore || transaction.SendTime < filter.SendTimeMore);
    })
        return { ...response, transactions: filteredTransactions };
    }))
    observable.subscribe((r:transactionCoinResponse)=>{
     this.#allTransactionsfilterBehaviorSubject.next(r)
    })
  }
  }
 

