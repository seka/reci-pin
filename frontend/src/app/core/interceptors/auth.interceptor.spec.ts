/**
 * @vitest-environment jsdom
 */
import { HttpRequest, HttpHandlerFn, HttpErrorResponse, HttpResponse } from '@angular/common/http';
import { authInterceptorInternal } from './auth.interceptor';
import { of, throwError, BehaviorSubject, firstValueFrom } from 'rxjs';
import { vi, expect, describe, it, beforeEach } from 'vitest';

describe('authInterceptorInternal', () => {
  let authServiceMock: {
    refresh: ReturnType<typeof vi.fn>;
    isRefreshing: boolean;
    refreshTokenSubject: BehaviorSubject<string | null>;
    clearAuth: ReturnType<typeof vi.fn>;
  };

  beforeEach(() => {
    authServiceMock = {
      refresh: vi.fn(),
      isRefreshing: false,
      refreshTokenSubject: new BehaviorSubject<string | null>(null),
      clearAuth: vi.fn(),
    };
  });

  it('should call next(req) for normal responses', async () => {
    const req = new HttpRequest('GET', '/api/test');
    const expectedResponse = new HttpResponse({ status: 200 });
    const next: HttpHandlerFn = vi.fn().mockReturnValue(of(expectedResponse));

    const result$ = authInterceptorInternal(req, next, authServiceMock as unknown as never);
    const result = await firstValueFrom(result$);

    expect(result).toBe(expectedResponse);
    expect(next).toHaveBeenCalledWith(req);
  });

  it('should trigger refresh on 401 error', async () => {
    const req = new HttpRequest('GET', '/api/test');
    const errorResponse = new HttpErrorResponse({ status: 401 });
    const next: HttpHandlerFn = vi
      .fn()
      .mockReturnValueOnce(throwError(() => errorResponse))
      .mockReturnValueOnce(of(new HttpResponse({ status: 200 })));

    authServiceMock.refresh.mockReturnValue(of({}));

    const result$ = authInterceptorInternal(req, next, authServiceMock as unknown as never);
    const result = (await firstValueFrom(result$)) as HttpResponse<unknown>;

    expect(result.status).toBe(200);
    expect(authServiceMock.refresh).toHaveBeenCalled();
    expect(next).toHaveBeenCalledTimes(2);
  });

  it('should clear auth and fail if refresh fails', async () => {
    const req = new HttpRequest('GET', '/api/test');
    const errorResponse = new HttpErrorResponse({ status: 401 });
    const next: HttpHandlerFn = vi.fn().mockReturnValue(throwError(() => errorResponse));

    authServiceMock.refresh.mockReturnValue(throwError(() => new Error('refresh failed')));

    const result$ = authInterceptorInternal(req, next, authServiceMock as unknown as never);

    await expect(firstValueFrom(result$)).rejects.toThrow();
    expect(authServiceMock.clearAuth).toHaveBeenCalled();
  });
});
