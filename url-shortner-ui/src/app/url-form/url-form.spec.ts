import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UrlForm } from './url-form';

describe('UrlForm', () => {
  let component: UrlForm;
  let fixture: ComponentFixture<UrlForm>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [UrlForm]
    })
    .compileComponents();

    fixture = TestBed.createComponent(UrlForm);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
