import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  FormBuilder,
  ReactiveFormsModule,
  Validators,
  AbstractControl,
  ValidationErrors,
} from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';
import { InputComponent } from '../../../shared/components/atoms/input/input.component';
import { ButtonComponent } from '../../../shared/components/atoms/button/button.component';
import { HeadlineComponent } from '../../../shared/components/atoms/headline/headline.component';
import { VALIDATION_RULES } from '../../../core/constants/validation.constants';

@Component({
  selector: 'app-change-password',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    InputComponent,
    ButtonComponent,
    HeadlineComponent,
  ],
  template: `
    <div class="change-password-container">
      <div class="back-link">
        <a routerLink="/settings">← 設定に戻る</a>
      </div>

      <app-headline level="1">パスワード変更</app-headline>

      <form [formGroup]="form" (ngSubmit)="onSubmit()" class="password-form">
        <app-input
          label="現在のパスワード"
          type="password"
          formControlName="currentPassword"
          [errorMessage]="getErrorMessage('currentPassword')"
          [required]="true"
          [maxLength]="VALIDATION_RULES.PASSWORD.MAX_LENGTH"
        ></app-input>
        <p class="forgot-password-link">
          <a routerLink="/password-reset/request">パスワードを忘れた場合</a>
        </p>

        <app-input
          label="新しいパスワード"
          type="password"
          formControlName="newPassword"
          [errorMessage]="getErrorMessage('newPassword')"
          [required]="true"
          [maxLength]="VALIDATION_RULES.PASSWORD.MAX_LENGTH"
          [showCounter]="true"
        ></app-input>
        <p class="hint">{{ VALIDATION_RULES.PASSWORD.MIN_LENGTH }}文字以上の英数字を入力してください。</p>

        <app-input
          label="新しいパスワード（確認）"
          type="password"
          formControlName="confirmNewPassword"
          [errorMessage]="getErrorMessage('confirmNewPassword')"
          [required]="true"
          [maxLength]="VALIDATION_RULES.PASSWORD.MAX_LENGTH"
        ></app-input>

        <div class="actions">
          <app-button type="submit" [disabled]="form.invalid || isProcessing">
            {{ isProcessing ? '変更中...' : 'パスワードを変更する' }}
          </app-button>
        </div>

        <p *ngIf="errorMessage" class="error-global">{{ errorMessage }}</p>
      </form>
    </div>
  `,
  styles: [
    `
      .change-password-container {
        max-width: 600px;
        margin: 0 auto;
        padding: var(--spacing-3);
      }

      .back-link {
        margin-bottom: var(--spacing-2);
      }

      .back-link a {
        color: var(--color-primary);
        text-decoration: none;
      }

      .back-link a:hover {
        text-decoration: underline;
      }

      .password-form {
        margin-top: var(--spacing-2);
        display: flex;
        flex-direction: column;
        gap: var(--spacing-1_5);
        background: var(--color-surface);
        padding: var(--spacing-3);
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      }

      .hint {
        font-size: var(--font-size-2);
        color: var(--color-text-secondary);
        margin-top: -8px;
        margin-left: 4px;
      }

      .forgot-password-link {
        font-size: var(--font-size-2);
        margin-top: -8px;
        margin-left: 4px;
        text-align: left;
      }

      .forgot-password-link a {
        color: var(--color-primary);
        text-decoration: none;
      }

      .forgot-password-link a:hover {
        text-decoration: underline;
      }

      .actions {
        margin-top: var(--spacing-1);
      }

      .error-global {
        color: var(--color-error);
        margin-top: var(--spacing-1_5);
        font-weight: bold;
      }
    `,
  ],
})
export class ChangePasswordComponent {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);

  isProcessing = false;
  errorMessage = '';

  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  form = this.fb.group(
    {
      currentPassword: ['', [Validators.required, Validators.maxLength(VALIDATION_RULES.PASSWORD.MAX_LENGTH)]],
      newPassword: [
        '',
        [
          Validators.required,
          Validators.minLength(VALIDATION_RULES.PASSWORD.MIN_LENGTH),
          Validators.maxLength(VALIDATION_RULES.PASSWORD.MAX_LENGTH),
        ],
      ],
      confirmNewPassword: ['', [Validators.required, Validators.maxLength(VALIDATION_RULES.PASSWORD.MAX_LENGTH)]],
    },
    { validators: this.passwordMatchValidator }
  );

  private passwordMatchValidator(control: AbstractControl): ValidationErrors | null {
    const newPassword = control.get('newPassword')?.value;
    const confirmPassword = control.get('confirmNewPassword')?.value;
    return newPassword === confirmPassword ? null : { mismatch: true };
  }

  getErrorMessage(controlName: string): string | null {
    const control = this.form.get(controlName);
    if (control?.touched && control?.errors) {
      if (control.errors['required']) return '必須項目です';
      if (control.errors['minlength']) return `${VALIDATION_RULES.PASSWORD.MIN_LENGTH}文字以上で入力してください`;
      if (control.errors['maxlength']) return `${VALIDATION_RULES.PASSWORD.MAX_LENGTH}文字以内で入力してください`;
      if (controlName === 'confirmNewPassword' && this.form.errors?.['mismatch']) {
        return 'パスワードが一致しません';
      }
    }
    return null;
  }

  onSubmit() {
    if (this.form.invalid) return;

    this.isProcessing = true;
    this.errorMessage = '';

    const { currentPassword, newPassword } = this.form.value;

    if (!currentPassword || !newPassword) return;

    this.authService
      .changePassword({
        current_password: currentPassword,
        new_password: newPassword,
      })
      .subscribe({
        next: () => {
          alert('パスワードを変更しました。');
          this.router.navigate(['/settings']);
        },
        error: (err) => {
          this.isProcessing = false;
          // Backend returns bad request for incorrect password or validation errors
          if (err.status === 400 || err.status === 401) {
            this.errorMessage = '現在のパスワードが間違っているか、新しいパスワードが無効です。';
          } else {
            this.errorMessage = 'パスワードの変更に失敗しました。時間をおいて再度お試しください。';
          }
          console.error(err);
        },
      });
  }
}
