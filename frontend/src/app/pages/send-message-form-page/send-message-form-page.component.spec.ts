import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SendMessageFormPageComponent } from './send-message-form-page.component';

describe('SendMessageFormPageComponent', () => {
  let component: SendMessageFormPageComponent;
  let fixture: ComponentFixture<SendMessageFormPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SendMessageFormPageComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(SendMessageFormPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
