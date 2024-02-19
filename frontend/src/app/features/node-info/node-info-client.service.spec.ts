import { TestBed } from '@angular/core/testing';

import { NodeInfoClientService } from './node-info-client.service';

describe('NodeInfoClientService', () => {
  let service: NodeInfoClientService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(NodeInfoClientService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
