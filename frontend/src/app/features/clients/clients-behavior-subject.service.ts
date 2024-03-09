import { Injectable } from '@angular/core';
import { ClientsModule } from './clients.module';
import { ClientsClientService } from './clients-client.service';
import { clientsInfoRsp, nodeDetails ,NodeInfoGraphQL} from '../../sharable';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: ClientsModule
})
export class ClientsBehaviorSubjectService {

   #ClientsBehaviorSubject:BehaviorSubject<Array<NodeInfoGraphQL>> = 
   
   new BehaviorSubject<Array<NodeInfoGraphQL>>([])
  constructor(private  clientsClientService:ClientsClientService) { }
  fetchClients(){
    this.clientsClientService.getClients().subscribe(r=>{
      const data = r.data?.clients;
      if (data){
      this.#ClientsBehaviorSubject.next(data)}
    } , err=>{} , ()=>{})
  } 
  getBehavior():BehaviorSubject<Array<NodeInfoGraphQL>> {
    return this.#ClientsBehaviorSubject
  }
}
