import { Routes } from '@angular/router';
import { HomePageComponent } from './pages/home-page/home-page.component';
import { MessengerPageComponent } from './pages/messenger-page/messenger-page.component';

export const routes: Routes = [
{path:'' , component:HomePageComponent},
{path:'messenger' , component:MessengerPageComponent}

];
