import { Pipe, PipeTransform } from '@angular/core';
import { TransactionCoinsRow } from '../../sharable';

@Pipe({
  name: 'sendCoinsPipe',
  standalone: true
})
export class SendCoinspipePipe implements PipeTransform {

  transform(value:TransactionCoinsRow ,indexId: number): string {
    console.log(value ,indexId )
 if (value.From === indexId) {
      return `<span style="color: red">${value.Coins}</span>`;
    } else if  (value.To == indexId){
      return `<span style="color: green">${value.Coins}</span>`;
    }else {
      return `<span>${value.Coins}</span>`;
    }
  }

}
