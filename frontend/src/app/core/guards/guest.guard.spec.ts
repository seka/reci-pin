import '@angular/compiler';
import { guestGuard } from './guest.guard';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { PLATFORM_ID, Injector, runInInjectionContext } from '@angular/core';

describe('guestGuard', () => {
  let mockAuthService: any;
  let mockRouter: any;
  let injector: Injector;

  beforeEach(() => {
    mockAuthService = { isLoggedIn: vi.fn() };
    mockRouter = { navigate: vi.fn() };

    injector = Injector.create({
      providers: [
        { provide: PLATFORM_ID, useValue: 'browser' },
        { provide: AuthService, useValue: mockAuthService },
        { provide: Router, useValue: mockRouter },
      ],
    });
  });

  it('should return true if user is not logged in', () => {
    mockAuthService.isLoggedIn.mockReturnValue(false);

    // Injector のコンテキスト内で guard を実行
    const result = runInInjectionContext(injector, () =>
      guestGuard(undefined as any, undefined as any),
    );

    expect(result).toBe(true);
    expect(mockRouter.navigate).not.toHaveBeenCalled();
  });

  it('should return false and navigate to /recipes if user is logged in', () => {
    mockAuthService.isLoggedIn.mockReturnValue(true);

    const result = runInInjectionContext(injector, () =>
      guestGuard(undefined as any, undefined as any),
    );

    expect(result).toBe(false);
    expect(mockRouter.navigate).toHaveBeenCalledWith(['/recipes']);
  });
});
