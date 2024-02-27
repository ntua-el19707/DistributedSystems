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
  selector: 'app-all-msg',
  standalone: true,
  imports: [...common , ...custom , ...material],
  templateUrl: './all-msg.component.html',
  styleUrl: './all-msg.component.scss'
})
export class AllMsgComponent {
readonly dataSource$:BehaviorSubject<transactionMsgResponse>
   readonly dataSourceList$:BehaviorSubject<Array<TransactionMsgRow>>
  typingTimer: any;

  constructor(private transactionMsgBehaviorService:TransactionMsgBehaviorService){
    this.transactionMsgBehaviorService.fetchAll()
    this.dataSource$ = transactionMsgBehaviorService.getAll()
    this.dataSourceList$ = transactionMsgBehaviorService.getAllTransactions()
  }
applyFilter(event: Event){
    clearTimeout(this.typingTimer); // Clear any existing timer
    const filterValue = (event.target as HTMLInputElement).value;


    this.typingTimer = setTimeout(() => {
 
     this.transactionMsgBehaviorService.filterAll(filterValue)
    
  
    }, 1000);
  }

}
