import {  TransactionMsgRowGraphQL } from '../../sharable';
import {

  Component,
  Input,
  OnInit,

} from '@angular/core';
import {  MatPaginatorModule } from '@angular/material/paginator';

import {  MatTableModule } from '@angular/material/table';
import { UnixDateFormatPipe } from '../../features/pipes/unix-date-format.pipe';

import { AsyncPipe, CommonModule } from '@angular/common';

const materail = [MatTableModule, MatPaginatorModule];
@Component({
  selector: 'app-transaction-msg-table',
  standalone: true,
  imports: [...materail, UnixDateFormatPipe, AsyncPipe, CommonModule],
  templateUrl: './transaction-msg-table.component.html',
  styleUrl: './transaction-msg-table.component.scss',
})
export class TransactionMsgTableComponent implements OnInit {
  displayedColumns: string[] = ['id', 'From', 'Nonce', 'To', 'Msg', 'SendTime'];

  @Input() data: Array<TransactionMsgRowGraphQL> = [];
  ngOnInit(): void {}
}
