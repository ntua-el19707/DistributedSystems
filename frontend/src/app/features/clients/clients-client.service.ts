import { Injectable } from '@angular/core';
import { ClientsModule } from './clients.module';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { clientsInfoRsp } from '../../sharable';

@Injectable({
  providedIn: ClientsModule
})
export class ClientsClientService {
  private  uri:string = "/api/v1/NodeDetails/Clients"
  constructor(private  http:HttpClient) { }
  getClients():Observable<clientsInfoRsp>{

     return  this.http.get(this.uri) as  Observable<clientsInfoRsp>
  }
}
