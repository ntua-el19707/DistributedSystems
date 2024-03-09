import { Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { GraphQLModule } from '../../features/GraphQL/graph-ql.module';
import { AsyncPipe, JsonPipe } from '@angular/common';
import { BehaviorSubject } from 'rxjs';
import { GraphQLResponse } from '../../sharable';
import { GraphQLBehaviorSubjectService } from '../../features/GraphQL/graph-qlbehavior-subject.service';
const material = [MatFormFieldModule, MatInputModule,MatButtonModule]
@Component({
  selector: 'app-graphql-rsp',
  standalone: true,
  imports: [GraphQLModule , ...material ,JsonPipe ,AsyncPipe],
  templateUrl: './graphql-rsp.component.html',
  styleUrl: './graphql-rsp.component.scss'
})
export class GraphqlRspComponent {
  readonly jsonBehaviorSubject$ :BehaviorSubject<GraphQLResponse> 
  constructor(private graphQLBehaviorSubjectService:GraphQLBehaviorSubjectService){
    this.jsonBehaviorSubject$ = this.graphQLBehaviorSubjectService.GetBehaviorSubject()
  }

}
