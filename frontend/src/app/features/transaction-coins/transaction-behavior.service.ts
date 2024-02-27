import { Injectable } from '@angular/core';
import { TransactionCoinsModule } from './transaction-coins.module';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, map} from 'rxjs';
import { BalanceRsp, transactionCoinResponse } from '../../sharable';
import { TransactionClientService } from './transaction-client.service';

@Injectable({
  providedIn: TransactionCoinsModule
})
export class TransactionBehaviorService {
  #balanceBehaviorSubject : BehaviorSubject<number> = new BehaviorSubject<number>(0.0)
  #myTransactionsBehaviorSubject:BehaviorSubject<transactionCoinResponse> = new BehaviorSubject<transactionCoinResponse>({nodeDetails:{indexId:0 , nodeId:"" , uri:""},transactions:[]})
  #myTransactionsfilterBehaviorSubject:BehaviorSubject<transactionCoinResponse> = new BehaviorSubject<transactionCoinResponse>({nodeDetails:{indexId:0 , nodeId:"" , uri:""},transactions:[]})
 
  
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
  filterCoins(coin:number){
    const  observable =  this.#myTransactionsBehaviorSubject.pipe(map((response:transactionCoinResponse )=>{
      const filteredTransactions = response.transactions.filter(transaction => transaction.Coins >= coin);
        return { ...response, transactions: filteredTransactions };
    }))
    observable.subscribe((r:transactionCoinResponse)=>{
     this.#myTransactionsfilterBehaviorSubject.next(r)
    })
  }
  getMyTransactions():BehaviorSubject<transactionCoinResponse>{
    return this.#myTransactionsBehaviorSubject
  }

  }
 

