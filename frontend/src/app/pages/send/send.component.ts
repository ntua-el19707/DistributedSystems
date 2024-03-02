
import { Component, OnDestroy } from '@angular/core';
import { TransactionMsgTableComponent } from '../../components/transaction-msg-table/transaction-msg-table.component';
import { AsyncPipe, CommonModule } from '@angular/common';
import { TransactionMsgRow, transactionMsgResponse } from '../../sharable';
import { BehaviorSubject, Subscription} from 'rxjs';
import { MatCardModule } from '@angular/material/card';
import { TransactionsMsgClientModule } from '../../features/TransactionsMsgClient/transactions-msg-client.module';
import { TransactionMsgBehaviorService } from '../../features/TransactionsMsgClient/transaction-msg-behavior.service';
import {MatInputModule} from '@angular/material/input';
import {MatFormFieldModule} from '@angular/material/form-field';
import { FilterTransactionMsg } from '../../features/filter-transaction-msg-toa-behavir-subject.service';
import { FilterTransactionMsgComponent } from '../../components/filter-transaction-msg/filter-transaction-msg.component';
const custom = [    FilterTransactionMsgComponent,TransactionMsgTableComponent ,TransactionsMsgClientModule]
const common = [AsyncPipe , CommonModule]
const material = [MatCardModule ,MatInputModule]
@Component({
  selector: 'app-send',
  standalone: true,
  imports: [...custom ,...material , ...common ],
  templateUrl: './send.component.html',
  styleUrl: './send.component.scss'
})
export class SendComponent implements OnDestroy{
  readonly dataSource$:BehaviorSubject<transactionMsgResponse>
  readonly dataSourceList$:BehaviorSubject<Array<TransactionMsgRow>>
  to:boolean = true  
  from:boolean = false 
  readonly #Subscription:Subscription = new Subscription()
  filter$ :BehaviorSubject<FilterTransactionMsg> = new  BehaviorSubject({})
  constructor(private transactionMsgBehaviorService:TransactionMsgBehaviorService){
    this.transactionMsgBehaviorService.fetchSend()
    this.dataSource$ = transactionMsgBehaviorService.getSend()
    this.dataSourceList$ = transactionMsgBehaviorService.getSendTransactions()
    this.#Subscription = this.filter$.subscribe(r=>{this.applyFilter(r)} )
  }
  ngOnDestroy(): void {
    this.#Subscription.unsubscribe()
  }
private applyFilter(filter:FilterTransactionMsg){
   this.transactionMsgBehaviorService.filterSend(filter)
}


}
