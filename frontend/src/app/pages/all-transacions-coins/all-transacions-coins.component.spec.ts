import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AllTransacionsCoinsComponent } from './all-transacions-coins.component';

describe('AllTransacionsCoinsComponent', () => {
  let component: AllTransacionsCoinsComponent;
  let fixture: ComponentFixture<AllTransacionsCoinsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AllTransacionsCoinsComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AllTransacionsCoinsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
