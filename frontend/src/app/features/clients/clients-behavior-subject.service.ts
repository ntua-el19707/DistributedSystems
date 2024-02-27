import { Injectable } from '@angular/core';
import { ClientsModule } from './clients.module';
import { ClientsClientService } from './clients-client.service';
import { clientsInfoRsp, nodeDetails } from '../../sharable';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: ClientsModule
})
export class ClientsBehaviorSubjectService {

   #ClientsBehaviorSubject:BehaviorSubject<Array<nodeDetails>> = 
   new BehaviorSubject<Array<nodeDetails>>([])
  constructor(private  clientsClientService:ClientsClientService) { }
  fetchClients(){

    this.clientsClientService.getClients().subscribe(r=>{
      this.#ClientsBehaviorSubject.next(r.clients)
    } , err=>{} , ()=>{})
  } 
  getBehavior():BehaviorSubject<Array<nodeDetails>> {
    return this.#ClientsBehaviorSubject
  }
}
