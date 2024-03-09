import { Component, OnDestroy, OnInit } from '@angular/core';
import { CoinCardComponent } from '../../components/coin-card/coin-card.component';
import { NodeInfoModule } from '../../features/node-info/node-info.module';
import { NodeInfoBehaviorService } from '../../features/node-info/node-info-behavior.service';
import { TransactionBehaviorService } from '../../features/transaction-coins/transaction-behavior.service';
import { TransactionCoinsModule } from '../../features/transaction-coins/transaction-coins.module';
import { TransactionCoinsAndNodeDetails, transactionCoinResponse } from '../../sharable';
import { BehaviorSubject, Subscription, interval } from 'rxjs';
import { AsyncPipe, CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { CoinTransactionTableComponent } from '../../components/coin-transaction-table/coin-transaction-table.component';
import { filterTranctionCoinSubject } from '../../features/filter-transaction-coin-toa-behavir-subject.service';
import { FilterTransactionCoinsComponent } from '../../components/filter-transaction-coins/filter-transaction-coins.component';
import { throwDialogContentAlreadyAttachedError } from '@angular/cdk/dialog';

const  custom = [ FilterTransactionCoinsComponent, CoinCardComponent , NodeInfoModule,TransactionCoinsModule,CoinTransactionTableComponent
]
const common = [AsyncPipe , CommonModule]
const material = [MatCardModule ,MatInputModule]
@Component({
  selector: 'app-home-page',
  standalone: true,
  imports: [...custom, ...material, ...common],
  templateUrl: './home-page.component.html',
  styleUrl: './home-page.component.scss',
})
export class HomePageComponent implements OnInit, OnDestroy {
  readonly dataSourceList$: BehaviorSubject<TransactionCoinsAndNodeDetails>;

  readonly #Subscription: Subscription = new Subscription();
  #SubscriptionFetchData: Subscription = new Subscription();
  filter$: BehaviorSubject<filterTranctionCoinSubject> = new BehaviorSubject(
    {}
  );
  typingTimer: any;
  constructor(
    private transactionBehaviorService: TransactionBehaviorService,
    private nodeInfoBehaviorService: NodeInfoBehaviorService
  ) {
    this.dataSourceList$ =
      this.transactionBehaviorService.getMyFilterTransactions();
    this.#Subscription = this.filter$.subscribe((r) => {
      this.applyFilter(r);
    });

  }
  ngOnInit(): void {
    const  fetcingTime = 10000 //10sec 
    this.nodeInfoBehaviorService.fetchNodeInfo();
    this.transactionBehaviorService.fetchBalance();
    this.transactionBehaviorService.fetchMyTransactions(true);
    this.#SubscriptionFetchData = interval(fetcingTime).subscribe(r=>{
      this.transactionBehaviorService.fetchMyTransactions(false)
      this.transactionBehaviorService.fetchBalance()
    })
  }
  ngOnDestroy(): void {
    this.#Subscription.unsubscribe();
    this.#SubscriptionFetchData.unsubscribe()
  }
  private applyFilter(filter: filterTranctionCoinSubject) {
    clearTimeout(this.typingTimer); // Clear any existing timer

    this.typingTimer = setTimeout(() => {
      this.transactionBehaviorService.filterCoins(filter);
    }, 1000);
  }
}
