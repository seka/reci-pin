import { Component, inject } from '@angular/core';
import { AsyncPipe } from '@angular/common';
import { RouterModule } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';
import { ButtonComponent } from '../../shared/components/atoms/button/button.component';
import { HeadlineComponent } from '../../shared/components/atoms/headline/headline.component';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [AsyncPipe, RouterModule, ButtonComponent, HeadlineComponent],
  template: `
    <div class="settings-container">
      <div class="back-link">
        <a routerLink="/recipes">← レシピ一覧に戻る</a>
      </div>

      <app-headline level="1">設定</app-headline>

      <section class="settings-section">
        <app-headline level="2">アカウント</app-headline>

        @if (currentUser$ | async; as user) {
          <div class="user-info">
            <p><strong>名前:</strong> {{ user.name }}</p>
            <p><strong>メールアドレス:</strong> {{ user.email }}</p>
          </div>
        }
        <div class="account-actions">
          <a routerLink="/settings/password" style="text-decoration: none;">
            <app-button variant="outline">パスワードを変更する</app-button>
          </a>
        </div>
      </section>

      <section class="settings-section danger-zone">
        <app-headline level="2">退会</app-headline>
        <p class="warning-text">
          アカウントを削除すると、すべてのレシピやデータが完全に削除されます。
          この操作は取り消せません。
        </p>
        <app-button variant="warn" (click)="onWithdraw()" [disabled]="isProcessing">
          {{ isProcessing ? '処理中...' : 'アカウントを削除する' }}
        </app-button>
      </section>
    </div>
  `,
  styles: [
    `
      .settings-container {
        max-width: 600px;
        margin: 0 auto;
        padding: var(--spacing-3);
      }

      .settings-section {
        margin-top: var(--spacing-4);
        padding: var(--spacing-3);
        background: var(--color-surface);
        border-radius: var(--radius-2);
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      }

      .user-info p {
        margin: 8px 0;
        color: var(--color-text-secondary);
      }

      .account-actions {
        margin-top: var(--spacing-2);
      }

      .danger-zone {
        border: 1px solid var(--color-error);
      }

      .warning-text {
        color: var(--color-text-secondary);
        margin-bottom: var(--spacing-2);
        line-height: 1.6;
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
    `,
  ],
})
export class SettingsComponent {
  private readonly authService = inject(AuthService);

  currentUser$ = this.authService.currentUser$;
  isProcessing = false;

  onWithdraw(): void {
    if (!confirm('本当に退会しますか？\n\nすべてのデータが削除され、この操作は取り消せません。')) {
      return;
    }

    this.isProcessing = true;
    this.authService.withdraw().subscribe({
      next: () => {
        alert('退会が完了しました。ご利用ありがとうございました。');
      },
      error: (err: Error) => {
        this.isProcessing = false;
        alert('退会処理に失敗しました。しばらくしてから再度お試しください。');
        console.error('Withdraw error:', err);
      },
    });
  }
}
