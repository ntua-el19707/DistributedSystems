import { Component, Input } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { nodeDetails, nodeInfoRsp } from '../../sharable';
import { NodeInfoModule } from '../../features/node-info/node-info.module';
import { NodeInfoBehaviorService } from '../../features/node-info/node-info-behavior.service';
import { BehaviorSubject } from 'rxjs';
import { AsyncPipe, CommonModule } from '@angular/common';
const material  = [MatButtonModule , MatCardModule]  
const  custom =  [NodeInfoModule]
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
  constructor(private nodeInfoBehaviorService :NodeInfoBehaviorService) {
    this.dataSource$ = this.nodeInfoBehaviorService.getNodeInfo()
  }
}
