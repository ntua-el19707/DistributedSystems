import { Component } from '@angular/core';
import { FormStakeComponent as  fc } from '../../components/form-stake/form-stake.component';
@Component({
  selector: 'app-stake-form-page',
  standalone: true,
  imports: [fc],
  templateUrl: './form-stake.component.html',
  styleUrl: './form-stake.component.scss'
})
export class FormStakeComponent {

}
