import { Directive, HostListener } from '@angular/core';

@Directive({
  selector: '[appTabDirective]',
  standalone: true
})
export class TabDirectiveDirective {

  @HostListener('keydown', ['$event'])
  onKeyDown(event: KeyboardEvent) {
    if (event.key === 'Tab') {
      event.preventDefault();
      const { selectionStart, selectionEnd, value } = event.target as HTMLTextAreaElement;
      const newValue = value.substring(0, selectionStart) + '\t' + value.substring(selectionEnd);
      (event.target as HTMLTextAreaElement).value = newValue;
  
      (event.target as HTMLTextAreaElement).setSelectionRange(selectionStart + 1, selectionStart + 1);
    }
  }

}
