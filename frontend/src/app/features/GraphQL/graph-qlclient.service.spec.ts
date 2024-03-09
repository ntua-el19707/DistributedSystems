import { TestBed } from '@angular/core/testing';

import { GraphQLClientService } from './graph-qlclient.service';

describe('GraphQLClientService', () => {
  let service: GraphQLClientService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(GraphQLClientService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
