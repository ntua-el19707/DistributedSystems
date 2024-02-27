import { Component, OnInit } from '@angular/core';
import { CoinCardComponent } from '../../components/coin-card/coin-card.component';
import { NodeInfoModule } from '../../features/node-info/node-info.module';
import { NodeInfoBehaviorService } from '../../features/node-info/node-info-behavior.service';
import { TransactionBehaviorService } from '../../features/transaction-coins/transaction-behavior.service';
import { TransactionCoinsModule } from '../../features/transaction-coins/transaction-coins.module';
import { transactionCoinResponse } from '../../sharable';
import { BehaviorSubject } from 'rxjs';
import { AsyncPipe, CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { CoinTransactionTableComponent } from '../../components/coin-transaction-table/coin-transaction-table.component';

const  custom = [CoinCardComponent , NodeInfoModule,TransactionCoinsModule,CoinTransactionTableComponent
]
const common = [AsyncPipe , CommonModule]
const material = [MatCardModule ,MatInputModule]
@Component({
  selector: 'app-home-page',
  standalone: true,
  imports: [...custom ,...material , ...common],
  templateUrl: './home-page.component.html',
  styleUrl: './home-page.component.scss'
})
export class HomePageComponent implements OnInit {

   readonly dataSourceList$:BehaviorSubject<transactionCoinResponse>
  typingTimer: any;
  constructor(private transactionBehaviorService:TransactionBehaviorService ,private nodeInfoBehaviorService:NodeInfoBehaviorService){
    this.dataSourceList$ = this.transactionBehaviorService.getMyFilterTransactions()
  }
ngOnInit(): void {
    this.nodeInfoBehaviorService.fetchNodeInfo()
    this.transactionBehaviorService.fetchBalance()
    this.transactionBehaviorService.fetchMyTransactions()
}
applyFilter(event: Event){
    clearTimeout(this.typingTimer); // Clear any existing timer
    const filterValue = (event.target as HTMLInputElement).value;
   const numericFilterValue = parseFloat(filterValue)
    this.typingTimer = setTimeout(() => { 
     this.transactionBehaviorService.filterCoins(numericFilterValue)
    
  
    }, 1000);
  }
}
