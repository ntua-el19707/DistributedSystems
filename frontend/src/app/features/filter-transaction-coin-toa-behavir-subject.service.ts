import { Injectable } from '@angular/core';
import { BehaviorSubject, ignoreElements } from 'rxjs';
import { TransactionCoinFilterForm } from '../components/filter-transaction-coins/filter-transaction-coins.component';
import { FormGroup } from '@angular/forms';

@Injectable({
  providedIn: 'root'
})
export class FilterTransactionCoinToaBehavirSubjectService {

  constructor() { }
  next(subject :BehaviorSubject<filterTranctionCoinSubject> , obj :  FormGroup<TransactionCoinFilterForm>){
  console.log(obj)
    if (obj.valid){
    let filter:filterTranctionCoinSubject ={ }
    if (obj.controls.to.value !== null  && obj.controls.to.value !== undefined) {
       if (obj.controls.to.value.indexId !== undefined)
         
         filter.To = obj.controls.to.value.indexId; 
    }
    if (
      obj.controls.from.value !== null &&
      obj.controls.from.value !== undefined
    ) {
      if (obj.controls.from.value.indexId !== undefined)
 
        filter.From = obj.controls.from.value.indexId;
    }
    if (obj.controls.reason.value !== null ){
    
      filter.Reason = obj.controls.reason.value
    }
        if (obj.controls.coinsMax.value !== null) {
      filter.CoinsMax =  parseInt(obj.controls.coinsMax.value,10) 
    }  
    if (obj.controls.coinsMin.value !== null) {
      filter.CoinsMin =  parseInt(obj.controls.coinsMin.value,10) 
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
export interface filterTranctionCoinSubject{
  To?: number ; 
  From?: number ; 
  CoinsMin?:number ;
  CoinsMax?:number ;
  Reason?:string ;
  SendTimeLess?:number ; 
  SendTimeMore?:number;  
}