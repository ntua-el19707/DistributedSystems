import { Injectable } from '@angular/core';
import { GraphQLModule } from './graph-ql.module';
import { BehaviorSubject } from 'rxjs';
import { GraphQLResponse } from '../../sharable';
import { GraphQLClientService } from './graph-qlclient.service';

@Injectable({
  providedIn: GraphQLModule
})
export class GraphQLBehaviorSubjectService {
  
  #dataPlayGrounndBehaviorSubject:BehaviorSubject<GraphQLResponse> = new BehaviorSubject<GraphQLResponse>({})
  constructor(private  graphQLClientService:GraphQLClientService){}
  fetchData(query:string , happening$:BehaviorSubject<boolean>){
    happening$.next(true)
    this.graphQLClientService.query(query).subscribe(r =>{
      console.log(r)
      this.#dataPlayGrounndBehaviorSubject.next(r)} , err=> {
      console.log(err)
      this.#dataPlayGrounndBehaviorSubject.next(err)} ,()=>{happening$.next(false)})
  }
  GetBehaviorSubject():BehaviorSubject<GraphQLResponse>{
    return this.#dataPlayGrounndBehaviorSubject
  }
 
}

