import { inject, Injectable, PLATFORM_ID } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject, map } from 'rxjs';
import { tap } from 'rxjs/operators';
import { Router } from '@angular/router';
import { isPlatformBrowser } from '@angular/common';
import { User } from '../models/user.model';
import {
  SignupRequest,
  LoginRequest,
  ChangePasswordRequest,
  PasswordResetRequest,
  PasswordResetConfirmRequest,
} from './requests/auth.request';
import { AuthResponse, MessageResponse, toUserModel } from './responses/auth.response';

export type RefreshState = 'success' | 'error' | null;

export interface SsrRequest {
  headers?: {
    get?: (name: string) => string;
    cookie?: string;
  };
  cookie?: string;
}

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly API_URL = '/api';
  private readonly USER_KEY = 'auth_user';

  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);
  private readonly platformId = inject(PLATFORM_ID);
  // SSR Context
  private readonly request = inject<SsrRequest>(
    'REQUEST' as unknown as import('@angular/core').InjectionToken<SsrRequest>,
    {
      optional: true,
    },
  );

  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$ = this.currentUserSubject.asObservable();

  public isRefreshing = false;
  public refreshTokenSubject = new BehaviorSubject<RefreshState>(null);

  constructor() {
    this.restoreSession();
  }

  private restoreSession() {
    try {
      let storedUser: string | null = null;

      if (isPlatformBrowser(this.platformId)) {
        storedUser = this.getCookie(this.USER_KEY);
      } else if (this.request) {
        // SSR context: read from request headers
        const req = this.request;
        const cookieHeader = req.headers?.get?.('cookie') || req.headers?.cookie;
        if (cookieHeader) {
          storedUser = this.parseCookieHeader(cookieHeader, this.USER_KEY);
        }
      }

      if (storedUser && storedUser !== 'undefined') {
        this.currentUserSubject.next(JSON.parse(decodeURIComponent(storedUser)));
      } else {
        this.clearAuth();
      }
    } catch (e) {
      console.error('Error restoring session', e);
      this.clearAuth();
    }
  }

  signup(data: SignupRequest): Observable<User> {
    return this.http
      .post<AuthResponse>(`${this.API_URL}/auth/signup`, data)
      .pipe(
        tap((response) => this.handleAuthResponse(response)),
        map(toUserModel),
      );
  }

  login(data: LoginRequest): Observable<User> {
    return this.http
      .post<AuthResponse>(`${this.API_URL}/auth/login`, data)
      .pipe(
        tap((response) => this.handleAuthResponse(response)),
        map(toUserModel),
      );
  }

  logout(): void {
    // サーバーサイドの Cookie (HttpOnly) も削除するために必要
    this.http.post(`${this.API_URL}/auth/logout`, {}).subscribe({
      next: () => this.handleLogoutSuccess(),
      error: () => this.handleLogoutSuccess(), // エラーでもフロント側はクリアする
    });
  }

  private handleLogoutSuccess(): void {
    this.clearAuth();
    this.router.navigate(['/login']);
  }

  private handleAuthResponse(response: AuthResponse) {
    // auth_token は HttpOnly Cookie としてサーバーから設定されるため、
    // ここで保存する必要はなくなりました。
    const user = toUserModel(response);
    this.saveUser(user);
    this.currentUserSubject.next(user);
  }

  private saveUser(user: User): void {
    if (isPlatformBrowser(this.platformId) && user) {
      this.setCookie(this.USER_KEY, JSON.stringify(user), 7);
    }
  }

  private setCookie(name: string, value: string, days: number): void {
    if (isPlatformBrowser(this.platformId)) {
      const date = new Date();
      date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
      const expires = '; expires=' + date.toUTCString();
      document.cookie = name + '=' + encodeURIComponent(value) + expires + '; path=/; SameSite=Lax';
    }
  }

  private getCookie(name: string): string | null {
    if (isPlatformBrowser(this.platformId)) {
      const nameEQ = name + '=';
      const ca = document.cookie.split(';');
      for (const element of ca) {
        let c = element;
        while (c.charAt(0) === ' ') c = c.substring(1, c.length);
        if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
      }
    }
    return null;
  }

  private parseCookieHeader(header: string, name: string): string | null {
    const nameEQ = name + '=';
    const ca = header.split(';');
    for (const element of ca) {
      let c = element;
      while (c.charAt(0) === ' ') c = c.substring(1, c.length);
      if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
    }
    return null;
  }

  clearAuth(): void {
    if (isPlatformBrowser(this.platformId)) {
      // HttpOnly 属性ではないものを削除する
      // (auth_token などの HttpOnly Cookie は JS から削除できないため、サーバー側でクリアする必要がある)
      this.setCookie(this.USER_KEY, '', -1);
    }
    this.currentUserSubject.next(null);
  }

  isLoggedIn(): boolean {
    return !!this.currentUserValue;
  }

  get currentUserValue(): User | null {
    return this.currentUserSubject.value;
  }

  changePassword(data: ChangePasswordRequest): Observable<void> {
    return this.http.put<void>(`${this.API_URL}/auth/password`, data);
  }

  withdraw(): Observable<void> {
    return this.http.delete<void>(`${this.API_URL}/auth/withdraw`).pipe(
      tap(() => {
        this.clearAuth();
        this.router.navigate(['/login']);
      }),
    );
  }

  requestPasswordReset(data: PasswordResetRequest): Observable<MessageResponse> {
    return this.http.post<MessageResponse>(`${this.API_URL}/auth/password-reset/request`, data);
  }

  resetPassword(data: PasswordResetConfirmRequest): Observable<MessageResponse> {
    return this.http.post<MessageResponse>(`${this.API_URL}/auth/password-reset`, data);
  }

  refresh(): Observable<User> {
    return this.http
      .post<AuthResponse>(`${this.API_URL}/auth/refresh`, {})
      .pipe(
        tap((response) => this.handleAuthResponse(response)),
        map(toUserModel),
      );
  }
}
