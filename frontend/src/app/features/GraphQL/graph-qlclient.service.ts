import { Injectable } from '@angular/core';
import { GraphQLModule } from './graph-ql.module';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { GraphQLResponse } from '../../sharable';

@Injectable({
  providedIn: GraphQLModule
})
export class GraphQLClientService {

   constructor(private  http:HttpClient) {}

  query(query:string) :Observable< GraphQLResponse>{
    const request = `/graphql?query=`+query
    console.log(request)
    return  this.http.get(request) as Observable<GraphQLResponse>
  }
}
