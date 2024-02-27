import { TestBed } from '@angular/core/testing';

import { ClientsBehaviorSubjectService } from './clients-behavior-subject.service';

describe('ClientsBehaviorSubjectService', () => {
  let service: ClientsBehaviorSubjectService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ClientsBehaviorSubjectService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
