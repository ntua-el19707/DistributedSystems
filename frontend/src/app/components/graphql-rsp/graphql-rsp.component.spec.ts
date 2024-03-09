import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GraphqlRspComponent } from './graphql-rsp.component';

describe('GraphqlRspComponent', () => {
  let component: GraphqlRspComponent;
  let fixture: ComponentFixture<GraphqlRspComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [GraphqlRspComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(GraphqlRspComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
