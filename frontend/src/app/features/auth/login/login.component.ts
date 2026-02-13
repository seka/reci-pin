import { Component, inject, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';
import { AuthCardComponent } from '../../../shared/components/organisms/auth-card/auth-card.component';
import { InputComponent } from '../../../shared/components/atoms/input/input.component';
import { ButtonComponent } from '../../../shared/components/atoms/button/button.component';
import { VALIDATION_RULES } from '../../../core/constants/validation.constants';

@Component({
  selector: 'app-login',
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
    <app-auth-card title="ログイン">
      <form [formGroup]="loginForm" (ngSubmit)="onSubmit()">
        <app-input
          label="メールアドレス"
          type="email"
          formControlName="email"
          [required]="true"
          [maxLength]="VALIDATION_RULES.EMAIL.MAX_LENGTH"
        ></app-input>

        <div style="margin-top: var(--spacing-2);">
          <app-input
            label="パスワード"
            type="password"
            formControlName="password"
            [required]="true"
            [maxLength]="VALIDATION_RULES.PASSWORD.MAX_LENGTH"
          ></app-input>
        </div>

        <div class="actions">
          <app-button variant="primary" type="submit" class="submit-btn" [disabled]="loginForm.invalid">ログイン</app-button>
        </div>

        <div style="text-align: center; margin-bottom: var(--spacing-2);">
          <a routerLink="/password-reset/request" style="font-size: var(--font-size-2); color: var(--color-text-secondary); text-decoration: none;">
            パスワードを忘れた場合
          </a>
        </div>

        @if (errorMessage) {
          <p class="error">{{ errorMessage }}</p>
        }
      </form>

      <div footer>
        <a routerLink="/signup" style="text-decoration: none;">
          <app-button variant="accent" type="button" class="full-width-btn"
            >アカウント作成はこちら</app-button
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
        text-align: center;
      }
    `,
  ],
})
export class LoginComponent implements OnInit {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);

  loginForm!: FormGroup;
  errorMessage = '';

  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  ngOnInit() {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email, Validators.maxLength(VALIDATION_RULES.EMAIL.MAX_LENGTH)]],
      password: ['', [Validators.required, Validators.maxLength(VALIDATION_RULES.PASSWORD.MAX_LENGTH)]],
    });
  }

  onSubmit() {
    if (this.loginForm.invalid) {
      return;
    }

    this.authService.login(this.loginForm.value).subscribe({
      next: () => {
        this.router.navigate(['/recipes']);
      },
      error: () => {
        this.errorMessage = 'ログインに失敗しました';
      },
    });
  }
}
