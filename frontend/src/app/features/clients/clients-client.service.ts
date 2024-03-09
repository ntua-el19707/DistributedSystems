import { Injectable } from '@angular/core';
import { ClientsModule } from './clients.module';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { GraphQLResponse, clientsInfoRsp } from '../../sharable';
import { GraphQLClientService } from '../GraphQL/graph-qlclient.service';

@Injectable({
  providedIn: ClientsModule,
})
export class ClientsClientService {
  constructor(
    private graphQLClientService: GraphQLClientService
  ) {}
  getClients(): Observable<GraphQLResponse> {
    const query = `
     {
       clients{
            nodeId,
            indexId,
            uri , 
            uriPublic, 
        }
     }
     `;
    return this.graphQLClientService.query(query) as Observable<GraphQLResponse>
  }
}
