import { TestBed } from '@angular/core/testing';

import { ClientsClientService } from './clients-client.service';

describe('ClientsClientService', () => {
  let service: ClientsClientService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ClientsClientService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
