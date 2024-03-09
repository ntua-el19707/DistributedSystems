import { Component, OnDestroy } from '@angular/core';
import { TransactionMsgTableComponent } from '../../components/transaction-msg-table/transaction-msg-table.component';
import { AsyncPipe, CommonModule } from '@angular/common';
import { TransactionMsgList, TransactionMsgListGraphQL, TransactionMsgRow, TransactionMsgRowGraphQL, transactionMsgResponse } from '../../sharable';
import { BehaviorSubject, Subscription, interval} from 'rxjs';
import { MatCardModule } from '@angular/material/card';
import { TransactionsMsgClientModule } from '../../features/TransactionsMsgClient/transactions-msg-client.module';
import { TransactionMsgBehaviorService } from '../../features/TransactionsMsgClient/transaction-msg-behavior.service';
import {MatInputModule} from '@angular/material/input';
import {MatFormFieldModule} from '@angular/material/form-field';
import { FilterTransactionMsg } from '../../features/filter-transaction-msg-toa-behavir-subject.service';
import { FilterTransactionMsgComponent } from '../../components/filter-transaction-msg/filter-transaction-msg.component';
const custom = [FilterTransactionMsgComponent, TransactionMsgTableComponent ,TransactionsMsgClientModule]
const common = [AsyncPipe , CommonModule]
const material = [MatCardModule ,MatInputModule]
@Component({
  selector: 'app-inbox',
  standalone: true,
  imports: [...custom, ...common, ...material],
  templateUrl: './inbox.component.html',
  styleUrl: './inbox.component.scss',
})
export class InboxComponent implements OnDestroy {
  readonly dataSource$: BehaviorSubject<TransactionMsgList>;
  readonly dataSourceList$: BehaviorSubject<Array<TransactionMsgRowGraphQL>>;

  to: boolean = false;
  from: boolean = true;
  readonly #Subscription: Subscription = new Subscription();
  #SubscriptionFetchData: Subscription = new Subscription();
  filter$: BehaviorSubject<FilterTransactionMsg> = new BehaviorSubject({});
  constructor(
    private transactionMsgBehaviorService: TransactionMsgBehaviorService
  ) {
    this.transactionMsgBehaviorService.fetchInbox(true);
    this.dataSource$ = transactionMsgBehaviorService.getInbox();
    this.dataSourceList$ = transactionMsgBehaviorService.getInboxTransactions();
    this.#Subscription = this.filter$.subscribe((r) => {
      this.applyFilter(r);
    });
    const fetcingTime = 10000; //10sec

    this.#SubscriptionFetchData = interval(fetcingTime).subscribe((r) => {
      this.transactionMsgBehaviorService.fetchInbox(false);
    });
  }
  ngOnDestroy(): void {
    this.#Subscription.unsubscribe();
    this.#SubscriptionFetchData.unsubscribe()
  }
  private applyFilter(filter: FilterTransactionMsg) {
    this.transactionMsgBehaviorService.filter(filter);
  }
}

