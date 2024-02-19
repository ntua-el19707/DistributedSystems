import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MessengerPageComponent } from './messenger-page.component';

describe('MessengerPageComponent', () => {
  let component: MessengerPageComponent;
  let fixture: ComponentFixture<MessengerPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MessengerPageComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(MessengerPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
