/**
 * @vitest-environment jsdom
 */
import { guestGuard } from './guest.guard';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { Router, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { Injector, runInInjectionContext } from '@angular/core';

describe('guestGuard', () => {
  let mockAuthService: { isLoggedIn: ReturnType<typeof vi.fn> };
  let mockRouter: { navigate: ReturnType<typeof vi.fn> };
  let injector: Injector;

  beforeEach(() => {
    mockAuthService = { isLoggedIn: vi.fn() };
    mockRouter = { navigate: vi.fn() };

    injector = Injector.create({
      providers: [
        { provide: AuthService, useValue: mockAuthService },
        { provide: Router, useValue: mockRouter },
      ],
    });
  });

  it('should return true if user is not logged in', () => {
    mockAuthService.isLoggedIn.mockReturnValue(false);

    // Injector のコンテキスト内で guard を実行
    const result = runInInjectionContext(injector, () =>
      guestGuard({} as ActivatedRouteSnapshot, {} as RouterStateSnapshot),
    );

    expect(result).toBe(true);
  });

  it('should redirect if user is logged in', () => {
    mockAuthService.isLoggedIn.mockReturnValue(true);

    const result = runInInjectionContext(injector, () =>
      guestGuard({} as ActivatedRouteSnapshot, {} as RouterStateSnapshot),
    );

    expect(result).toBe(false);
    expect(mockRouter.navigate).toHaveBeenCalledWith(['/recipes']);
  });
});
