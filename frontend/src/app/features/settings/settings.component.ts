import { Component, inject } from '@angular/core';
import { AsyncPipe } from '@angular/common';
import { RouterModule } from '@angular/router';
import { TranslocoPipe, TranslocoService } from '@jsverse/transloco';
import { AuthService } from '../../core/services/auth.service';
import { ButtonComponent } from '../../shared/components/atoms/button/button.component';
import { HeadlineComponent } from '../../shared/components/atoms/headline/headline.component';
import { LinkComponent } from '../../shared/components/atoms/link/link.component';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [
    AsyncPipe,
    RouterModule,
    TranslocoPipe,
    ButtonComponent,
    HeadlineComponent,
    LinkComponent,
  ],
  templateUrl: './settings.component.html',
  styleUrl: './settings.component.scss',
})
export class SettingsComponent {
  private readonly authService = inject(AuthService);
  private readonly translate = inject(TranslocoService);

  currentUser$ = this.authService.currentUser$;
  isProcessing = false;

  onWithdraw(): void {
    if (!confirm(this.translate.translate('FEATURES.SETTINGS.SETTINGS.WITHDRAW_CONFIRM'))) {
      return;
    }

    this.isProcessing = true;
    this.authService.withdraw().subscribe({
      next: () => {
        alert(this.translate.translate('FEATURES.SETTINGS.SETTINGS.WITHDRAW_SUCCESS'));
      },
      error: (err: Error) => {
        this.isProcessing = false;
        alert(this.translate.translate('FEATURES.SETTINGS.SETTINGS.WITHDRAW_FAILED'));
        console.error('Withdraw error:', err);
      },
    });
  }
}
