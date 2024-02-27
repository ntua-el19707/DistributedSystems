import { Injectable } from '@angular/core';
import { TransactionCoinsModule } from './transaction-coins.module';
import { HttpClient } from '@angular/common/http';
import { BalanceRsp, transactionCoinRequest, transactionCoinResponse } from '../../sharable';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: TransactionCoinsModule
})
export class TransactionClientService {

 constructor(private http:HttpClient) {

  }
  getBalance():Observable<BalanceRsp>{
    return this.http.get("/api/v1/balance") as  Observable<BalanceRsp>
  }
  getMyTransactions():Observable<transactionCoinResponse>{
    return this.http.get("/api/v1/transactions") as  Observable<transactionCoinResponse>
  }
  postTransaction(to:number , coins:number){
    const body:transactionCoinRequest = {
      to:to , 
      amount:coins
    }
    this.http.post("/api/v1/transfer" , body).subscribe((r)=>{} ,err=>{console.log(err)} ,()=>{} )


  }
}
