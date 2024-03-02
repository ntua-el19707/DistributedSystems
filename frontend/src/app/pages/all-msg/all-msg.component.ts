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
const custom = [ FilterTransactionMsgComponent, TransactionMsgTableComponent ,TransactionsMsgClientModule]
const common = [AsyncPipe , CommonModule]
const material = [MatCardModule ,MatInputModule]
@Component({
  selector: 'app-all-msg',
  standalone: true,
  imports: [...common , ...custom , ...material],
  templateUrl: './all-msg.component.html',
  styleUrl: './all-msg.component.scss'
})
export class AllMsgComponent implements OnDestroy {
readonly dataSource$:BehaviorSubject<transactionMsgResponse>
   readonly dataSourceList$:BehaviorSubject<Array<TransactionMsgRow>>

to:boolean = true
  from:boolean = true
  readonly #Subscription:Subscription = new Subscription()
  filter$ :BehaviorSubject<FilterTransactionMsg> = new  BehaviorSubject({})
  constructor(private transactionMsgBehaviorService:TransactionMsgBehaviorService){
    this.transactionMsgBehaviorService.fetchAll()
    this.dataSource$ = transactionMsgBehaviorService.getAll()
    this.dataSourceList$ = transactionMsgBehaviorService.getAllTransactions()
   this.#Subscription = this.filter$.subscribe(r=>{this.applyFilter(r)} )
  }
  ngOnDestroy(): void {
    this.#Subscription.unsubscribe()
  }
private applyFilter(filter:FilterTransactionMsg){
   this.transactionMsgBehaviorService.filterAll(filter)
}

}
