import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AllMsgComponent } from './all-msg.component';

describe('AllMsgComponent', () => {
  let component: AllMsgComponent;
  let fixture: ComponentFixture<AllMsgComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AllMsgComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AllMsgComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
