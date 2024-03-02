import { TestBed } from '@angular/core/testing';

import { FilterTransactionCoinToaBehavirSubjectService } from './filter-transaction-coin-toa-behavir-subject.service';

describe('FilterTransactionCoinToaBehavirSubjectService', () => {
  let service: FilterTransactionCoinToaBehavirSubjectService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(FilterTransactionCoinToaBehavirSubjectService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
