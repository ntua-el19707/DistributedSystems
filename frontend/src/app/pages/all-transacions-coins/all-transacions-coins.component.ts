import { Component, OnInit } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { transactionCoinResponse } from '../../sharable';
import { TransactionBehaviorService } from '../../features/transaction-coins/transaction-behavior.service';
import { TransactionCoinsModule } from '../../features/transaction-coins/transaction-coins.module';
import { CoinTransactionTableComponent } from '../../components/coin-transaction-table/coin-transaction-table.component';
import { AsyncPipe, CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
const  custom = [TransactionCoinsModule,CoinTransactionTableComponent]
const common = [AsyncPipe , CommonModule]
const material = [MatCardModule ,MatInputModule]
@Component({
  selector: 'app-all-transacions-coins',
  standalone: true,
  imports: [...custom , ...common , ...material],
  templateUrl: './all-transacions-coins.component.html',
  styleUrl: './all-transacions-coins.component.scss'
})
export class AllTransacionsCoinsComponent implements OnInit{
 readonly dataSourceList$:BehaviorSubject<transactionCoinResponse>
  typingTimer: any;
  constructor(private transactionBehaviorService:TransactionBehaviorService ){
    this.dataSourceList$ = this.transactionBehaviorService.getMyAllFilterTransactions()
  }
ngOnInit(): void {
    
    this.transactionBehaviorService.fetchAllTransactions()
}
applyFilter(event: Event){
    clearTimeout(this.typingTimer); // Clear any existing timer
    const filterValue = (event.target as HTMLInputElement).value;
   const numericFilterValue = parseFloat(filterValue)
    this.typingTimer = setTimeout(() => { 
     this.transactionBehaviorService.filterAllCoins(numericFilterValue)
    
  
    }, 1000);
  }
}
