import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { AuthService, SignupRequest } from '../../../core/services/auth.service';

@Component({
  selector: 'app-signup',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    RouterModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule
  ],
  template: `
    <div class="signup-container">
      <mat-card>
        <mat-card-header>
          <mat-card-title>アカウント作成</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <form (ngSubmit)="onSubmit()">
            <mat-form-field appearance="outline" class="full-width">
              <mat-label>名前</mat-label>
              <input matInput type="text" [(ngModel)]="user.name" name="name" required />
            </mat-form-field>

            <mat-form-field appearance="outline" class="full-width">
              <mat-label>メールアドレス</mat-label>
              <input matInput type="email" [(ngModel)]="user.email" name="email" required />
            </mat-form-field>

            <mat-form-field appearance="outline" class="full-width">
              <mat-label>パスワード</mat-label>
              <input matInput type="password" [(ngModel)]="user.password" name="password" required />
            </mat-form-field>

            <div class="actions">
              <button mat-flat-button color="primary" type="submit" class="submit-btn">登録</button>
            </div>

            <p *ngIf="errorMessage" class="error">{{ errorMessage }}</p>
          </form>
        </mat-card-content>
        <mat-card-footer>
           <a mat-button color="accent" routerLink="/login">ログインページへ</a>
        </mat-card-footer>
      </mat-card>
    </div>
  `,
  styles: [`
    .signup-container { max-width: 400px; margin: 80px auto; padding: 0 20px; }
    mat-card { padding: 24px; text-align: center; }
    mat-card-title { margin-bottom: 24px; font-weight: 700; color: #e91e63; font-size: 1.5rem; justify-content: center; }
    .full-width { width: 100%; margin-bottom: 8px; }
    .actions { margin-top: 16px; margin-bottom: 16px; }
    .submit-btn { width: 100%; font-size: 1.1em; padding: 24px 0; }
    mat-card-footer { padding: 16px; margin-top: 0; }
    .error { color: #f44336; margin-top: 16px; }
    a { width: 100%; }
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
