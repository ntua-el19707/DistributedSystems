import { Component, Input } from '@angular/core';
import {FormGroup, FormControl, FormsModule, ReactiveFormsModule, Validators, AbstractControl, ValidatorFn, ValidationErrors, Validator} from '@angular/forms';
import {provideNativeDateAdapter} from '@angular/material/core';
import {MatDatepickerModule} from '@angular/material/datepicker';
import {MatFormFieldModule} from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { BehaviorSubject } from 'rxjs';

const material = [MatFormFieldModule, MatInputModule ,MatDatepickerModule]
const forms = [FormsModule, ReactiveFormsModule]
@Component({
  selector: 'app-date-and-time-picker',
  standalone: true,
  providers:[provideNativeDateAdapter()],
  imports: [...forms , material],
  templateUrl: './date-and-time-picker.component.html',
  styleUrl: './date-and-time-picker.component.scss'
})
export class DateAndTimePickerComponent {

@Input () DateAndTimeForm$ :BehaviorSubject< FormGroup<DateAndTimeForm> >= new  BehaviorSubject<FormGroup<DateAndTimeForm>> (new FormGroup<DateAndTimeForm>({
  DatePicker: new FormGroup<DatePickerForm>({
    start: new FormControl<Date |  null>(null),
    end: new FormControl<Date |  null>(null),
    }) ,
Start: new FormControl<string | null>(null, { validators:[Validators.pattern(timePattern)   ,TimeValidator]}),
  End: new FormControl<string | null>(null, [Validators.pattern(timePattern) ,TimeValidator]),

}))
DateAndTimeForm : FormGroup<DateAndTimeForm> = new FormGroup<DateAndTimeForm>({
  DatePicker: new FormGroup<DatePickerForm>({
    start: new FormControl<Date |  null>(null),
    end: new FormControl<Date |  null>(null),
    }) ,
Start: new FormControl<string | null>(null, { validators:[Validators.pattern(timePattern)   ,TimeValidator]}),
  End: new FormControl<string | null>(null, [Validators.pattern(timePattern) ,TimeValidator]),

})
typingTimer:any
applyFilter(){
    clearTimeout(this.typingTimer); // Clear any existing timer
   
    this.typingTimer = setTimeout(() => { 
     this.DateAndTimeForm$.next(this.DateAndTimeForm)
    
  
    }, 1000);
}
notify(){
  this.DateAndTimeForm$.next(this.DateAndTimeForm)
}
}
export interface DatePickerForm {
  start:FormControl<Date  | null > ,
  end:FormControl<Date | null > 
}
export interface DateAndTimeForm {
  DatePicker:FormGroup<DatePickerForm> ;
  Start:FormControl<string | null> ;
  End: FormControl<string |  null  > ;

} 
export const timePattern = /^(?:[01]\d|2[0-3]):(?:[0-5]\d):(?:[0-5]\d)$/;


export  const  TimeValidator:  ValidatorFn  = (control:AbstractControl):  ValidationErrors|  null  =>{
 const value = control.value;
  if (!value) {
    return null; 
  }
  if (!timePattern.test(value)) {
    return { invalidFormat: true };
  }
  const [hours, minutes, seconds] = value.split(':').map(Number);
  if (hours >= 24 || minutes >= 60 || seconds >= 60) {
    return { invalidTime: true };
  }
  return null;
 
}