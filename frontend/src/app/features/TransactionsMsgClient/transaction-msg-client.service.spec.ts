import { TestBed } from '@angular/core/testing';

import { TransactionMsgClientService } from './transaction-msg-client.service';

describe('TransactionMsgClientService', () => {
  let service: TransactionMsgClientService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TransactionMsgClientService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
