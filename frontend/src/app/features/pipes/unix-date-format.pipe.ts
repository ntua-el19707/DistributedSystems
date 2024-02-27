import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'unixDateFormat',
  standalone: true
})
export class UnixDateFormatPipe implements PipeTransform {

  transform(unixTime: number): string  {
    if ( !unixTime )return ''
    const date  =   new Date(unixTime * 1000)
    return date.toLocaleString();
  }

}
