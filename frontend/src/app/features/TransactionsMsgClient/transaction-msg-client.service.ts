import { Injectable } from '@angular/core';
import { TransactionsMsgClientModule } from './transactions-msg-client.module';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { transactionMsgRequest, transactionMsgResponse } from '../../sharable';

@Injectable({
  providedIn: TransactionsMsgClientModule
})
export class TransactionMsgClientService {

  constructor(private http:HttpClient) { }
  getINBOX():Observable<transactionMsgResponse>{

    return  this.http.get("/api/v1/inbox") as Observable<transactionMsgResponse>
  }
    getSend():Observable<transactionMsgResponse>{

    return  this.http.get("/api/v1/send") as Observable<transactionMsgResponse>
  }
   getAll():Observable<transactionMsgResponse>{

    return  this.http.get("/api/v1/allMsg") as Observable<transactionMsgResponse>
  }
  postTransaction(To:number , Msg:string){
   const  body: transactionMsgRequest = {
    to:To ,
    msg:Msg
   }
   this.http.post("/api/v1/send" , body).subscribe(r=>{},err=>{} , ()=>{})
    
  }
}
