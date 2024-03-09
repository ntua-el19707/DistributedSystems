import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { NodeInfoGraphQL, nodeDetails } from '../../sharable';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatInputModule } from '@angular/material/input';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { DateAndTimeForm, DateAndTimePickerComponent, DatePickerForm, TimeValidator, timePattern } from '../date-and-time-picker/date-and-time-picker.component';
import { BehaviorSubject, Subscription } from 'rxjs';
import { FilterTransactionCoinToaBehavirSubjectService, filterTranctionCoinSubject } from '../../features/filter-transaction-coin-toa-behavir-subject.service';
import { AllNodesBehaviorSubjectService } from '../../features/allNodes/all-nodes-behavior-subject.service';
import { AsyncPipe, CommonModule } from '@angular/common';
const material = [MatFormFieldModule, MatSelectModule, MatInputModule,MatCardModule ,MatButtonModule]
const forms = [FormsModule, ReactiveFormsModule]
@Component({
  selector: 'app-filter-transaction-coins',
  standalone: true,
  imports: [...material, ...forms, DateAndTimePickerComponent ,CommonModule, AsyncPipe],
  templateUrl: './filter-transaction-coins.component.html',
  styleUrl: './filter-transaction-coins.component.scss',
})
export class FilterTransactionCoinsComponent implements OnDestroy {

  reasons: Array<string> = ['fee', 'Transfer', 'BootStrap'];
  readonly #Subscription: Subscription = new Subscription();
  readonly clients$ :BehaviorSubject<Array<NodeInfoGraphQL>>
  @Input() filter$: BehaviorSubject<filterTranctionCoinSubject> =
    new BehaviorSubject({});
  constructor(
    private allNodesBehaviorSubjectService: AllNodesBehaviorSubjectService,
    private filterTransactionCoinToaBehavirSubjectService: FilterTransactionCoinToaBehavirSubjectService
  ) {
    this.#Subscription = this.DateAndTimeForm$.subscribe((r) => {
      this.filterForm.controls.DateAndTimeForm = r;
      this.filterTransactionCoinToaBehavirSubjectService.next(
        this.filter$,
        this.filterForm
      );
    });
    this.clients$ = this.allNodesBehaviorSubjectService.getBehaviorSubject()
  }
  ngOnDestroy(): void {
    this.#Subscription.unsubscribe();
  }
  DateAndTimeForm$: BehaviorSubject<FormGroup<DateAndTimeForm>> =
    new BehaviorSubject<FormGroup<DateAndTimeForm>>(
      new FormGroup<DateAndTimeForm>({
        DatePicker: new FormGroup<DatePickerForm>({
          start: new FormControl<Date | null>(null),
          end: new FormControl<Date | null>(null),
        }),
        Start: new FormControl<string | null>(null, {
          validators: [Validators.pattern(timePattern), TimeValidator],
        }),
        End: new FormControl<string | null>(null, [
          Validators.pattern(timePattern),
          TimeValidator,
        ]),
      })
    );
  filterForm: FormGroup<TransactionCoinFilterForm> =
    new FormGroup<TransactionCoinFilterForm>({
      to: new FormControl<nodeDetails | null>(null),
      from: new FormControl<nodeDetails | null>(null),
      coinsMax: new FormControl<string | null>(null, [
        Validators.pattern(/^[0-9]+(\.[0-9]+)?$/),
      ]),
      coinsMin: new FormControl<string | null>(null, [
        Validators.pattern(/^[0-9]+(\.[0-9]+)?$/),
      ]),
      reason: new FormControl<string | null>(null, []),
      DateAndTimeForm: new FormGroup<DateAndTimeForm>({
        DatePicker: new FormGroup<DatePickerForm>({
          start: new FormControl<Date | null>(null),
          end: new FormControl<Date | null>(null),
        }),
        Start: new FormControl<string | null>(null, {
          validators: [Validators.pattern(timePattern), TimeValidator],
        }),
        End: new FormControl<string | null>(null, [
          Validators.pattern(timePattern),
          TimeValidator,
        ]),
      }),
    });
  typingTimer: any;
  applyFilter() {
    clearTimeout(this.typingTimer); // Clear any existing timer

    this.typingTimer = setTimeout(() => {
      this.filterTransactionCoinToaBehavirSubjectService.next(
        this.filter$,
        this.filterForm
      );
    }, 1000);
  }
  notify() {
    this.filterTransactionCoinToaBehavirSubjectService.next(
      this.filter$,
      this.filterForm
    );
  }
}
export   interface TransactionCoinFilterForm{
 coinsMin:FormControl<string |  null>  ; 
 coinsMax:FormControl<string |  null > ; 
 reason:FormControl<string | null>;
 from:FormControl<nodeDetails | null>  ;
 to:FormControl<nodeDetails | null > ; 
 DateAndTimeForm:FormGroup<DateAndTimeForm>

}

     

