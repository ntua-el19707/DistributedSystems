import { Component, Input } from '@angular/core';
import { TransactionCoinsAndNodeDetails, TransactionCoinsRowGraphQL, transactionCoinResponse } from '../../sharable';
import { SendCoinspipePipe } from '../../features/pipes/send-coinspipe.pipe';
import { UnixDateFormatPipe } from '../../features/pipes/unix-date-format.pipe';
import { MatTableModule } from '@angular/material/table';
const materail = [  MatTableModule]
@Component({
  selector: 'app-coin-transaction-table',
  standalone: true,
  imports: [UnixDateFormatPipe ,...materail],
  templateUrl: './coin-transaction-table.component.html',
  styleUrl: './coin-transaction-table.component.scss'
})
export class CoinTransactionTableComponent {

  displayedColumns: string[] = ['id', 'From', 'To' ,'Nonce', 'Coins' , 'Reason', 'SendTime'];
  
  @Input() data : TransactionCoinsAndNodeDetails= {transactions:[] , nodeDetails:{indexId:0 , nodeId:"" ,uri:"" ,uriPublic:""}}
  ngOnInit(): void {
   
  }
}
