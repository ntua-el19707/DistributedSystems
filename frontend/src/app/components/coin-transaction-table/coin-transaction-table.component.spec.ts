import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CoinTransactionTableComponent } from './coin-transaction-table.component';

describe('CoinTransactionTableComponent', () => {
  let component: CoinTransactionTableComponent;
  let fixture: ComponentFixture<CoinTransactionTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CoinTransactionTableComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(CoinTransactionTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
