import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject, map } from 'rxjs';
import { tap } from 'rxjs/operators';
import { Router } from '@angular/router';
import { User } from '../models/user.model';
import {
  toSignupRequest,
  toLoginRequest,
  toChangePasswordRequest,
  toPasswordResetRequest,
  toPasswordResetConfirmRequest,
} from './requests/auth.request';
import { AuthResponse, MessageResponse, toUserModel } from './responses/auth.response';
import {
  LoginFormModel,
  SignupFormModel,
  ChangePasswordFormModel,
  PasswordResetFormModel,
  PasswordResetConfirmFormModel,
} from '../models/auth.model';

export type RefreshState = 'success' | 'error' | null;

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly API_URL = '/api';
  private readonly USER_KEY = 'auth_user';

  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);

  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$ = this.currentUserSubject.asObservable();

  public isRefreshing = false;
  public refreshTokenSubject = new BehaviorSubject<RefreshState>(null);

  constructor() {
    this.restoreSession();
  }

  private restoreSession() {
    try {
      const storedUser = this.getCookie(this.USER_KEY);

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

  signup(data: SignupFormModel): Observable<User> {
    const request = toSignupRequest(data);
    return this.http.post<AuthResponse>(`${this.API_URL}/auth/signup`, request).pipe(
      tap((response) => this.handleAuthResponse(response)),
      map(toUserModel),
    );
  }

  login(data: LoginFormModel): Observable<User> {
    const request = toLoginRequest(data);
    return this.http.post<AuthResponse>(`${this.API_URL}/auth/login`, request).pipe(
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
    if (user) {
      this.setCookie(this.USER_KEY, JSON.stringify(user), 7);
    }
  }

  private setCookie(name: string, value: string, days: number): void {
    const date = new Date();
    date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
    const expires = '; expires=' + date.toUTCString();
    document.cookie = name + '=' + encodeURIComponent(value) + expires + '; path=/; SameSite=Lax';
  }

  private getCookie(name: string): string | null {
    const nameEQ = name + '=';
    const ca = document.cookie.split(';');
    for (const element of ca) {
      let c = element;
      while (c.charAt(0) === ' ') c = c.substring(1, c.length);
      if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
    }
    return null;
  }

  clearAuth(): void {
    // HttpOnly 属性ではないものを削除する
    // (auth_token などの HttpOnly Cookie は JS から削除できないため、サーバー側でクリアする必要がある)
    this.setCookie(this.USER_KEY, '', -1);
    this.currentUserSubject.next(null);
  }

  isLoggedIn(): boolean {
    return !!this.currentUserValue;
  }

  get currentUserValue(): User | null {
    return this.currentUserSubject.value;
  }

  changePassword(data: ChangePasswordFormModel): Observable<void> {
    const request = toChangePasswordRequest(data);
    return this.http.put<void>(`${this.API_URL}/auth/password`, request);
  }

  withdraw(): Observable<void> {
    return this.http.delete<void>(`${this.API_URL}/auth/withdraw`).pipe(
      tap(() => {
        this.clearAuth();
        this.router.navigate(['/login']);
      }),
    );
  }

  requestPasswordReset(data: PasswordResetFormModel): Observable<MessageResponse> {
    const request = toPasswordResetRequest(data);
    return this.http.post<MessageResponse>(`${this.API_URL}/auth/password-reset/request`, request);
  }

  resetPassword(data: PasswordResetConfirmFormModel): Observable<MessageResponse> {
    const request = toPasswordResetConfirmRequest(data);
    return this.http.post<MessageResponse>(`${this.API_URL}/auth/password-reset`, request);
  }

  refresh(): Observable<User> {
    return this.http.post<AuthResponse>(`${this.API_URL}/auth/refresh`, {}).pipe(
      tap((response) => this.handleAuthResponse(response)),
      map(toUserModel),
    );
  }
}
