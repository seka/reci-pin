import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
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
    CommonModule,
    FormsModule,
    RouterModule,
    ReactiveFormsModule,
    AuthCardComponent,
    InputComponent,
    ButtonComponent
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
        ></app-input>

        <div style="margin-top: 16px;">
          <app-input 
            label="メールアドレス" 
            type="email" 
            [(ngModel)]="user.email" 
            name="email" 
            [required]="true"
          ></app-input>
        </div>

        <div style="margin-top: 16px;">
          <app-input 
            label="パスワード" 
            type="password" 
            [(ngModel)]="user.password" 
            name="password" 
            [required]="true"
          ></app-input>
        </div>

        <div class="actions">
          <app-button variant="primary" type="submit" class="submit-btn">登録</app-button>
        </div>

        <p *ngIf="errorMessage" class="error">{{ errorMessage }}</p>
      </form>

      <div footer>
        <a routerLink="/login" style="text-decoration: none;">
          <app-button variant="accent" type="button" class="full-width-btn">ログインページへ</app-button>
        </a>
      </div>
    </app-auth-card>
  `,
  styles: [`
    .actions { margin-top: 24px; margin-bottom: 16px; }
    .submit-btn { width: 100%; }
    .full-width-btn { width: 100%; }
    .error { color: #f44336; margin-top: 16px; text-align: center; }
  `]
})
export class SignupComponent {
  user: SignupRequest = { name: '', email: '', password: '' };
  errorMessage = '';

  constructor(
    private authService: AuthService,
    private router: Router
  ) { }

  onSubmit() {
    this.authService.signup(this.user).subscribe({
      next: (response) => {
        this.router.navigate(['/recipes']);
      },
      error: (err) => {
        this.errorMessage = '登録に失敗しました';
      }
    });
  }
}
