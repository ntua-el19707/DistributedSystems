import { Component } from '@angular/core';
import { TranferCoinsFormComponent } from '../../components/tranfer-coins-form/tranfer-coins-form.component';
import { BehaviorSubject } from 'rxjs';
import { nodeDetails } from '../../sharable';
import { ClientsBehaviorSubjectService } from '../../features/clients/clients-behavior-subject.service';
import { ClientsModule } from '../../features/clients/clients.module';
import { AsyncPipe, CommonModule } from '@angular/common';

@Component({
  selector: 'app-trasfer-coins',
  standalone: true,
  imports: [TranferCoinsFormComponent , ClientsModule ,AsyncPipe ,CommonModule
  ],
  templateUrl: './trasfer-coins.component.html',
  styleUrl: './trasfer-coins.component.scss'
})
export class TrasferCoinsComponent {
readonly dataSource$:BehaviorSubject<Array<nodeDetails>>
constructor(private clientsBehaviorSubjectService:ClientsBehaviorSubjectService){
  this.clientsBehaviorSubjectService.fetchClients()
this.dataSource$ = this.clientsBehaviorSubjectService.getBehavior()
}
}
