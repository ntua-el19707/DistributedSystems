import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { NodeInfoGraphQL } from '../../sharable';

@Injectable({
  providedIn: 'root'
})
export class AllNodesBehaviorSubjectService {
  readonly #allNodesBheaviorSubjec:BehaviorSubject<Array<NodeInfoGraphQL>> = new BehaviorSubject<Array<NodeInfoGraphQL>>([]) 
  
   getBehaviorSubject():BehaviorSubject<Array<NodeInfoGraphQL>> {
    return this.#allNodesBheaviorSubjec
   }
   next(list : Array<NodeInfoGraphQL>){
    this.#allNodesBheaviorSubjec.next(list)
  }
}
