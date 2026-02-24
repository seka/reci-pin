import { Component, inject } from '@angular/core';
import { AsyncPipe } from '@angular/common';
import { RouterModule } from '@angular/router';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { AuthService } from '../../core/services/auth.service';
import { ButtonComponent } from '../../shared/components/atoms/button/button.component';
import { HeadlineComponent } from '../../shared/components/atoms/headline/headline.component';
import { LinkComponent } from '../../shared/components/atoms/link/link.component';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [AsyncPipe, RouterModule, TranslatePipe, ButtonComponent, HeadlineComponent, LinkComponent],
  template: `
    <div class="settings-container">
      <div class="back-link">
        <app-link routerLink="/recipes">{{ 'SETTINGS.BACK_TO_RECIPES' | translate }}</app-link>
      </div>

      <app-headline level="1">{{ 'SETTINGS.TITLE' | translate }}</app-headline>

      <section class="settings-section">
        <app-headline level="2">{{ 'SETTINGS.ACCOUNT' | translate }}</app-headline>

        @if (currentUser$ | async; as user) {
          <div class="user-info">
            <p><strong>{{ 'SETTINGS.NAME' | translate }}:</strong> {{ user.name }}</p>
            <p><strong>{{ 'SETTINGS.EMAIL' | translate }}:</strong> {{ user.email }}</p>
          </div>
        }
        <div class="account-actions">
          <app-button routerLink="/settings/password" variant="outline">{{ 'SETTINGS.CHANGE_PASSWORD_BUTTON' | translate }}</app-button>
        </div>
      </section>

      <section class="settings-section danger-zone">
        <app-headline level="2">{{ 'SETTINGS.WITHDRAW_TITLE' | translate }}</app-headline>
        <p class="warning-text" [innerHTML]="'SETTINGS.WITHDRAW_WARNING' | translate">
        </p>
        <app-button variant="warn" (click)="onWithdraw()" [disabled]="isProcessing">
          {{ isProcessing ? ('SETTINGS.WITHDRAWING' | translate) : ('SETTINGS.WITHDRAW_BUTTON' | translate) }}
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
    `,
  ],
})
export class SettingsComponent {
  private readonly authService = inject(AuthService);
  private readonly translate = inject(TranslateService);

  currentUser$ = this.authService.currentUser$;
  isProcessing = false;

  onWithdraw(): void {
    if (!confirm(this.translate.instant('SETTINGS.WITHDRAW_CONFIRM'))) {
      return;
    }

    this.isProcessing = true;
    this.authService.withdraw().subscribe({
      next: () => {
        alert(this.translate.instant('SETTINGS.WITHDRAW_SUCCESS'));
      },
      error: (err: Error) => {
        this.isProcessing = false;
        alert(this.translate.instant('SETTINGS.WITHDRAW_FAILED'));
        console.error('Withdraw error:', err);
      },
    });
  }
}
