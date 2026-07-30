/**
 * @vitest-environment jsdom
 */
import { TestBed } from '@angular/core/testing';
import { provideRouter, Router } from '@angular/router';
import { TranslocoTestingModule } from '@jsverse/transloco';
import { Subject } from 'rxjs';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { AuthService } from '../../../core/services/auth.service';
import { ChangePasswordComponent } from './change-password.component';

describe('ChangePasswordComponent', () => {
  let response$: Subject<unknown>;
  const changePassword = vi.fn();

  beforeEach(async () => {
    response$ = new Subject<unknown>();
    changePassword.mockReset();
    changePassword.mockReturnValue(response$.asObservable());

    await TestBed.configureTestingModule({
      imports: [
        ChangePasswordComponent,
        TranslocoTestingModule.forRoot({
          langs: {
            ja: {
              FEATURES: {
                SETTINGS: {
                  CHANGE_PASSWORD: {
                    CHANGING: '変更中',
                    FAILED_INVALID: '現在のパスワードが正しくありません',
                    FAILED_ERROR: '変更に失敗しました',
                    SUCCESS: '変更しました',
                  },
                },
              },
            },
          },
          translocoConfig: { availableLangs: ['ja'], defaultLang: 'ja' },
        }),
      ],
      providers: [provideRouter([]), { provide: AuthService, useValue: { changePassword } }],
    }).compileComponents();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  function fillForm(element: HTMLElement): void {
    const inputs = element.querySelectorAll<HTMLInputElement>('input[type="password"]');
    const values = ['current-password1', 'new-password1', 'new-password1'];
    inputs.forEach((input, index) => {
      input.value = values[index];
      input.dispatchEvent(new Event('input', { bubbles: true }));
    });
  }

  it('submitすると変更APIを呼び、処理中表示へ切り替える', async () => {
    vi.stubGlobal('alert', vi.fn());
    vi.spyOn(TestBed.inject(Router), 'navigate').mockResolvedValue(true);
    const fixture = TestBed.createComponent(ChangePasswordComponent);
    fixture.autoDetectChanges();
    fillForm(fixture.nativeElement);
    await fixture.whenStable();

    const form = fixture.nativeElement.querySelector('form') as HTMLFormElement;
    form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));

    expect(changePassword).toHaveBeenCalledWith({
      currentPassword: 'current-password1',
      newPassword: 'new-password1',
      confirmNewPassword: 'new-password1',
    });
    expect(fixture.nativeElement.textContent).toContain('変更中');
  });

  it('認証エラーをDOMへ表示する', async () => {
    vi.spyOn(console, 'error').mockImplementation(() => undefined);
    const fixture = TestBed.createComponent(ChangePasswordComponent);
    fixture.autoDetectChanges();
    fillForm(fixture.nativeElement);
    await fixture.whenStable();

    const form = fixture.nativeElement.querySelector('form') as HTMLFormElement;
    form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));
    response$.error({ status: 401 });
    await fixture.whenStable();

    expect(fixture.nativeElement.textContent).toContain('現在のパスワードが正しくありません');
  });

  it('フォームが無効な場合は変更APIを呼ばない', () => {
    const fixture = TestBed.createComponent(ChangePasswordComponent);
    fixture.autoDetectChanges();

    const form = fixture.nativeElement.querySelector('form') as HTMLFormElement;
    form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));

    expect(changePassword).not.toHaveBeenCalled();
  });
});
