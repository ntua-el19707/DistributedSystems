import { TestBed } from '@angular/core/testing';

import { TransactionBehaviorService } from './transaction-behavior.service';

describe('TransactionBehaviorService', () => {
  let service: TransactionBehaviorService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TransactionBehaviorService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
