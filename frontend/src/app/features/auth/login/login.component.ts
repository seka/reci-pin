import { Component, inject } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { AuthService, LoginRequest } from '../../../core/services/auth.service';
import { AuthCardComponent } from '../../../shared/components/organisms/auth-card/auth-card.component';
import { InputComponent } from '../../../shared/components/atoms/input/input.component';
import { ButtonComponent } from '../../../shared/components/atoms/button/button.component';

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
      <form (ngSubmit)="onSubmit()">
        <app-input
          label="メールアドレス"
          type="email"
          [(ngModel)]="credentials.email"
          name="email"
          [required]="true"
        ></app-input>

        <div style="margin-top: 16px;">
          <app-input
            label="パスワード"
            type="password"
            [(ngModel)]="credentials.password"
            name="password"
            [required]="true"
          ></app-input>
        </div>

        <div class="actions">
          <app-button variant="primary" type="submit" class="submit-btn">ログイン</app-button>
        </div>

        <div style="text-align: center; margin-bottom: 16px;">
          <a routerLink="/password-reset/request" style="font-size: 14px; color: #666; text-decoration: none;">
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
        margin-top: 24px;
        margin-bottom: 16px;
      }
      .submit-btn {
        width: 100%;
      }
      .full-width-btn {
        width: 100%;
      }
      .error {
        color: #f44336;
        margin-top: 16px;
        text-align: center;
      }
    `,
  ],
})
export class LoginComponent {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  credentials: LoginRequest = { email: '', password: '' };
  errorMessage = '';

  onSubmit() {
    this.authService.login(this.credentials).subscribe({
      next: () => {
        this.router.navigate(['/recipes']);
      },
      error: () => {
        this.errorMessage = 'ログインに失敗しました';
      },
    });
  }
}
