import { TestBed } from '@angular/core/testing';

import { TransactionMsgBehaviorService } from './transaction-msg-behavior.service';

describe('TransactionMsgBehaviorService', () => {
  let service: TransactionMsgBehaviorService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TransactionMsgBehaviorService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
