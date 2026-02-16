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
  private readonly API_URL = 'http://localhost:8080/api';
  private readonly TOKEN_KEY = 'auth_token';
  private readonly USER_KEY = 'auth_user';

  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);
  private readonly platformId = inject(PLATFORM_ID);

  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$ = this.currentUserSubject.asObservable();

  constructor() {
    this.restoreSession();
  }

  private restoreSession() {
    if (!isPlatformBrowser(this.platformId)) {
      return;
    }
    try {
      const storedUser = localStorage.getItem(this.USER_KEY);
      const token = localStorage.getItem(this.TOKEN_KEY);
      if (storedUser && storedUser !== 'undefined' && token && token !== 'undefined') {
        this.currentUserSubject.next(JSON.parse(storedUser));
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
    // Optimistic logout
    this.clearAuth();
    this.router.navigate(['/login']);
    // Optional: Call backend to invalidate token if needed
    // this.http.post(`${this.API_URL}/auth/logout`, {}).subscribe();
  }

  private handleAuthResponse(response: AuthResponse) {
    this.saveToken(response.token);
    this.saveUser(response.user);
    this.currentUserSubject.next(response.user);
  }

  private saveToken(token: string): void {
    if (isPlatformBrowser(this.platformId) && token) {
      localStorage.setItem(this.TOKEN_KEY, token);
    }
  }

  private saveUser(user: User): void {
    if (isPlatformBrowser(this.platformId) && user) {
      localStorage.setItem(this.USER_KEY, JSON.stringify(user));
    }
  }

  getToken(): string | null {
    if (isPlatformBrowser(this.platformId)) {
      return localStorage.getItem(this.TOKEN_KEY);
    }
    return null;
  }

  clearAuth(): void {
    if (isPlatformBrowser(this.platformId)) {
      localStorage.removeItem(this.TOKEN_KEY);
      localStorage.removeItem(this.USER_KEY);
    }
    this.currentUserSubject.next(null);
  }

  isLoggedIn(): boolean {
    return !!this.getToken();
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
