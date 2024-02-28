
import { Component, Input } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { nodeDetails } from '../../sharable';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatCardModule } from '@angular/material/card';

import { MatButtonModule } from '@angular/material/button';
import { TransactionMsgClientService } from '../../features/TransactionsMsgClient/transaction-msg-client.service';
import { TransactionsMsgClientModule } from '../../features/TransactionsMsgClient/transactions-msg-client.module';
const material = [MatFormFieldModule, MatSelectModule, MatInputModule,MatCardModule ,MatButtonModule]
const forms = [ FormsModule, ReactiveFormsModule,]
@Component({
  selector: 'app-send-message-form',
  standalone: true,
  imports: [...material , ...forms,TransactionsMsgClientModule],
  templateUrl: './send-message-form.component.html',
  styleUrl: './send-message-form.component.scss'
})
export class SendMessageFormComponent {

TransferForm: FormGroup<transactionMsgForm> = new FormGroup<transactionMsgForm>({To:new FormControl<nodeDetails | null>(null,Validators.required ) ,  Msg:new FormControl<string | null>("",[Validators.required])})
  @Input()clients: nodeDetails[] = []
  constructor(private  transactionMsgClientService:TransactionMsgClientService){}
  submit(){
    const target = this.TransferForm.controls.To.value?.indexId 

    const msg =  this.TransferForm.controls.Msg.value 
    if (target && msg) {
      console.log(target , msg )
    this.transactionMsgClientService.postTransaction(target , msg)
  }

  }

}

``

interface  transactionMsgForm {
  To: FormControl<nodeDetails |null> 
  Msg: FormControl<string| null >
}
