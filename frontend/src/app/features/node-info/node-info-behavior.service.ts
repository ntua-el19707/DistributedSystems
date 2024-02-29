import { Injectable } from '@angular/core';
import { NodeInfoModule } from './node-info.module';
import { NodeInfoClientService } from './node-info-client.service';
import { nodeDetails, nodeInfoRsp } from '../../sharable';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: NodeInfoModule
})
export class NodeInfoBehaviorService { 
  #nodeInfoBehaviorSubject:BehaviorSubject<nodeDetails> = new BehaviorSubject({nodeId:"" , indexId:0 , uri:"", uriPublic:"string"} )
  #nodeInfoOthersBehaviorSubject:BehaviorSubject<Array<number>> = new BehaviorSubject<Array<number>>([]) 
  constructor(private  nodeInfoClientService:NodeInfoClientService) { }
  fetchNodeInfo(){
    this.nodeInfoClientService.getInfo().subscribe(r =>{
      const node = r.client.indexId
      let  other :Array<number> = [] 
      for(let i=0 ; i<r.total ;i++  ) {
       if(  i !==  node) {
        other.push(i)
      }
      console.log(r.client)
      this.#nodeInfoBehaviorSubject.next(r.client)
      this.#nodeInfoOthersBehaviorSubject.next(other)
      }
    } , err=>{} , ()=>{})

  }
  getNodeInfo():BehaviorSubject<nodeDetails>{
    return this.#nodeInfoBehaviorSubject
  }
}
