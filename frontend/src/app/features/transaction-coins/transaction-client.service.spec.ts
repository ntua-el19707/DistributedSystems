import { TestBed } from '@angular/core/testing';

import { TransactionClientService } from './transaction-client.service';

describe('TransactionClientService', () => {
  let service: TransactionClientService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TransactionClientService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
