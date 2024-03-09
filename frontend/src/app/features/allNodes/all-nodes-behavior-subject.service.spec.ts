import { TestBed } from '@angular/core/testing';

import { AllNodesBehaviorSubjectService } from './all-nodes-behavior-subject.service';

describe('AllNodesBehaviorSubjectService', () => {
  let service: AllNodesBehaviorSubjectService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(AllNodesBehaviorSubjectService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
