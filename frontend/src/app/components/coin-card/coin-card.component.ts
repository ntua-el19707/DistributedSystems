import { Component, Input } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { nodeDetails, nodeInfoRsp } from '../../sharable';
import { NodeInfoModule } from '../../features/node-info/node-info.module';
import { NodeInfoBehaviorService } from '../../features/node-info/node-info-behavior.service';
import { BehaviorSubject } from 'rxjs';
import { AsyncPipe, CommonModule } from '@angular/common';
import { TransactionCoinsModule } from '../../features/transaction-coins/transaction-coins.module';
import { TransactionBehaviorService } from '../../features/transaction-coins/transaction-behavior.service';
import { TranferCoinsFormComponent } from '../tranfer-coins-form/tranfer-coins-form.component';
const material  = [MatButtonModule , MatCardModule]  
const  custom =  [NodeInfoModule ,TransactionCoinsModule]
const common = [CommonModule , AsyncPipe] 
@Component({
  selector: 'app-coin-card',
  standalone: true,
  imports: [...material , ...custom ,...common],
  templateUrl: './coin-card.component.html',
  styleUrl: './coin-card.component.scss'
})
export class CoinCardComponent {
  readonly dataSource$ :BehaviorSubject<nodeDetails>
  readonly coins$:BehaviorSubject<number>
  constructor(private nodeInfoBehaviorService :NodeInfoBehaviorService , private transactionBehaviorService:TransactionBehaviorService) {
    this.dataSource$ = this.nodeInfoBehaviorService.getNodeInfo()
    this.coins$ = this.transactionBehaviorService.getBalanceSubject()
  }
}
