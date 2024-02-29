import { Injectable } from '@angular/core';
import { TransactionCoinsModule } from './transaction-coins.module';
import { HttpClient } from '@angular/common/http';
import { BalanceRsp, transactionCoinRequest, transactionCoinResponse } from '../../sharable';
import { BehaviorSubject, Observable } from 'rxjs';
import { OpenErrordialogService } from '../open-errordialog.service';
import { NavigateService } from '../navigate.service';

@Injectable({
  providedIn: TransactionCoinsModule
})
export class TransactionClientService {

 constructor(private http:HttpClient,private  navigateService:NavigateService, private openErrordialogService:OpenErrordialogService) {

  }
  getBalance():Observable<BalanceRsp>{
    return this.http.get("/api/v1/balance") as  Observable<BalanceRsp>
  }
  getMyTransactions():Observable<transactionCoinResponse>{
    return this.http.get("/api/v1/transactions") as  Observable<transactionCoinResponse>
  }
  getAllTransactions():Observable<transactionCoinResponse>{
    return this.http.get("/api/v1/transactionsAll") as  Observable<transactionCoinResponse>
  }
  postTransaction(to:number , coins:number ,happening$:BehaviorSubject<boolean>){
   happening$.next(true)
    const body:transactionCoinRequest = {
      to:to , 
      amount:coins
    }
    this.http.post("/api/v1/transfer" , body).subscribe((r)=>{
      this.navigateService.navigateTo("/")

    } ,err=>{ 

      this.openErrordialogService.errorDialog(err.error.Message)
       happening$.next(false)
    } ,()=>{} )


  }
}
