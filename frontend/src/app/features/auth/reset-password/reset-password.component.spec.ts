/**
 * @vitest-environment jsdom
 */
import { TestBed } from '@angular/core/testing';
import { ActivatedRoute, provideRouter, Router } from '@angular/router';
import { TranslocoTestingModule } from '@jsverse/transloco';
import { of, Subject } from 'rxjs';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { AuthService } from '../../../core/services/auth.service';
import { ResetPasswordComponent } from './reset-password.component';

describe('ResetPasswordComponent', () => {
  let response$: Subject<{ message: string }>;
  const resetPassword = vi.fn();

  beforeEach(async () => {
    response$ = new Subject<{ message: string }>();
    resetPassword.mockReset();
    resetPassword.mockReturnValue(response$.asObservable());

    await TestBed.configureTestingModule({
      imports: [
        ResetPasswordComponent,
        TranslocoTestingModule.forRoot({
          langs: {
            ja: {
              FEATURES: {
                AUTH: {
                  RESET_PASSWORD: {
                    INVALID_LINK: 'リンクが無効です',
                    FAILED_EXPIRED: '再設定に失敗しました',
                    SETTING_PASSWORD: '設定中',
                    SUBMIT: '設定する',
                  },
                },
              },
            },
          },
          translocoConfig: { availableLangs: ['ja'], defaultLang: 'ja' },
        }),
      ],
      providers: [
        provideRouter([]),
        { provide: ActivatedRoute, useValue: { queryParams: of({ token: 'valid-token' }) } },
        { provide: AuthService, useValue: { resetPassword } },
      ],
    }).compileComponents();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('submitするとリセットAPIを呼び、成功メッセージを表示する', async () => {
    vi.spyOn(TestBed.inject(Router), 'navigate').mockResolvedValue(true);
    const fixture = TestBed.createComponent(ResetPasswordComponent);
    fixture.autoDetectChanges();

    const input = fixture.nativeElement.querySelector('input[type="password"]') as HTMLInputElement;
    input.value = 'new-password1';
    input.dispatchEvent(new Event('input', { bubbles: true }));
    await fixture.whenStable();

    const form = fixture.nativeElement.querySelector('form') as HTMLFormElement;
    form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));

    expect(resetPassword).toHaveBeenCalledWith({
      token: 'valid-token',
      newPassword: 'new-password1',
    });
    expect(fixture.nativeElement.textContent).toContain('設定中');

    response$.next({ message: 'パスワードを設定しました' });
    await fixture.whenStable();
    expect(fixture.nativeElement.textContent).toContain('パスワードを設定しました');
  });

  it('APIが失敗するとエラーメッセージを表示する', async () => {
    vi.spyOn(console, 'error').mockImplementation(() => undefined);
    const fixture = TestBed.createComponent(ResetPasswordComponent);
    fixture.autoDetectChanges();

    const input = fixture.nativeElement.querySelector('input[type="password"]') as HTMLInputElement;
    input.value = 'new-password1';
    input.dispatchEvent(new Event('input', { bubbles: true }));
    await fixture.whenStable();

    const form = fixture.nativeElement.querySelector('form') as HTMLFormElement;
    form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));
    response$.error(new Error('request failed'));
    await fixture.whenStable();

    expect(fixture.nativeElement.textContent).toContain('再設定に失敗しました');
  });

  it('tokenがない場合は無効リンクのメッセージを表示する', async () => {
    TestBed.overrideProvider(ActivatedRoute, { useValue: { queryParams: of({}) } });
    const fixture = TestBed.createComponent(ResetPasswordComponent);
    fixture.autoDetectChanges();
    await fixture.whenStable();

    expect(fixture.nativeElement.querySelector('.alert.error')?.textContent).toBeTruthy();
  });
});
