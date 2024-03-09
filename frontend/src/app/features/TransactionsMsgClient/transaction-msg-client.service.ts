import { Injectable } from '@angular/core';
import { TransactionsMsgClientModule } from './transactions-msg-client.module';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import {
  GraphQLResponse,
  transactionMsgRequest,
  transactionMsgResponse,
} from '../../sharable';
import { OpenErrordialogService } from '../open-errordialog.service';
import { NavigateService } from '../navigate.service';
import { GraphQLClientService } from '../GraphQL/graph-qlclient.service';

@Injectable({
  providedIn: TransactionsMsgClientModule,
})
export class TransactionMsgClientService {
  constructor(
    private http: HttpClient,
    private graphQLClientService: GraphQLClientService,
    private navigateService: NavigateService,
    private openErrordialogService: OpenErrordialogService
  ) {}
  getINBOX(mode: boolean): Observable<GraphQLResponse> {
    let query = `{
    inbox{
        transactions{
            From , 
            To ,
            Msg , 
            Time , 
            Nonce, 
            TransactionId
        }
    },  self{
          client{  
            nodeId,
            indexId,
            uri , 
            uriPublic
          }
    }
`;
    if (mode) {
      query += `allNodes{
        nodeId,
            indexId,
            uri , 
            uriPublic
  }}`;
    } else {
      query += '}';
    }
    return this.graphQLClientService.query(
      query
    ) as Observable<GraphQLResponse>;
  }
  getSend(mode: boolean): Observable<GraphQLResponse> {
    let query = `{
    send{
        transactions{
            From , 
            To ,
            Msg , 
            Time , 
            Nonce, 
            TransactionId
        }
    },  self{
          client{  
            nodeId,
            indexId,
            uri , 
            uriPublic
          }
    }
`;
    if (mode) {
      query += `allNodes{
        nodeId,
            indexId,
            uri , 
            uriPublic
  }}`;
    } else {
      query += '}';
    }
    return this.graphQLClientService.query(
      query
    ) as Observable<GraphQLResponse>;
  }
  getAll(mode: boolean): Observable<GraphQLResponse> {
    let query = `{
    allTransactionMsg{
        transactions{
            From , 
            To ,
            Msg , 
            Time ,  
            Nonce, 
            TransactionId
        } 
    },  self{
          client{  
            nodeId,
            indexId,
            uri , 
            uriPublic
          }
    }
`;
    if (mode) {
      query += `allNodes{
        nodeId,
            indexId,
            uri , 
            uriPublic
  }}`;
    } else {
      query += '}';
    }
    return this.graphQLClientService.query(
      query
    ) as Observable<GraphQLResponse>;
  }
  postTransaction(
    To: number,
    Msg: string,
    happening: BehaviorSubject<boolean>
  ) {
    happening.next(true);
    const body: transactionMsgRequest = {
      to: To,
      msg: Msg,
    };
    this.http.post('/api/v1/send', body).subscribe(
      (r) => {
        this.navigateService.navigateTo('/send');
      },
      (err) => {
        this.openErrordialogService.errorDialog(err.error.Message);
        happening.next(false);
      },
      () => {}
    );
  }
}
