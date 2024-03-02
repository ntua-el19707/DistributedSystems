import { Component, Input } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { nodeDetails } from '../../sharable';
import { DateAndTimeForm, DatePickerForm, TimeValidator, timePattern } from '../date-and-time-picker/date-and-time-picker.component';
import { BehaviorSubject, Subscription } from 'rxjs';
import { FilterTransactionMsg, FilterTransactionMsgToaBehavirSubjectService } from '../../features/filter-transaction-msg-toa-behavir-subject.service';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatInputModule } from '@angular/material/input';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import {DateAndTimePickerComponent,  } from '../date-and-time-picker/date-and-time-picker.component';


const material = [MatFormFieldModule, MatSelectModule, MatInputModule,MatCardModule ,MatButtonModule]
const forms = [FormsModule, ReactiveFormsModule]
@Component({
  selector: 'app-filter-transaction-msg',
  standalone: true,
  imports: [...material ,  ...forms ,DateAndTimePickerComponent],
  templateUrl: './filter-transaction-msg.component.html',
  styleUrl: './filter-transaction-msg.component.scss'
})
export class FilterTransactionMsgComponent {
  clients:Array<nodeDetails > = []
@Input() from : boolean = false  
@Input() to  : boolean = true 
 @Input() filter$ :BehaviorSubject<FilterTransactionMsg> = new  BehaviorSubject({})
readonly #Subscription:Subscription = new Subscription()
 
  constructor(private  filterTransactionMsgToaBehavirSubjectService: FilterTransactionMsgToaBehavirSubjectService){
    this.#Subscription =  this.DateAndTimeForm$.subscribe(r=>{
      this.filterForm.controls.DateAndTimeForm = r  
      this.filterTransactionMsgToaBehavirSubjectService.next( this.filter$,  this.filterForm)
    })
  }
  ngOnDestroy(): void {
    this.#Subscription.unsubscribe()
  }
   DateAndTimeForm$ :BehaviorSubject< FormGroup<DateAndTimeForm> >= new  BehaviorSubject<FormGroup<DateAndTimeForm>> (new FormGroup<DateAndTimeForm>({
  DatePicker: new FormGroup<DatePickerForm>({
    start: new FormControl<Date |  null>(null),
    end: new FormControl<Date |  null>(null),
    }) ,
Start: new FormControl<string | null>(null, { validators:[Validators.pattern(timePattern)   ,TimeValidator]}),
  End: new FormControl<string | null>(null, [Validators.pattern(timePattern) ,TimeValidator]),
}))
  filterForm: FormGroup<TransactionMsgFilterForm> = 
 new FormGroup<TransactionMsgFilterForm>({
  to:new FormControl<nodeDetails | null>(null ) ,
  from:new FormControl<nodeDetails | null>(null ) , 
 
  message:new FormControl<string | null>(null,[]) ,
  DateAndTimeForm:new  FormGroup<DateAndTimeForm>({DatePicker: new FormGroup<DatePickerForm>({
    start: new FormControl<Date |  null>(null),
    end: new FormControl<Date |  null>(null),
    }) ,
Start: new FormControl<string | null>(null, { validators:[Validators.pattern(timePattern)   ,TimeValidator]}),
  End: new FormControl<string | null>(null, [Validators.pattern(timePattern) ,TimeValidator]),}
  )
})
typingTimer:any
applyFilter(){
    clearTimeout(this.typingTimer); // Clear any existing timer
   
    this.typingTimer = setTimeout(() => { 
   this.filterTransactionMsgToaBehavirSubjectService.next( this.filter$,  this.filterForm)
    
  
    }, 1000);
}
notify(){
  this.filterTransactionMsgToaBehavirSubjectService.next( this.filter$,  this.filterForm)
}
}

export   interface TransactionMsgFilterForm{
 message:FormControl<string | null>;
 from:FormControl<nodeDetails | null>  ;
 to:FormControl<nodeDetails | null > ; 
 DateAndTimeForm:FormGroup<DateAndTimeForm>

}
