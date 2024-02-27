import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TrasferCoinsComponent } from './trasfer-coins.component';

describe('TrasferCoinsComponent', () => {
  let component: TrasferCoinsComponent;
  let fixture: ComponentFixture<TrasferCoinsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TrasferCoinsComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(TrasferCoinsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
