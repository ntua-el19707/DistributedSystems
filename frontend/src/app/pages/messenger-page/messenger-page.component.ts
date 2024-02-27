import { Component } from '@angular/core';
import { TransactionMsgTableComponent } from '../../components/transaction-msg-table/transaction-msg-table.component';
const custom = [TransactionMsgTableComponent]
@Component({
  selector: 'app-messenger-page',
  standalone: true,
  imports: [],
  templateUrl: './messenger-page.component.html',
  styleUrl: './messenger-page.component.scss'
})
export class MessengerPageComponent {

}
