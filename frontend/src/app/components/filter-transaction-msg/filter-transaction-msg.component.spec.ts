import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FilterTransactionMsgComponent } from './filter-transaction-msg.component';

describe('FilterTransactionMsgComponent', () => {
  let component: FilterTransactionMsgComponent;
  let fixture: ComponentFixture<FilterTransactionMsgComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [FilterTransactionMsgComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(FilterTransactionMsgComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
