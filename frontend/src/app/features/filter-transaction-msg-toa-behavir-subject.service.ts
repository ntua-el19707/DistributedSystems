import { Injectable } from '@angular/core';
import { FormGroup } from '@angular/forms';
import { BehaviorSubject } from 'rxjs';
import { TransactionMsgFilterForm } from '../components/filter-transaction-msg/filter-transaction-msg.component';

@Injectable({
  providedIn: 'root'
})
export class FilterTransactionMsgToaBehavirSubjectService {

  constructor() { }
  next(subject :BehaviorSubject<FilterTransactionMsg> , obj :  FormGroup<TransactionMsgFilterForm>){
  if (obj.valid){
    let filter:FilterTransactionMsg ={ }
    if (obj.controls.to.value !== null && obj.controls.to.value !== undefined) {
      if (obj.controls.to.value.indexId !== undefined)
        filter.To = obj.controls.to.value.indexId;
    }
    if (obj.controls.from.value !== null && obj.controls.from.value !== undefined) {
      if (obj.controls.from.value.indexId !== undefined)
        filter.From = obj.controls.from.value.indexId;
    }
    if (obj.controls.message.value !== null ){
      filter.Message = obj.controls.message.value
    }
       
    const  timeform =  obj.controls.DateAndTimeForm 
    if (timeform.valid){
      if( timeform.controls.DatePicker.value.start){
        let time ="00:00:00"
        if  (timeform.controls.Start.value !== "" &&  timeform.controls.Start.value != null){
          time = timeform.controls.Start.value
        }
        const date =  new Date(timeform.controls.DatePicker.value.start)
        const isoTimestamp = new Date(date.getFullYear(), date.getMonth(), date.getDate(),
                                  parseInt(time.split(':')[0]), parseInt(time.split(':')[1]), parseInt(time.split(':')[2])).toISOString();

       
        const unixTime = Date.parse(isoTimestamp) / 1000; 
        filter.SendTimeLess = unixTime
      }
      if( timeform.controls.DatePicker.value.end){
        let time ="23:59:59"
        if  (timeform.controls.End.value !== "" &&  timeform.controls.End.value != null){
          time = timeform.controls.End.value
        }
        const date =  new Date(timeform.controls.DatePicker.value.end)
        const isoTimestamp = new Date(date.getFullYear(), date.getMonth(), date.getDate() , 
                                  parseInt(time.split(':')[0]), parseInt(time.split(':')[1]), parseInt(time.split(':')[2])).toISOString();

        const unixTime = Date.parse(isoTimestamp) / 1000; 
      
        filter.SendTimeMore = unixTime
      }
      subject.next(filter)
    }
    
  }

  }
}
export  interface FilterTransactionMsg {
  To ?: number  ; 
  From?:number ;
  Message?:string ;
  SendTimeLess?:number ; 
  SendTimeMore?:number;  
} 