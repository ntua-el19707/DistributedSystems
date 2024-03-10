import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormStakeComponent } from './form-stake.component';

describe('FormStakeComponent', () => {
  let component: FormStakeComponent;
  let fixture: ComponentFixture<FormStakeComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [FormStakeComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(FormStakeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
