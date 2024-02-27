import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TransactionMsgTableComponent } from './transaction-msg-table.component';

describe('TransactionMsgTableComponent', () => {
  let component: TransactionMsgTableComponent;
  let fixture: ComponentFixture<TransactionMsgTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TransactionMsgTableComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(TransactionMsgTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
