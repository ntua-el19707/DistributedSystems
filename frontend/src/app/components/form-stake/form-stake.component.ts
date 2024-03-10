import { Component } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { BehaviorSubject } from 'rxjs';
import { TransactionClientService } from '../../features/transaction-coins/transaction-client.service';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatFormFieldModule } from '@angular/material/form-field';

import { MatInputModule } from '@angular/material/input';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { TransactionCoinsModule } from '../../features/transaction-coins/transaction-coins.module';
import { CommonModule } from '@angular/common';
const material = [
  MatProgressSpinnerModule,
  MatFormFieldModule,
   MatInputModule , 
  MatCardModule,
  MatButtonModule,
];
const forms = [FormsModule, ReactiveFormsModule];
@Component({
  selector: 'app-form-stake',
  standalone: true,
  imports: [...material ,  ...forms , TransactionCoinsModule , CommonModule ],
  templateUrl: './form-stake.component.html',
  styleUrl: './form-stake.component.scss',
})
export class FormStakeComponent {
  happening$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  TransferForm: FormGroup<transactionStakeForm> =
    new FormGroup<transactionStakeForm>({
      Coins: new FormControl<string | null>(null, [
        Validators.required,
        Validators.pattern(/^[0-9]+(\.[0-9]+)?$/),
      ]),
    });

  constructor(
    private transactionClientService: TransactionClientService,

  ) {}
  submit() {
    const isHappening = this.happening$.getValue();

    if (isHappening) {
      return;
    }


    const coins = this.TransferForm.controls.Coins.value;

    if (coins !== null) {
      this.transactionClientService.postStake(
        parseFloat(coins),
        this.happening$
      );
    }
  }
}

interface transactionStakeForm {
  Coins: FormControl<string | null>;
}