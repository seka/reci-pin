import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { AuthService, LoginRequest } from '../../core/services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  template: `
    <div class="login-container">
      <h2>ログイン</h2>
      <form (ngSubmit)="onSubmit()">
        <div>
          <label>メールアドレス</label>
          <input type="email" [(ngModel)]="credentials.email" name="email" required />
        </div>
        <div>
          <label>パスワード</label>
          <input type="password" [(ngModel)]="credentials.password" name="password" required />
        </div>
        <button type="submit">ログイン</button>
        <p *ngIf="errorMessage" class="error">{{ errorMessage }}</p>
      </form>
      <a routerLink="/signup">アカウント作成</a>
    </div>
  `,
  styles: [`
    .login-container { max-width: 400px; margin: 50px auto; padding: 20px; }
    input { width: 100%; padding: 8px; margin: 5px 0 15px; }
    button { width: 100%; padding: 10px; background: #007bff; color: white; border: none; cursor: pointer; }
    .error { color: red; }
  `]
})
export class LoginComponent {
  credentials: LoginRequest = { email: '', password: '' };
  errorMessage = '';

  constructor(
    private authService: AuthService,
    private router: Router
  ) { }

  onSubmit() {
    this.authService.login(this.credentials).subscribe({
      next: (response) => {
        this.authService.saveToken(response.token);
        this.router.navigate(['/recipes']);
      },
      error: (err) => {
        this.errorMessage = 'ログインに失敗しました';
      }
    });
  }
}
