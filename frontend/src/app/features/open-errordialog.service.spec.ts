import { TestBed } from '@angular/core/testing';

import { OpenErrordialogService } from './open-errordialog.service';

describe('OpenErrordialogService', () => {
  let service: OpenErrordialogService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(OpenErrordialogService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
