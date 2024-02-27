import { Routes } from '@angular/router';
import { HomePageComponent } from './pages/home-page/home-page.component';

import { InboxComponent } from './pages/inbox/inbox.component';
import { SendComponent } from './pages/send/send.component';
import { AllMsgComponent } from './pages/all-msg/all-msg.component';
import { TrasferCoinsComponent } from './pages/trasfer-coins/trasfer-coins.component';
import { SendMessageFormPageComponent } from './pages/send-message-form-page/send-message-form-page.component';

export const routes: Routes = [
{path:'' , component:HomePageComponent},
{path:'inbox' , component:InboxComponent},
{path:'send' , component:SendComponent},
{path:'allMessages' , component:AllMsgComponent},
{path:'transfer' , component:TrasferCoinsComponent},
{path:'sendMessage' , component:SendMessageFormPageComponent}
];
