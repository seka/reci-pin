/**
 * @vitest-environment jsdom
 */
import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting, HttpTestingController } from '@angular/common/http/testing';
import { Router } from '@angular/router';
import { AuthService } from './auth.service';
import { User, UserData } from '../models/user.model';
import { AuthResponse } from './responses/auth.response';
import { LoginFormModel } from '../models/auth.model';
import { vi, expect, describe, it, beforeEach, afterEach } from 'vitest';

describe('AuthService', () => {
  let service: AuthService;
  let httpMock: HttpTestingController;
  let routerMock: { navigate: ReturnType<typeof vi.fn> };

  const userData: UserData = {
    id: 1,
    email: 'test@example.com',
    name: 'Test User',
    createdAt: '',
    updatedAt: '',
  };
  const mockUser = new User(userData);

  function initTestBed() {
    routerMock = { navigate: vi.fn() };
    TestBed.configureTestingModule({
      providers: [
        AuthService,
        provideHttpClient(),
        provideHttpClientTesting(),
        { provide: Router, useValue: routerMock },
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
    const mockResponse: AuthResponse = { token: '', user: userData };

    service
      .login({ email: 'test@example.com', password: 'password' } as LoginFormModel)
      .subscribe((user) => {
        expect(user).toEqual(mockUser);
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
