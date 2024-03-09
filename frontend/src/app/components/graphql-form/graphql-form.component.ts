import { Component } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { BehaviorSubject } from 'rxjs';
import { GraphQLModule } from '../../features/GraphQL/graph-ql.module';
import { GraphQLBehaviorSubjectService } from '../../features/GraphQL/graph-qlbehavior-subject.service';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';
import { CommonModule } from '@angular/common';
import { TabDirectiveDirective } from '../../features/Directives/tab-directive.directive';
const material = [MatProgressSpinnerModule ,MatFormFieldModule, MatInputModule,MatCardModule ,MatButtonModule]
const forms = [ FormsModule, ReactiveFormsModule,]
@Component({
  selector: 'app-graphql-form',
  standalone: true,
  imports: [GraphQLModule ,TabDirectiveDirective,...material , ...forms , CommonModule],
  templateUrl: './graphql-form.component.html',
  styleUrl: './graphql-form.component.scss'
})
export class GraphqlFormComponent {
  constructor(private graphQLBehaviorSubjectService:GraphQLBehaviorSubjectService){}
  GraphQLForm:FormGroup<GraphQLFormInterface> =  new  FormGroup<GraphQLFormInterface> ({
  query:new FormControl<string |  null >(null ,  Validators.required)
  })
  happening$:BehaviorSubject< boolean>  = new BehaviorSubject<boolean>(false)
  submit() {
    const isHappening = this.happening$.getValue()

  if (isHappening) {
    return; 
  }
    if (this.GraphQLForm.valid){
      const  query =  this.GraphQLForm.controls.query.value
      if  (query){
        console.log(query)
        this.graphQLBehaviorSubjectService.fetchData(query , this.happening$)
      }
    }
  }
}

interface  GraphQLFormInterface {
  query: FormControl<string| null >
}
