import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { GraphQLModule } from '../GraphQL/graph-ql.module';



@NgModule({
  declarations: [],
  imports: [
    HttpClientModule ,GraphQLModule
  ]
})
export class TransactionsMsgClientModule { }
