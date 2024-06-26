import { Component, inject } from '@angular/core';
import { BreakpointObserver, Breakpoints } from '@angular/cdk/layout';
import { AsyncPipe, CommonModule, JsonPipe } from '@angular/common';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';
import { BehaviorSubject, Observable } from 'rxjs';
import { map, shareReplay } from 'rxjs/operators';
import { RouterOutlet } from '@angular/router';
import { ClientsModule } from '../../features/clients/clients.module';
import { ClientsBehaviorSubjectService } from '../../features/clients/clients-behavior-subject.service';
import {  NodeInfoGraphQL } from '../../sharable';

import {MatExpansionModule} from '@angular/material/expansion';
const custom = [ClientsModule]
const material = [MatExpansionModule, MatIconModule , MatButtonModule , MatToolbarModule,MatSidenavModule,MatListModule]
@Component({
  selector: 'app-nav-bar',
  templateUrl: './nav-bar.component.html',
  styleUrl: './nav-bar.component.scss',
  standalone: true,
  imports: [
    ...custom ,
    AsyncPipe,
    RouterOutlet,
    CommonModule,ClientsModule , JsonPipe ,...material

  ],
})
export class NavBarComponent {
  readonly dataSource$:BehaviorSubject<Array<NodeInfoGraphQL>>
  private breakpointObserver = inject(BreakpointObserver);
constructor(private clientsBehaviorSubjectService:ClientsBehaviorSubjectService){
  this.clientsBehaviorSubjectService.fetchClients()
this.dataSource$ = this.clientsBehaviorSubjectService.getBehavior()
}
  isHandset$: Observable<boolean> = this.breakpointObserver
    .observe(Breakpoints.Handset)
    .pipe(
      map((result) => result.matches),
      shareReplay()
    );
    openInNewTab(link :string | undefined){
      if(link){
        window.open(link, '_blank');}
    }
}
