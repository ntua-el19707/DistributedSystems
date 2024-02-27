import { Component } from '@angular/core';
import { ClientsModule } from '../../features/clients/clients.module';
import { AsyncPipe, CommonModule } from '@angular/common';
import { SendMessageFormComponent } from '../../components/send-message-form/send-message-form.component';
import { BehaviorSubject } from 'rxjs';
import { nodeDetails } from '../../sharable';
import { ClientsBehaviorSubjectService } from '../../features/clients/clients-behavior-subject.service';

@Component({
  selector: 'app-send-message-form-page',
  standalone: true,
  imports: [SendMessageFormComponent , ClientsModule ,AsyncPipe ,CommonModule],
  templateUrl: './send-message-form-page.component.html',
  styleUrl: './send-message-form-page.component.scss'
})
export class SendMessageFormPageComponent {
readonly dataSource$:BehaviorSubject<Array<nodeDetails>>
constructor(private clientsBehaviorSubjectService:ClientsBehaviorSubjectService){
    this.clientsBehaviorSubjectService.fetchClients()
  this.dataSource$ = this.clientsBehaviorSubjectService.getBehavior()
}
}
