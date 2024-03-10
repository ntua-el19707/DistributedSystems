import { Injectable } from '@angular/core';
import { TransactionCoinsModule } from './transaction-coins.module';
import { HttpClient } from '@angular/common/http';
import {
  BalanceRsp,
  GraphQLResponse,
  transactionCoinRequest,
} from '../../sharable';
import { BehaviorSubject, Observable } from 'rxjs';
import { OpenErrordialogService } from '../open-errordialog.service';
import { NavigateService } from '../navigate.service';
import { GraphQLClientService } from '../GraphQL/graph-qlclient.service';

@Injectable({
  providedIn: TransactionCoinsModule,
})
export class TransactionClientService {
  constructor(
    private http: HttpClient,
    private graphQLClientService: GraphQLClientService,
    private navigateService: NavigateService,
    private openErrordialogService: OpenErrordialogService
  ) {}
  getBalance(): Observable<GraphQLResponse> {
    const query = `{balance {availableBalance}}`;
    return this.graphQLClientService.query(
      query
    ) as Observable<GraphQLResponse>;
  }

  getMyTransactions(mode: boolean): Observable<GraphQLResponse> {
    let query = `{
  nodeTransactions {
      Transactions{
        From , 
        To , 
        Coins , 
        Reason , 
        Nonce  , 
        Time , TransactionId
      }
  }
   self{
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

  getAllTransactions(mode: boolean): Observable<GraphQLResponse> {
    let query = `{
  getTransactionsCoins{
      Transactions{
        From , 
        To , 
        Coins , 
        Reason , 
        Nonce  , 
        Time ,TransactionId
      }
  }
   self{
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
    to: number,
    coins: number,
    happening$: BehaviorSubject<boolean>
  ) {
    happening$.next(true);
    const body: transactionCoinRequest = {
      to: to,
      amount: coins,
    };
    this.http.post('/api/v1/transfer', body).subscribe(
      (r) => {
        this.navigateService.navigateTo('/');
      },
      (err) => {
        this.openErrordialogService.errorDialog(err.error.Message);
        happening$.next(false);
      },
      () => {}
    );
  }
  postStake(
    coins: number,
    happening$: BehaviorSubject<boolean>
  ) {
    happening$.next(true);
    const body: {stake:number} = {
      stake: coins,
    };
    this.http.post('/api/v1/stake', body).subscribe(
      (r) => {
        this.navigateService.navigateTo('/');
      },
      (err) => {
        this.openErrordialogService.errorDialog(err.error.Message);
        happening$.next(false);
      },
      () => {}
    );
  }
}
