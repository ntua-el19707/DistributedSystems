import { Component, OnInit } from '@angular/core';
import { CoinCardComponent } from '../../components/coin-card/coin-card.component';
import { NodeInfoModule } from '../../features/node-info/node-info.module';
import { NodeInfoBehaviorService } from '../../features/node-info/node-info-behavior.service';

const  custom = [CoinCardComponent , NodeInfoModule]
@Component({
  selector: 'app-home-page',
  standalone: true,
  imports: [...custom],
  templateUrl: './home-page.component.html',
  styleUrl: './home-page.component.scss'
})
export class HomePageComponent implements OnInit {
  constructor(private nodeInfoBehaviorService:NodeInfoBehaviorService){}
ngOnInit(): void {
    this.nodeInfoBehaviorService.fetchNodeInfo()
}
}
