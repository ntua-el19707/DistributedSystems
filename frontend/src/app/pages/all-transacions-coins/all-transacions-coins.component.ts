import { Component, OnDestroy, OnInit } from '@angular/core';
import { BehaviorSubject, Subscription, interval } from 'rxjs';
import { TransactionCoinsAndNodeDetails, TransactionCoinsRowGraphQL, transactionCoinResponse } from '../../sharable';
import { TransactionBehaviorService } from '../../features/transaction-coins/transaction-behavior.service';
import { TransactionCoinsModule } from '../../features/transaction-coins/transaction-coins.module';
import { CoinTransactionTableComponent } from '../../components/coin-transaction-table/coin-transaction-table.component';
import { AsyncPipe, CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { FilterTransactionCoinsComponent } from '../../components/filter-transaction-coins/filter-transaction-coins.component';
import { filterTranctionCoinSubject } from '../../features/filter-transaction-coin-toa-behavir-subject.service';
const  custom = [TransactionCoinsModule,CoinTransactionTableComponent,FilterTransactionCoinsComponent]
const common = [AsyncPipe , CommonModule]
const material = [MatCardModule ,MatInputModule]
@Component({
  selector: 'app-all-transacions-coins',
  standalone: true,
  imports: [...custom, ...common, ...material],
  templateUrl: './all-transacions-coins.component.html',
  styleUrl: './all-transacions-coins.component.scss',
})
export class AllTransacionsCoinsComponent implements OnInit, OnDestroy {
  readonly dataSourceList$: BehaviorSubject<TransactionCoinsAndNodeDetails>;
  readonly #Subscription: Subscription = new Subscription();
  #SubscriptionFetchData: Subscription = new Subscription();
  filter$: BehaviorSubject<filterTranctionCoinSubject> = new BehaviorSubject(
    {}
  );
  typingTimer: any;
  constructor(private transactionBehaviorService: TransactionBehaviorService) {
    this.dataSourceList$ =
      this.transactionBehaviorService.getMyAllFilterTransactions();
    this.#Subscription = this.filter$.subscribe((r) => {
      this.applyFilter(r);
    });
  }
  ngOnInit(): void {
    const fetcingTime = 10000; //10sec
    this.transactionBehaviorService.fetchAllTransactions(true);
    this.#SubscriptionFetchData = interval(fetcingTime).subscribe((r) => {
      this.transactionBehaviorService.fetchAllTransactions(false);
    });
  }
  ngOnDestroy(): void {
    this.#Subscription.unsubscribe();
    this.#SubscriptionFetchData.unsubscribe()
  }
  private applyFilter(filter: filterTranctionCoinSubject) {
    clearTimeout(this.typingTimer); // Clear any existing timer

    this.typingTimer = setTimeout(() => {
      this.transactionBehaviorService.filterAllCoins(filter);
    }, 1000);
  }
}
