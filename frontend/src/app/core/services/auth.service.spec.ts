/**
 * @vitest-environment jsdom
 */
import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { Router } from '@angular/router';
import { PLATFORM_ID } from '@angular/core';
import { AuthService, User, AuthResponse } from './auth.service';
import { vi, expect, describe, it, beforeEach, afterEach } from 'vitest';

describe('AuthService', () => {
    let service: AuthService;
    let httpMock: HttpTestingController;
    let routerMock: any;

    const mockUser: User = {
        id: 1, email: 'test@example.com', name: 'Test User', created_at: '', updated_at: ''
    };

    function initTestBed() {
        TestBed.resetTestingModule();
        routerMock = { navigate: vi.fn() };
        TestBed.configureTestingModule({
            imports: [HttpClientTestingModule],
            providers: [
                AuthService,
                { provide: Router, useValue: routerMock },
                { provide: PLATFORM_ID, useValue: 'browser' }
            ]
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

        service.login({ email: 'test@example.com', password: 'password' }).subscribe(res => {
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
