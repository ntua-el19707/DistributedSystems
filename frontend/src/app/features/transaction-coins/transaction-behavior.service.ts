import { Injectable } from '@angular/core';
import { TransactionCoinsModule } from './transaction-coins.module';

import { BehaviorSubject, Observable, map } from 'rxjs';
import {
  BalanceRsp,
  GraphQLResponse,
  TransactionCoinsAndNodeDetails,
  TransactionCoinsRowGraphQL,
} from '../../sharable';
import { TransactionClientService } from './transaction-client.service';
import { filterTranctionCoinSubject } from '../filter-transaction-coin-toa-behavir-subject.service';
import { AllNodesBehaviorSubjectService } from '../allNodes/all-nodes-behavior-subject.service';

@Injectable({
  providedIn: TransactionCoinsModule,
})
export class TransactionBehaviorService {
  #balanceBehaviorSubject: BehaviorSubject<number> =
    new BehaviorSubject<number>(0.0);
  #myTransactionsBehaviorSubject: BehaviorSubject<TransactionCoinsAndNodeDetails> =
    new BehaviorSubject<TransactionCoinsAndNodeDetails>({
      transactions: [],
      nodeDetails: { indexId: -1, nodeId: '', uri: '', uriPublic: '' },
    });
  #myTransactionsfilterBehaviorSubject: BehaviorSubject<TransactionCoinsAndNodeDetails> =
    new BehaviorSubject<TransactionCoinsAndNodeDetails>({
      transactions: [],
      nodeDetails: { indexId: -1, nodeId: '', uri: '', uriPublic: '' },
    });

  #allTransactionsBehaviorSubject: BehaviorSubject<TransactionCoinsAndNodeDetails> =
    new BehaviorSubject<TransactionCoinsAndNodeDetails>({
      transactions: [],
      nodeDetails: { indexId: -1, nodeId: '', uri: '', uriPublic: '' },
    });
  #allTransactionsfilterBehaviorSubject: BehaviorSubject<TransactionCoinsAndNodeDetails> =
    new BehaviorSubject<TransactionCoinsAndNodeDetails>({
      transactions: [],
      nodeDetails: { indexId: -1, nodeId: '', uri: '', uriPublic: '' },
    });
#Filter:filterTranctionCoinSubject ={}
  constructor(
    private transactionClientService: TransactionClientService,
    private allNodesBehaviorSubjectService: AllNodesBehaviorSubjectService
  ) {}
  fetchBalance() {
    this.transactionClientService.getBalance().subscribe(
      (r: GraphQLResponse) => {
        const b = r.data?.balance;
        if (b) {
          if (b.availableBalance) {
            this.#balanceBehaviorSubject.next(b.availableBalance);
          }
        }
      },
      (err) => {},
      () => {}
    );
  }

  getBalanceSubject(): BehaviorSubject<number> {
    return this.#balanceBehaviorSubject;
  }

  fetchMyTransactions(mode: boolean) {
    this.transactionClientService.getMyTransactions(mode).subscribe(
      (r) => {
        const transactions = r.data?.nodeTransactions;
        const client = r.data?.self?.client;
        if (transactions && client) {
          if (transactions.Transactions) {
            this.#myTransactionsBehaviorSubject.next({
              transactions: transactions.Transactions,
              nodeDetails: client,
            });
           
           this.filterCoins(this.#Filter)
          }
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
  getMyFilterTransactions(): BehaviorSubject<TransactionCoinsAndNodeDetails> {
    return this.#myTransactionsfilterBehaviorSubject;
  }
  filterCoins(filter: filterTranctionCoinSubject) {

    const observable = this.filter(filter, this.#myTransactionsBehaviorSubject);

    observable.subscribe((r: TransactionCoinsAndNodeDetails) => {
      this.#myTransactionsfilterBehaviorSubject.next(r);
    });
  }
  private filter(
    filter: filterTranctionCoinSubject,
    subject: BehaviorSubject<TransactionCoinsAndNodeDetails>
  ): Observable<TransactionCoinsAndNodeDetails> {
    console.log(filter)
    const filterComp: (
      row: TransactionCoinsRowGraphQL,
      filter: filterTranctionCoinSubject
    ) => boolean = (
      row: TransactionCoinsRowGraphQL,
      filter: filterTranctionCoinSubject
    ): boolean => {
      let defaultRsp = true;

      if (row.To !== undefined) {
        if  (filter.To !== undefined){
        defaultRsp = defaultRsp && (row.To === filter.To);
      }
      }

      if (row.Reason !== undefined) {
        defaultRsp =
          defaultRsp && (!filter.Reason || row.Reason === filter.Reason);
      }

      if (row.From !== undefined) {
              if  (filter.From !== undefined){
        defaultRsp = defaultRsp && ( row.From === filter.From)};
      }

      if (row.Coins !== undefined) {
        defaultRsp =
          defaultRsp && (!filter.CoinsMin || row.Coins >= filter.CoinsMin);
        defaultRsp =
          defaultRsp && (!filter.CoinsMax || row.Coins <= filter.CoinsMax);
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
    this.#Filter  = filter
    const observable = subject.pipe(
      map((response: TransactionCoinsAndNodeDetails) => {
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

  getMyTransactions(): BehaviorSubject<TransactionCoinsAndNodeDetails> {
    return this.#myTransactionsBehaviorSubject;
  }
  getMyAllFilterTransactions(): BehaviorSubject<TransactionCoinsAndNodeDetails> {
    return this.#allTransactionsfilterBehaviorSubject;
  }
  fetchAllTransactions(mode: boolean) {
    this.transactionClientService.getAllTransactions(mode).subscribe(
      (r) => {
        const transactions = r.data?.getTransactionsCoins;
        const client = r.data?.self?.client;
        if (transactions && client) {
          if (transactions.Transactions) {
            this.#allTransactionsBehaviorSubject.next({
              transactions: transactions.Transactions,
              nodeDetails: client,
            });
           this.filterAllCoins(this.#Filter)
          }
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
  filterAllCoins(filter: filterTranctionCoinSubject) {
    const observable = this.filter(
      filter,
      this.#allTransactionsBehaviorSubject
    );
    observable.subscribe((r: TransactionCoinsAndNodeDetails) => {
      this.#allTransactionsfilterBehaviorSubject.next(r);
    });
  }
}
