/**
 * @vitest-environment jsdom
 */
import { TestBed } from '@angular/core/testing';
import { App } from './app';
import { TranslateModule } from '@ngx-translate/core';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AuthService } from './core/services/auth.service';
import { of, BehaviorSubject } from 'rxjs';
import { PLATFORM_ID, NO_ERRORS_SCHEMA } from '@angular/core';
import { provideRouter } from '@angular/router';

describe.skip('App', () => {
  beforeEach(async () => {
    TestBed.resetTestingModule();
    await TestBed.configureTestingModule({
      imports: [App, TranslateModule.forRoot()],
      providers: [
        provideRouter([]),
        { provide: PLATFORM_ID, useValue: 'browser' },
        {
          provide: AuthService,
          useValue: {
            currentUser$: of(null),
            refreshTokenSubject: new BehaviorSubject(null),
            isRefreshing: false,
          },
        },
      ],
      schemas: [NO_ERRORS_SCHEMA],
    }).compileComponents();
  });

  it('should create the app', () => {
    const fixture = TestBed.createComponent(App);
    const app = fixture.componentInstance;
    expect(app).toBeTruthy();
  });
});
