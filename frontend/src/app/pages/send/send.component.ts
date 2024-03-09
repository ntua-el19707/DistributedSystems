
import { Component, OnDestroy } from '@angular/core';
import { TransactionMsgTableComponent } from '../../components/transaction-msg-table/transaction-msg-table.component';
import { AsyncPipe, CommonModule } from '@angular/common';
import { TransactionMsgList, TransactionMsgListGraphQL, TransactionMsgRowGraphQL } from '../../sharable';
import { BehaviorSubject, Subscription, interval} from 'rxjs';
import { MatCardModule } from '@angular/material/card';
import { TransactionsMsgClientModule } from '../../features/TransactionsMsgClient/transactions-msg-client.module';
import { TransactionMsgBehaviorService } from '../../features/TransactionsMsgClient/transaction-msg-behavior.service';
import {MatInputModule} from '@angular/material/input';

import { FilterTransactionMsg } from '../../features/filter-transaction-msg-toa-behavir-subject.service';
import { FilterTransactionMsgComponent } from '../../components/filter-transaction-msg/filter-transaction-msg.component';
const custom = [    FilterTransactionMsgComponent,TransactionMsgTableComponent ,TransactionsMsgClientModule]
const common = [AsyncPipe , CommonModule]
const material = [MatCardModule ,MatInputModule]
@Component({
  selector: 'app-send',
  standalone: true,
  imports: [...custom, ...material, ...common],
  templateUrl: './send.component.html',
  styleUrl: './send.component.scss',
})
export class SendComponent implements OnDestroy {
  readonly dataSource$: BehaviorSubject<TransactionMsgList>;
  readonly dataSourceList$: BehaviorSubject<Array<TransactionMsgRowGraphQL>>;
  to: boolean = true;
  from: boolean = false;
  readonly #Subscription: Subscription = new Subscription();
  #SubscriptionFetchData: Subscription = new Subscription();
  filter$: BehaviorSubject<FilterTransactionMsg> = new BehaviorSubject({});
  constructor(
    private transactionMsgBehaviorService: TransactionMsgBehaviorService
  ) {
    this.transactionMsgBehaviorService.fetchSend(true);
    this.dataSource$ = transactionMsgBehaviorService.getSend();
    this.dataSourceList$ = transactionMsgBehaviorService.getSendTransactions();
    this.#Subscription = this.filter$.subscribe((r) => {
      this.applyFilter(r);
    });
    const fetcingTime = 10000; //10sec

    this.#SubscriptionFetchData = interval(fetcingTime).subscribe((r) => {
      this.transactionMsgBehaviorService.fetchSend(false);
    });
  }
  ngOnDestroy(): void {
    this.#Subscription.unsubscribe();
    this.#SubscriptionFetchData.unsubscribe();
  }
  private applyFilter(filter: FilterTransactionMsg) {
    this.transactionMsgBehaviorService.filterSend(filter);
  }
}
