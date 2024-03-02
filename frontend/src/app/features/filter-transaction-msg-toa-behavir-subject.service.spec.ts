import { TestBed } from '@angular/core/testing';

import { FilterTransactionMsgToaBehavirSubjectService } from './filter-transaction-msg-toa-behavir-subject.service';

describe('FilterTransactionMsgToaBehavirSubjectService', () => {
  let service: FilterTransactionMsgToaBehavirSubjectService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(FilterTransactionMsgToaBehavirSubjectService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
