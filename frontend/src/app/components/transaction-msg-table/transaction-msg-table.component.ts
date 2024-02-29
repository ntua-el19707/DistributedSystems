
import { TransactionMsgRow } from '../../sharable';
import {AfterViewInit, Component, Input, OnInit, ViewChild} from '@angular/core';
import {MatPaginator, MatPaginatorModule} from '@angular/material/paginator';

import {MatTableDataSource, MatTableModule} from '@angular/material/table';
import { UnixDateFormatPipe } from '../../features/pipes/unix-date-format.pipe';
import { BehaviorSubject } from 'rxjs';
import { AsyncPipe, CommonModule } from '@angular/common';
import { TransactionsMsgClientModule } from '../../features/TransactionsMsgClient/transactions-msg-client.module';


const materail = [  MatTableModule, MatPaginatorModule]
@Component({
  selector: 'app-transaction-msg-table',
  standalone: true,
  imports: [...materail , UnixDateFormatPipe ,AsyncPipe,CommonModule],
  templateUrl: './transaction-msg-table.component.html',
  styleUrl: './transaction-msg-table.component.scss'
})
export class TransactionMsgTableComponent implements OnInit {
 
  displayedColumns: string[] = ['id', 'From','Nonce' , 'To', 'Msg' , 'SendTime'];
  
  @Input() data :Array<TransactionMsgRow> = [] 
  ngOnInit(): void {
   
  }



}
