import { inject, Injectable, PLATFORM_ID } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { tap } from 'rxjs/operators';
import { Router } from '@angular/router';
import { isPlatformBrowser } from '@angular/common';

export interface User {
  id: number;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface SignupRequest {
  email: string;
  password: string;
  name: string;
}

export interface LoginRequest {
  email: string;
  password: string;
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
  private readonly request = inject('REQUEST' as any, { optional: true });

  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$ = this.currentUserSubject.asObservable();

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
        const req = this.request as any;
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

  signup(data: SignupRequest): Observable<AuthResponse> {
    return this.http
      .post<AuthResponse>(`${this.API_URL}/auth/signup`, data)
      .pipe(tap((response) => this.handleAuthResponse(response)));
  }

  login(data: LoginRequest): Observable<AuthResponse> {
    return this.http
      .post<AuthResponse>(`${this.API_URL}/auth/login`, data)
      .pipe(tap((response) => this.handleAuthResponse(response)));
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
    this.saveUser(response.user);
    this.currentUserSubject.next(response.user);
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
      for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) === ' ') c = c.substring(1, c.length);
        if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
      }
    }
    return null;
  }

  private parseCookieHeader(header: string, name: string): string | null {
    const nameEQ = name + '=';
    const ca = header.split(';');
    for (let i = 0; i < ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) === ' ') c = c.substring(1, c.length);
      if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
    }
    return null;
  }

  clearAuth(): void {
    if (isPlatformBrowser(this.platformId)) {
      // USER_KEY Cookie を削除
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

  changePassword(data: { current_password: string; new_password: string }): Observable<void> {
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

  requestPasswordReset(email: string): Observable<MessageResponse> {
    return this.http.post<MessageResponse>(`${this.API_URL}/auth/password-reset/request`, { email });
  }

  resetPassword(token: string, newPassword: string): Observable<MessageResponse> {
    return this.http.post<MessageResponse>(`${this.API_URL}/auth/password-reset`, {
      token,
      new_password: newPassword,
    });
  }
}

export interface MessageResponse {
  message: string;
}
