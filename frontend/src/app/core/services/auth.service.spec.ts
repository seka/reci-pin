/**
 * @vitest-environment jsdom
 */
import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting, HttpTestingController } from '@angular/common/http/testing';
import { Router } from '@angular/router';
import { PLATFORM_ID } from '@angular/core';
import { AuthService } from './auth.service';
import { User } from '../models/user.model';
import { AuthResponse } from './responses/auth.response';
import { vi, expect, describe, it, beforeEach, afterEach } from 'vitest';

describe('AuthService', () => {
  let service: AuthService;
  let httpMock: HttpTestingController;
  let routerMock: { navigate: ReturnType<typeof vi.fn> };

  const mockUser: User = {
    id: 1,
    email: 'test@example.com',
    name: 'Test User',
    createdAt: '',
    updatedAt: '',
  };

  function initTestBed() {
    routerMock = { navigate: vi.fn() };
    TestBed.configureTestingModule({
      providers: [
        AuthService,
        provideHttpClient(),
        provideHttpClientTesting(),
        { provide: Router, useValue: routerMock },
        { provide: PLATFORM_ID, useValue: 'browser' },
      ],
    });
    service = TestBed.inject(AuthService);
    httpMock = TestBed.inject(HttpTestingController);
  }

  beforeEach(() => {
    // Clear cookies before each test
    document.cookie = 'auth_user=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
  });

  afterEach(() => {
    if (httpMock) {
      httpMock.verify();
    }
  });

  it('should restore user from cookie on initialization', () => {
    const userString = encodeURIComponent(JSON.stringify(mockUser));
    document.cookie = `auth_user=${userString}; path=/;`;

    initTestBed();
    expect(service.currentUserValue).toEqual(mockUser);
  });

  it('should login and set current user', () => {
    initTestBed();
    const mockResponse: AuthResponse = { token: '', user: mockUser };

    service.login({ email: 'test@example.com', password: 'password' }).subscribe((res) => {
      expect(res.user).toEqual(mockUser);
    });

    const req = httpMock.expectOne('/api/auth/login');
    req.flush(mockResponse);
    expect(service.currentUserValue).toEqual(mockUser);
  });

  it('should logout and clear state', () => {
    initTestBed();
    service.logout();

    const req = httpMock.expectOne('/api/auth/logout');
    req.flush({});

    expect(service.currentUserValue).toBeNull();
    expect(routerMock.navigate).toHaveBeenCalledWith(['/login']);
  });
});
