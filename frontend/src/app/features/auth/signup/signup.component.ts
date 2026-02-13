import { Component, inject } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { AuthService, SignupRequest } from '../../../core/services/auth.service';
import { AuthCardComponent } from '../../../shared/components/organisms/auth-card/auth-card.component';
import { InputComponent } from '../../../shared/components/atoms/input/input.component';
import { ButtonComponent } from '../../../shared/components/atoms/button/button.component';

@Component({
  selector: 'app-signup',
  standalone: true,
  imports: [
    FormsModule,
    RouterModule,
    ReactiveFormsModule,
    AuthCardComponent,
    InputComponent,
    ButtonComponent,
  ],
  template: `
    <app-auth-card title="アカウント作成">
      <form (ngSubmit)="onSubmit()">
          <app-input
            label="名前"
            type="text"
            [(ngModel)]="user.name"
            name="name"
            [required]="true"
            [errorMessage]="fieldErrors['name']"
          ></app-input>

        <div style="margin-top: var(--spacing-2);">
          <app-input
            label="メールアドレス"
            type="email"
            [(ngModel)]="user.email"
            name="email"
            [required]="true"
            [errorMessage]="fieldErrors['email']"
          ></app-input>
        </div>

        <div style="margin-top: var(--spacing-2);">
          <app-input
            label="パスワード"
            type="password"
            [(ngModel)]="user.password"
            name="password"
            [required]="true"
            [errorMessage]="fieldErrors['password']"
          ></app-input>
        </div>

        <div class="actions">
          <app-button variant="primary" type="submit" class="submit-btn">登録</app-button>
        </div>

        @if (errorMessage) {
          <p class="error">{{ errorMessage }}</p>
        }
      </form>

      <div footer>
        <a routerLink="/login" style="text-decoration: none;">
          <app-button variant="accent" type="button" class="full-width-btn"
            >ログインページへ</app-button
          >
        </a>
      </div>
    </app-auth-card>
  `,
  styles: [
    `
      .actions {
        margin-top: var(--spacing-3);
        margin-bottom: var(--spacing-2);
      }
      .submit-btn {
        width: 100%;
      }
      .full-width-btn {
        width: 100%;
      }
      .error {
        color: var(--color-error);
        margin-top: var(--spacing-2);
        text-align: left !important;
        width: 100%;
      }
    `,
  ],
})
export class SignupComponent {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  user: SignupRequest = { name: '', email: '', password: '' };
  fieldErrors: { [key: string]: string[] } = {};
  errorMessage = '';

  onSubmit() {
    this.fieldErrors = {}; // Reset errors
    this.errorMessage = '';



    this.authService.signup(this.user).subscribe({
      next: () => {
        this.router.navigate(['/recipes']);
      },
      error: (err) => {
        // バックエンドからのエラーレスポンスを処理
        if (err.error && err.error.error && err.error.error.details) {
          const details = err.error.error.details;

          // 各フィールドのエラーをマッピング
          Object.keys(details).forEach(field => {
            const messages = (details as any)[field].map((d: any) => {
              switch (d.code) {
                case 'PASSWORD_TOO_SHORT':
                  return `パスワードは${d.params?.min || 8}文字以上である必要があります`;
                case 'PASSWORD_NO_ALPHA':
                  return 'パスワードには少なくとも1つの英字を含める必要があります';
                case 'PASSWORD_NO_NUMERIC':
                  return 'パスワードには少なくとも1つの数字を含める必要があります';
                case 'EMAIL_INVALID_FORMAT':
                  return 'メールアドレスの形式が正しくありません';
                case 'EMAIL_TOO_LONG':
                  return `メールアドレスは${d.params?.max || 254}文字以下である必要があります`;
                case 'REQUIRED':
                  return 'この項目は必須です';
                default:
                  return '入力内容が正しくありません';
              }
            });
            this.fieldErrors[field] = messages;
          });

          if (Object.keys(this.fieldErrors).length > 0) {
            return;
          }
        }
        this.errorMessage = '登録に失敗しました';
      },
    });
  }
}
