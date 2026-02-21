/**
 * @vitest-environment jsdom
 */
import { HttpRequest, HttpHandlerFn, HttpErrorResponse, HttpResponse } from '@angular/common/http';
import { authInterceptorInternal } from './auth.interceptor';
import { AuthService } from '../services/auth.service';
import { of, throwError, BehaviorSubject, firstValueFrom } from 'rxjs';
import { vi, expect, describe, it, beforeEach } from 'vitest';

describe('authInterceptorInternal', () => {
    let authServiceMock: any;

    beforeEach(() => {
        authServiceMock = {
            refresh: vi.fn(),
            isRefreshing: false,
            refreshTokenSubject: new BehaviorSubject<any>(null),
            clearAuth: vi.fn()
        };
    });

    it('should call next(req) for normal responses', async () => {
        const req = new HttpRequest('GET', '/api/test');
        const expectedResponse = new HttpResponse({ status: 200 });
        const next: HttpHandlerFn = vi.fn().mockReturnValue(of(expectedResponse));

        const result$ = authInterceptorInternal(req, next, authServiceMock);
        const result = await firstValueFrom(result$);

        expect(result).toBe(expectedResponse);
        expect(next).toHaveBeenCalledWith(req);
    });

    it('should trigger refresh on 401 error', async () => {
        const req = new HttpRequest('GET', '/api/test');
        const errorResponse = new HttpErrorResponse({ status: 401 });
        const next: HttpHandlerFn = vi.fn()
            .mockReturnValueOnce(throwError(() => errorResponse))
            .mockReturnValueOnce(of(new HttpResponse({ status: 200 })));

        authServiceMock.refresh.mockReturnValue(of({}));

        const result$ = authInterceptorInternal(req, next, authServiceMock);
        const result = await firstValueFrom(result$) as HttpResponse<any>;

        expect(result.status).toBe(200);
        expect(authServiceMock.refresh).toHaveBeenCalled();
        expect(next).toHaveBeenCalledTimes(2);
    });

    it('should clear auth and fail if refresh fails', async () => {
        const req = new HttpRequest('GET', '/api/test');
        const errorResponse = new HttpErrorResponse({ status: 401 });
        const next: HttpHandlerFn = vi.fn().mockReturnValue(throwError(() => errorResponse));

        authServiceMock.refresh.mockReturnValue(throwError(() => new Error('refresh failed')));

        const result$ = authInterceptorInternal(req, next, authServiceMock);

        await expect(firstValueFrom(result$)).rejects.toThrow();
        expect(authServiceMock.clearAuth).toHaveBeenCalled();
    });
});
