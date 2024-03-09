import { Component } from '@angular/core';
import { GraphqlRspComponent } from '../../components/graphql-rsp/graphql-rsp.component';
import { GraphqlFormComponent } from '../../components/graphql-form/graphql-form.component';

@Component({
  selector: 'app-playground',
  standalone: true,
  imports: [GraphqlFormComponent , GraphqlRspComponent],
  templateUrl: './playground.component.html',
  styleUrl: './playground.component.scss'
})
export class PlaygroundComponent {

}
