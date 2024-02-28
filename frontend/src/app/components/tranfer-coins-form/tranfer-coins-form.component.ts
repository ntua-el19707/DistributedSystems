import { Component, Input } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { nodeDetails } from '../../sharable';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatCardModule } from '@angular/material/card';

import { MatButtonModule } from '@angular/material/button';
import { TransactionClientService } from '../../features/transaction-coins/transaction-client.service';
import { TransactionCoinsModule } from '../../features/transaction-coins/transaction-coins.module';
import { ErrorDialogComponent } from '../error-dialog/error-dialog.component';
import { MatDialog } from '@angular/material/dialog';
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import { BehaviorSubject } from 'rxjs';
import { CommonModule } from '@angular/common';

const material = [MatProgressSpinnerModule ,MatFormFieldModule, MatSelectModule, MatInputModule,MatCardModule ,MatButtonModule]
const forms = [ FormsModule, ReactiveFormsModule,]
@Component({
  selector: 'app-tranfer-coins-form',
  standalone: true,
  imports: [...forms ,  ...material ,CommonModule,TransactionCoinsModule ,ErrorDialogComponent],
  templateUrl: './tranfer-coins-form.component.html',
  styleUrl: './tranfer-coins-form.component.scss'
})
export class TranferCoinsFormComponent {
happening$:BehaviorSubject< boolean>  = new BehaviorSubject<boolean>(false)
  TransferForm: FormGroup<transactionCoinsForm> = new FormGroup<transactionCoinsForm>({To:new FormControl<nodeDetails | null>(null,Validators.required ) ,  Coins:new FormControl<string | null>(null,[Validators.required,Validators.pattern(/^[0-9]+(\.[0-9]+)?$/)])})
  @Input()clients: nodeDetails[] = []
    constructor(private  transactionClientService:TransactionClientService , public dialog: MatDialog){}
submit(){
    const isHappening = this.happening$.getValue()

  if (isHappening) {
    return; 
  }
    const target = this.TransferForm.controls.To.value?.indexId 

    const coins = this.TransferForm.controls.Coins.value 
       if (target && coins) {

        this.transactionClientService.postTransaction(target , parseFloat(coins) ,this.happening$)
     

  }

  }
 


}


interface  transactionCoinsForm {
  To: FormControl<nodeDetails |null> 
  Coins: FormControl<string| null >
}