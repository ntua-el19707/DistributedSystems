import { TestBed } from '@angular/core/testing';

import { GraphQLBehaviorSubjectService } from './graph-qlbehavior-subject.service';

describe('GraphQLBehaviorSubjectService', () => {
  let service: GraphQLBehaviorSubjectService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(GraphQLBehaviorSubjectService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
