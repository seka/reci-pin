import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService, SignupRequest } from '../../core/services/auth.service';

@Component({
    selector: 'app-signup',
    standalone: true,
    imports: [CommonModule, FormsModule],
    template: `
    <div class="signup-container">
      <h2>アカウント作成</h2>
      <form (ngSubmit)="onSubmit()">
        <div>
          <label>名前</label>
          <input type="text" [(ngModel)]="user.name" name="name" required />
        </div>
        <div>
          <label>メールアドレス</label>
          <input type="email" [(ngModel)]="user.email" name="email" required />
        </div>
        <div>
          <label>パスワード</label>
          <input type="password" [(ngModel)]="user.password" name="password" required />
        </div>
        <button type="submit">登録</button>
        <p *ngIf="errorMessage" class="error">{{ errorMessage }}</p>
      </form>
      <a routerLink="/login">ログインページへ</a>
    </div>
  `,
    styles: [`
    .signup-container { max-width: 400px; margin: 50px auto; padding: 20px; }
    input { width: 100%; padding: 8px; margin: 5px 0 15px; }
    button { width: 100%; padding: 10px; background: #28a745; color: white; border: none; cursor: pointer; }
    .error { color: red; }
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
                this.authService.saveToken(response.token);
                this.router.navigate(['/recipes']);
            },
            error: (err) => {
                this.errorMessage = '登録に失敗しました';
            }
        });
    }
}
