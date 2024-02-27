import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

import { NavBarComponent } from './components/nav-bar/nav-bar.component';
import { TranferCoinsFormComponent } from './components/tranfer-coins-form/tranfer-coins-form.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule,NavBarComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  title = 'frontend';
}
