
import { Component, Input } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { NodeInfoGraphQL, nodeDetails, nodeInfoRsp } from '../../sharable';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatCardModule } from '@angular/material/card';
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';

import { MatButtonModule } from '@angular/material/button';
import { TransactionMsgClientService } from '../../features/TransactionsMsgClient/transaction-msg-client.service';
import { TransactionsMsgClientModule } from '../../features/TransactionsMsgClient/transactions-msg-client.module';
import { BehaviorSubject } from 'rxjs';
import { CommonModule } from '@angular/common';
const material = [MatProgressSpinnerModule ,MatFormFieldModule, MatSelectModule, MatInputModule,MatCardModule ,MatButtonModule]
const forms = [ FormsModule, ReactiveFormsModule,]
@Component({
  selector: 'app-send-message-form',
  standalone: true,
  imports: [...material , ...forms,TransactionsMsgClientModule,CommonModule],
  templateUrl: './send-message-form.component.html',
  styleUrl: './send-message-form.component.scss'
})
export class SendMessageFormComponent {

TransferForm: FormGroup<transactionMsgForm> = new FormGroup<transactionMsgForm>({To:new FormControl<nodeDetails | null>(null,Validators.required ) ,  Msg:new FormControl<string | null>("",[Validators.required])})
  
  @Input()clients: NodeInfoGraphQL[] = []
  happening$:BehaviorSubject< boolean>  = new BehaviorSubject<boolean>(false)
  constructor(private  transactionMsgClientService:TransactionMsgClientService){}
  submit(){
  
  const isHappening = this.happening$.getValue()

  if (isHappening) {
    return; 
  }
    const target = this.TransferForm.controls.To.value?.indexId 

    const msg =  this.TransferForm.controls.Msg.value 
    if (target !== undefined && msg !== null) {
      this.transactionMsgClientService.postTransaction(
        target,
        msg,
        this.happening$
      );
    }
    
  

  }

}

``

interface  transactionMsgForm {
  To: FormControl<NodeInfoGraphQL |null> 
  Msg: FormControl<string| null >
}
