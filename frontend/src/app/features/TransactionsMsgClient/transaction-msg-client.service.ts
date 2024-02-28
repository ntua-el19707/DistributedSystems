import { Injectable } from '@angular/core';
import { TransactionsMsgClientModule } from './transactions-msg-client.module';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import { transactionMsgRequest, transactionMsgResponse } from '../../sharable';
import { OpenErrordialogService } from '../open-errordialog.service';
import { NavigateService } from '../navigate.service';

@Injectable({
  providedIn: TransactionsMsgClientModule
})
export class TransactionMsgClientService {

  constructor(private http:HttpClient, private navigateService:NavigateService ,private openErrordialogService:OpenErrordialogService) { }
  getINBOX():Observable<transactionMsgResponse>{

    return  this.http.get("/api/v1/inbox") as Observable<transactionMsgResponse>
  }
    getSend():Observable<transactionMsgResponse>{

    return  this.http.get("/api/v1/send") as Observable<transactionMsgResponse>
  }
   getAll():Observable<transactionMsgResponse>{

    return  this.http.get("/api/v1/allMsg") as Observable<transactionMsgResponse>
  }
  
  postTransaction(To:number , Msg:string ,happening:BehaviorSubject<boolean>){
    happening.next(true)
   const  body: transactionMsgRequest = {
    to:To ,
    msg:Msg
   }
   this.http.post("/api/v1/send" , body).subscribe(r=>{ 
    this.navigateService.navigateTo("/#/send")
},err=>{
  this.openErrordialogService.errorDialog(err.error.Message)
  happening.next(false)

} , ()=>{
 
})
    
  }
}
