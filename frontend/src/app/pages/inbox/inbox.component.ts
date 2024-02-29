import { Component } from '@angular/core';
import { TransactionMsgTableComponent } from '../../components/transaction-msg-table/transaction-msg-table.component';
import { AsyncPipe, CommonModule } from '@angular/common';
import { TransactionMsgRow, transactionMsgResponse } from '../../sharable';
import { BehaviorSubject} from 'rxjs';
import { MatCardModule } from '@angular/material/card';
import { TransactionsMsgClientModule } from '../../features/TransactionsMsgClient/transactions-msg-client.module';
import { TransactionMsgBehaviorService } from '../../features/TransactionsMsgClient/transaction-msg-behavior.service';
import {MatInputModule} from '@angular/material/input';
import {MatFormFieldModule} from '@angular/material/form-field';
const custom = [TransactionMsgTableComponent ,TransactionsMsgClientModule]
const common = [AsyncPipe , CommonModule]
const material = [MatCardModule ,MatInputModule]
@Component({
  selector: 'app-inbox',
  standalone: true,
  imports: [...custom ,...common,...material],
  templateUrl: './inbox.component.html',
  styleUrl: './inbox.component.scss'
})
export class InboxComponent {
  readonly dataSource$:BehaviorSubject<transactionMsgResponse>
   readonly dataSourceList$:BehaviorSubject<Array<TransactionMsgRow>>
  typingTimer: any;

  constructor(private transactionMsgBehaviorService:TransactionMsgBehaviorService){
    this.transactionMsgBehaviorService.fetchInbox()
    this.dataSource$ = transactionMsgBehaviorService.getInbox()
    this.dataSourceList$ = transactionMsgBehaviorService.getInboxTransactions()
  }
applyFilter(event: Event){
    clearTimeout(this.typingTimer); // Clear any existing timer
    const filterValue = (event.target as HTMLInputElement).value;


    this.typingTimer = setTimeout(() => {
 
     this.transactionMsgBehaviorService.filter(filterValue)
    
  
    }, 1000);
  }


}

