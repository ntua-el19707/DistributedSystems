import { TestBed } from '@angular/core/testing';

import { NodeInfoBehaviorService } from './node-info-behavior.service';

describe('NodeInfoBehaviorService', () => {
  let service: NodeInfoBehaviorService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(NodeInfoBehaviorService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
