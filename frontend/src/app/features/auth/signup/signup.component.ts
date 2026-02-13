import { Component, inject, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';
import { AuthCardComponent } from '../../../shared/components/organisms/auth-card/auth-card.component';
import { InputComponent } from '../../../shared/components/atoms/input/input.component';
import { ButtonComponent } from '../../../shared/components/atoms/button/button.component';
import { VALIDATION_RULES } from '../../../core/constants/validation.constants';

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
      <form [formGroup]="signupForm" (ngSubmit)="onSubmit()">
          <app-input
            label="名前"
            type="text"
            formControlName="name"
            [required]="true"
            [maxLength]="VALIDATION_RULES.RECIPE.NAME_MAX_LENGTH"
            [showCounter]="true"
            [errorMessage]="fieldErrors['name']"
          ></app-input>

        <div style="margin-top: var(--spacing-2);">
          <app-input
            label="メールアドレス"
            type="email"
            formControlName="email"
            [required]="true"
            [maxLength]="VALIDATION_RULES.EMAIL.MAX_LENGTH"
            [showCounter]="true"
            [errorMessage]="fieldErrors['email']"
          ></app-input>
        </div>

        <div style="margin-top: var(--spacing-2);">
          <app-input
            label="パスワード"
            type="password"
            formControlName="password"
            [required]="true"
            [maxLength]="VALIDATION_RULES.PASSWORD.MAX_LENGTH"
            [showCounter]="true"
            [errorMessage]="fieldErrors['password']"
          ></app-input>
        </div>

        <div class="actions">
          <app-button variant="primary" type="submit" class="submit-btn" [disabled]="signupForm.invalid">登録</app-button>
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
      .submit-btn:disabled {
        opacity: 0.6;
        cursor: not-allowed;
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
export class SignupComponent implements OnInit {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);

  signupForm!: FormGroup;
  fieldErrors: { [key: string]: string[] } = {};
  errorMessage = '';

  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  ngOnInit() {
    this.signupForm = this.fb.group({
      name: ['', [Validators.required, Validators.maxLength(VALIDATION_RULES.RECIPE.NAME_MAX_LENGTH)]],
      email: ['', [Validators.required, Validators.email, Validators.maxLength(VALIDATION_RULES.EMAIL.MAX_LENGTH)]],
      password: ['', [
        Validators.required,
        Validators.minLength(VALIDATION_RULES.PASSWORD.MIN_LENGTH),
        Validators.maxLength(VALIDATION_RULES.PASSWORD.MAX_LENGTH),
        Validators.pattern(/[a-zA-Z]/),
        Validators.pattern(/[0-9]/)
      ]],
    });
  }

  onSubmit() {
    if (this.signupForm.invalid) {
      return;
    }

    this.fieldErrors = {}; // Reset errors
    this.errorMessage = '';

    this.authService.signup(this.signupForm.value).subscribe({
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
                case 'PASSWORD_TOO_LONG':
                  return `パスワードは${d.params?.max || VALIDATION_RULES.PASSWORD.MAX_LENGTH}文字以下である必要があります`;
                case 'PASSWORD_NO_ALPHA':
                  return 'パスワードには少なくとも1つの英字を含める必要があります';
                case 'PASSWORD_NO_NUMERIC':
                  return 'パスワードには少なくとも1つの数字を含める必要があります';
                case 'EMAIL_INVALID_FORMAT':
                  return 'メールアドレスの形式が正しくありません';
                case 'EMAIL_TOO_LONG':
                  return `メールアドレスは${d.params?.max || VALIDATION_RULES.EMAIL.MAX_LENGTH}文字以下である必要があります`;
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
