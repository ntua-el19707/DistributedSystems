import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FilterTransactionCoinsComponent } from './filter-transaction-coins.component';

describe('FilterTransactionCoinsComponent', () => {
  let component: FilterTransactionCoinsComponent;
  let fixture: ComponentFixture<FilterTransactionCoinsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [FilterTransactionCoinsComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(FilterTransactionCoinsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
