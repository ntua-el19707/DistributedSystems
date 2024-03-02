import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DateAndTimePickerComponent } from './date-and-time-picker.component';

describe('DateAndTimePickerComponent', () => {
  let component: DateAndTimePickerComponent;
  let fixture: ComponentFixture<DateAndTimePickerComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DateAndTimePickerComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(DateAndTimePickerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
