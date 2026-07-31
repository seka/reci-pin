/**
 * @vitest-environment jsdom
 */
import { TestBed } from '@angular/core/testing';
import { provideRouter, Router } from '@angular/router';
import { TranslocoTestingModule } from '@jsverse/transloco';
import { Subject } from 'rxjs';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { AuthService } from '../../../core/services/auth.service';
import { SignupComponent } from './signup.component';

describe('SignupComponent', () => {
  let response$: Subject<unknown>;
  const signup = vi.fn();

  beforeEach(async () => {
    response$ = new Subject<unknown>();
    signup.mockReset();
    signup.mockReturnValue(response$.asObservable());

    await TestBed.configureTestingModule({
      imports: [
        SignupComponent,
        TranslocoTestingModule.forRoot({
          langs: {
            ja: {
              FEATURES: { AUTH: { SIGNUP: { FAILED: '登録に失敗しました' } } },
              VALIDATION: {
                INVALID_EMAIL: 'メールアドレスが不正です',
                INVALID_INPUT: '入力内容が不正です',
              },
            },
          },
          translocoConfig: { availableLangs: ['ja'], defaultLang: 'ja' },
        }),
      ],
      providers: [provideRouter([]), { provide: AuthService, useValue: { signup } }],
    }).compileComponents();
  });

  function fillForm(element: HTMLElement): void {
    const inputs = element.querySelectorAll('input');
    const values = ['テストユーザー', 'user@example.com', 'password1'];
    inputs.forEach((input, index) => {
      input.value = values[index];
      input.dispatchEvent(new Event('input', { bubbles: true }));
    });
  }

  it('有効なフォームをsubmitすると登録APIを呼び、成功後に遷移する', async () => {
    const navigate = vi.spyOn(TestBed.inject(Router), 'navigate').mockResolvedValue(true);
    const fixture = TestBed.createComponent(SignupComponent);
    fixture.autoDetectChanges();
    fillForm(fixture.nativeElement);
    await fixture.whenStable();

    const form = fixture.nativeElement.querySelector('form') as HTMLFormElement;
    form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));

    expect(signup).toHaveBeenCalledWith({
      name: 'テストユーザー',
      email: 'user@example.com',
      password: 'password1',
    });

    response$.next({});
    await fixture.whenStable();
    expect(navigate).toHaveBeenCalledWith(['/recipes']);
  });

  it('サーバーのフィールドエラーを対応する入力欄へ表示する', async () => {
    const fixture = TestBed.createComponent(SignupComponent);
    fixture.autoDetectChanges();
    fillForm(fixture.nativeElement);
    await fixture.whenStable();

    const form = fixture.nativeElement.querySelector('form') as HTMLFormElement;
    form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));
    response$.error({
      error: {
        error: {
          message: 'validation failed',
          details: {
            email: [{ code: 'EMAIL_INVALID_FORMAT' }],
          },
        },
      },
    });
    await fixture.whenStable();

    expect(fixture.nativeElement.textContent).toContain('メールアドレスが不正です');
  });

  it('フォームが無効な場合は登録APIを呼ばない', () => {
    const fixture = TestBed.createComponent(SignupComponent);
    fixture.autoDetectChanges();

    const form = fixture.nativeElement.querySelector('form') as HTMLFormElement;
    form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));

    expect(signup).not.toHaveBeenCalled();
  });
});
