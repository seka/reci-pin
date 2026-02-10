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
        ></app-input>
        <p class="hint">8文字以上の英数字を入力してください。</p>

        <app-input
          label="新しいパスワード（確認）"
          type="password"
          formControlName="confirmNewPassword"
          [errorMessage]="getErrorMessage('confirmNewPassword')"
          [required]="true"
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
        padding: 24px;
      }

      .back-link {
        margin-bottom: 16px;
      }

      .back-link a {
        color: var(--primary-color, #1976d2);
        text-decoration: none;
      }

      .back-link a:hover {
        text-decoration: underline;
      }

      .password-form {
        margin-top: 16px;
        display: flex;
        flex-direction: column;
        gap: 12px;
        background: var(--surface-color, #fff);
        padding: 24px;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      }

      .hint {
        font-size: 0.85rem;
        color: var(--text-secondary, #666);
        margin-top: -8px;
        margin-left: 4px;
      }

      .forgot-password-link {
        font-size: 0.85rem;
        margin-top: -8px;
        margin-left: 4px;
        text-align: left;
      }

      .forgot-password-link a {
        color: var(--primary-color, #1976d2);
        text-decoration: none;
      }

      .forgot-password-link a:hover {
        text-decoration: underline;
      }

      .actions {
        margin-top: 8px;
      }

      .error-global {
        color: #f44336;
        margin-top: 12px;
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

  form = this.fb.group(
    {
      currentPassword: ['', Validators.required],
      newPassword: ['', [Validators.required, Validators.minLength(8)]],
      confirmNewPassword: ['', Validators.required],
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
      if (control.errors['minlength']) return '8文字以上で入力してください';
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
