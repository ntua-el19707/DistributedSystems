import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TranferCoinsFormComponent } from './tranfer-coins-form.component';

describe('TranferCoinsFormComponent', () => {
  let component: TranferCoinsFormComponent;
  let fixture: ComponentFixture<TranferCoinsFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TranferCoinsFormComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(TranferCoinsFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
