/**
 * @vitest-environment jsdom
 */
import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { TranslocoTestingModule } from '@jsverse/transloco';
import { Subject } from 'rxjs';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { AuthService } from '../../../core/services/auth.service';
import { RequestPasswordResetComponent } from './request-password-reset.component';

describe('RequestPasswordResetComponent', () => {
  let response$: Subject<{ message: string }>;
  const requestPasswordReset = vi.fn();

  beforeEach(async () => {
    response$ = new Subject<{ message: string }>();
    requestPasswordReset.mockReset();
    requestPasswordReset.mockReturnValue(response$.asObservable());

    await TestBed.configureTestingModule({
      imports: [
        RequestPasswordResetComponent,
        TranslocoTestingModule.forRoot({
          langs: {
            ja: {
              COMMON: {
                SEND: '送信',
                SENDING: '送信中',
              },
              FEATURES: {
                AUTH: {
                  REQUEST_PASSWORD_RESET: {
                    FAILED: '送信に失敗しました',
                  },
                },
              },
            },
          },
          translocoConfig: {
            availableLangs: ['ja'],
            defaultLang: 'ja',
          },
        }),
      ],
      providers: [
        provideRouter([]),
        {
          provide: AuthService,
          useValue: { requestPasswordReset },
        },
      ],
    }).compileComponents();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('有効なメールアドレスでsubmitするとリクエストを送信し、成功メッセージを表示する', async () => {
    const fixture = TestBed.createComponent(RequestPasswordResetComponent);
    fixture.autoDetectChanges();

    const input = fixture.nativeElement.querySelector('input[type="email"]') as HTMLInputElement;
    input.value = 'user@example.com';
    input.dispatchEvent(new Event('input', { bubbles: true }));
    await fixture.whenStable();

    const form = fixture.nativeElement.querySelector('form') as HTMLFormElement;
    form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));

    expect(requestPasswordReset).toHaveBeenCalledWith({ email: 'user@example.com' });
    fixture.detectChanges();
    expect(fixture.nativeElement.textContent).toContain('送信中');

    response$.next({ message: 'メールを送信しました' });
    await fixture.whenStable();

    expect(fixture.nativeElement.textContent).toContain('メールを送信しました');
    expect(fixture.nativeElement.textContent).not.toContain('送信中');
  });

  it('リクエストが失敗するとエラーメッセージを表示する', async () => {
    vi.spyOn(console, 'error').mockImplementation(() => undefined);
    const fixture = TestBed.createComponent(RequestPasswordResetComponent);
    fixture.autoDetectChanges();

    const input = fixture.nativeElement.querySelector('input[type="email"]') as HTMLInputElement;
    input.value = 'user@example.com';
    input.dispatchEvent(new Event('input', { bubbles: true }));
    await fixture.whenStable();

    const form = fixture.nativeElement.querySelector('form') as HTMLFormElement;
    form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));
    response$.error(new Error('request failed'));
    await fixture.whenStable();

    expect(fixture.nativeElement.textContent).toContain('送信に失敗しました');
  });

  it('フォームが無効な場合はリクエストを送信しない', () => {
    const fixture = TestBed.createComponent(RequestPasswordResetComponent);
    fixture.autoDetectChanges();

    const form = fixture.nativeElement.querySelector('form') as HTMLFormElement;
    form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));

    expect(requestPasswordReset).not.toHaveBeenCalled();
  });
});
