import { Injectable } from '@angular/core';
import { TransactionsMsgClientModule } from './transactions-msg-client.module';
import { TransactionMsgClientService } from './transaction-msg-client.service';
import { BehaviorSubject, Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import {
  GraphQLResponse,
  TransactionMsgList,
  TransactionMsgRowGraphQL,
} from '../../sharable';
import { FilterTransactionMsg } from '../filter-transaction-msg-toa-behavir-subject.service';
import { AllNodesBehaviorSubjectService } from '../allNodes/all-nodes-behavior-subject.service';

@Injectable({
  providedIn: TransactionsMsgClientModule,
})
export class TransactionMsgBehaviorService {
  #inboxBehaviorSubject: BehaviorSubject<TransactionMsgList> =
    new BehaviorSubject<TransactionMsgList>({
      transactions: [],
      nodeDetails: { indexId: -1 },
    });
  #inboxBehaviorTransactionsSubject: BehaviorSubject<
    Array<TransactionMsgRowGraphQL>
  > = new BehaviorSubject<Array<TransactionMsgRowGraphQL>>([]);
  #sendBehaviorSubject: BehaviorSubject<TransactionMsgList> =
    new BehaviorSubject<TransactionMsgList>({
      transactions: [],
      nodeDetails: { indexId: -1 },
    });
  #sendBehaviorTransactionsSubject: BehaviorSubject<
    Array<TransactionMsgRowGraphQL>
  > = new BehaviorSubject<Array<TransactionMsgRowGraphQL>>([]);
  #allBehaviorSubject: BehaviorSubject<TransactionMsgList> =
    new BehaviorSubject<TransactionMsgList>({
      transactions: [],
      nodeDetails: { indexId: -1 },
    });
  #allBehaviorTransactionsSubject: BehaviorSubject<
    Array<TransactionMsgRowGraphQL>
  > = new BehaviorSubject<Array<TransactionMsgRowGraphQL>>([]);
  #Filter:FilterTransactionMsg = {}
  constructor(
    private transactionMsgClientService: TransactionMsgClientService,
    private allNodesBehaviorSubjectService: AllNodesBehaviorSubjectService
  ) {}
  fetchSend(mode:boolean) {
    this.transactionMsgClientService.getSend(mode).subscribe(
      (r: GraphQLResponse) => {
        console.log(r);
        const transactions = r.data?.send;
        const self = r.data?.self;
        if (transactions && self?.client) {
          this.#sendBehaviorSubject.next({
            transactions: transactions.transactions,
            nodeDetails: self.client,
          });
          this.filterSend(this.#Filter)
        }
           const nodes = r.data?.allNodes;
           if (nodes) {
             this.allNodesBehaviorSubjectService.next(nodes);
           }
      },
      (err) => {
        console.log(err);
      },
      () => {}
    );
  }
  getSend(): BehaviorSubject<TransactionMsgList> {
    return this.#sendBehaviorSubject;
  }
  filterSend(filter: FilterTransactionMsg) {
    const observable = this.FilterBehaviorSubject(
      filter,
      this.#sendBehaviorSubject
    );
    observable.subscribe((r: TransactionMsgList) => {
      this.#sendBehaviorTransactionsSubject.next(r.transactions);
    });
  }
  getSendTransactions(): BehaviorSubject<Array<TransactionMsgRowGraphQL>> {
    return this.#sendBehaviorTransactionsSubject;
  }
  fetchAll(mode: boolean) {
    this.transactionMsgClientService.getAll(mode).subscribe(
      (r: GraphQLResponse) => {
        const transactions = r.data?.allTransactionMsg;

        const self = r.data?.self;
        if (transactions && self?.client) {
          this.#allBehaviorSubject.next({
            transactions: transactions.transactions,
            nodeDetails: self.client,
          });
        
            
          this.filterAll(this.#Filter)
        }
        const nodes = r.data?.allNodes;
        if (nodes) {
          this.allNodesBehaviorSubjectService.next(nodes);
        }
      },
      (err) => {
        console.log(err);
      },
      () => {}
    );
  }
  getAll(): BehaviorSubject<TransactionMsgList> {
    return this.#allBehaviorSubject;
  }
  filterAll(filter: FilterTransactionMsg) {
    const observable = this.FilterBehaviorSubject(
      filter,
      this.#allBehaviorSubject
    );
    observable.subscribe((r: TransactionMsgList) => {
      this.#allBehaviorTransactionsSubject.next(r.transactions);
    });
  }
  getAllTransactions(): BehaviorSubject<Array<TransactionMsgRowGraphQL>> {
    return this.#allBehaviorTransactionsSubject;
  }
  fetchInbox(mode:boolean) {
    this.transactionMsgClientService.getINBOX(mode).subscribe(
      (r: GraphQLResponse) => {
        const transactions = r.data?.inbox;
        const self = r.data?.self;
        if (transactions && self?.client) {
          this.#inboxBehaviorSubject.next({
            transactions: transactions.transactions,
            nodeDetails: self.client,
          });
         this.filter(this.#Filter)
        }
         const nodes = r.data?.allNodes;
         if (nodes) {
           this.allNodesBehaviorSubjectService.next(nodes);
         }
      },
      (err) => {
        console.log(err);
      },
      () => {}
    );
  }
  getInbox(): BehaviorSubject<TransactionMsgList> {
    return this.#inboxBehaviorSubject;
  }
  filter(filter: FilterTransactionMsg) {
    const observable = this.FilterBehaviorSubject(
      filter,
      this.#inboxBehaviorSubject
    );
    observable.subscribe((r: TransactionMsgList) => {
      this.#inboxBehaviorTransactionsSubject.next(r.transactions);
    });
  }
  getInboxTransactions(): BehaviorSubject<Array<TransactionMsgRowGraphQL>> {
    return this.#inboxBehaviorTransactionsSubject;
  }
  FilterBehaviorSubject(
    filter: FilterTransactionMsg,
    subject: BehaviorSubject<TransactionMsgList>
  ): Observable<TransactionMsgList> {
    this.#Filter = filter
    const filterComp: (
      row: TransactionMsgRowGraphQL,
      filter: FilterTransactionMsg
    ) => boolean = (
      row: TransactionMsgRowGraphQL,
      filter: FilterTransactionMsg
    ): boolean => {
      let defaultRsp = true;
      
      if (row.To !== undefined) {
        if (filter.To !== undefined){
        defaultRsp = defaultRsp && ( row.To === filter.To)};
      }

      if (row.Msg !== undefined) {
        defaultRsp =
          defaultRsp && (!filter.Message || row.Msg.includes(filter.Message));
      }

      if (row.From !== undefined) {
           if (filter.From !== undefined){
        defaultRsp = defaultRsp && (row.From === filter.From)}
        
      }

      if (row.Time !== undefined) {
        defaultRsp =
          defaultRsp &&
          (!filter.SendTimeLess || row.Time >= filter.SendTimeLess);
        defaultRsp =
          defaultRsp &&
          (!filter.SendTimeMore || row.Time < filter.SendTimeMore);
      }

      return defaultRsp;
    };

    const observable = subject.pipe(
      map((response: TransactionMsgList) => {
        const filteredTransactions = response.transactions.filter(
          (transaction) => {
            return filterComp(transaction, filter);
          }
        );
        return { ...response, transactions: filteredTransactions };
      })
    );
    return observable;
  }
}
