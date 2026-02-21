import { HttpErrorResponse, HttpEvent, HttpHandlerFn, HttpRequest, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { catchError, filter, switchMap, take, throwError, Observable } from 'rxjs';
import { AuthService } from '../services/auth.service';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
    const authService = inject(AuthService);

    return next(req).pipe(
        catchError((error) => {
            if (error instanceof HttpErrorResponse && error.status === 401) {
                // ログインやリフレッシュ自体の失敗はそのままエラーとして流す
                if (req.url.includes('/auth/login') || req.url.includes('/auth/refresh')) {
                    authService.clearAuth();
                    return throwError(() => error);
                }

                return handleUnauthorizedError(authService, req, next);
            }
            return throwError(() => error);
        })
    );
};

const handleUnauthorizedError = (authService: AuthService, req: HttpRequest<any>, next: HttpHandlerFn): Observable<HttpEvent<any>> => {
    if (!authService.isRefreshing) {
        authService.isRefreshing = true;
        authService.refreshTokenSubject.next(null);

        return authService.refresh().pipe(
            switchMap(() => {
                authService.isRefreshing = false;
                authService.refreshTokenSubject.next('success'); // 成功を通知
                return next(req);
            }),
            catchError((err) => {
                authService.isRefreshing = false;
                authService.refreshTokenSubject.next('error'); // 待機中の他リクエストに失敗を通知
                authService.clearAuth();
                return throwError(() => err);
            })
        );
    } else {
        // すでにリフレッシュ中の場合は、信号を待ってから再試行
        return authService.refreshTokenSubject.pipe(
            filter(token => token !== null),
            take(1),
            switchMap((token) => {
                if (token === 'error') {
                    return throwError(() => new HttpErrorResponse({ status: 401, statusText: 'Unauthorized' }));
                }
                return next(req);
            })
        );
    }
};
