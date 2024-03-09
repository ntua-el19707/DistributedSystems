import { Injectable } from '@angular/core';
import { NodeInfoModule } from './node-info.module';

import { Observable } from 'rxjs';
import { GraphQLResponse} from '../../sharable';
import { GraphQLClientService } from '../GraphQL/graph-qlclient.service';

@Injectable({
  providedIn: NodeInfoModule,
})
export class NodeInfoClientService {
  constructor(
    private graphQLClientService: GraphQLClientService
  ) {}
  getInfo(): Observable<GraphQLResponse> {
    const query = `
     {
       self{
          client{  
            nodeId,
            indexId,
            uri , 
            uriPublic
          }
          total
        }
     }
     `;
    return this.graphQLClientService.query(query) as Observable<GraphQLResponse>
  }
}
