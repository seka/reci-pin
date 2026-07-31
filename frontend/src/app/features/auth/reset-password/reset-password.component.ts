import { Component, inject, signal } from '@angular/core';
import {
  disabled,
  form,
  FormField,
  FormRoot,
  maxLength,
  minLength,
  required,
} from '@angular/forms/signals';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { TranslocoPipe, TranslocoService } from '@jsverse/transloco';
import { AuthService } from '../../../core/services/auth.service';
import { PasswordResetConfirmFormModel } from '../../../core/models/auth.model';
import { AuthCardComponent } from '../../../shared/components/organisms/auth-card/auth-card.component';
import { InputComponent } from '../../../shared/components/atoms/input/input.component';
import { ButtonComponent } from '../../../shared/components/atoms/button/button.component';
import { VALIDATION_RULES } from '../../../core/constants/validation.constants';
import { AlertComponent } from '../../../shared/components/atoms/alert/alert.component';

@Component({
  selector: 'app-reset-password',
  imports: [
    FormField,
    FormRoot,
    RouterModule,
    TranslocoPipe,
    AuthCardComponent,
    InputComponent,
    ButtonComponent,
    AlertComponent,
  ],
  templateUrl: './reset-password.component.html',
  styleUrls: ['./reset-password.component.scss'],
})
export class ResetPasswordComponent {
  private readonly authService = inject(AuthService);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly translate = inject(TranslocoService);

  private readonly model = signal<PasswordResetConfirmFormModel>({ token: '', newPassword: '' });

  protected readonly form = form(this.model, (path) => {
    required(path.token);

    required(path.newPassword);
    minLength(path.newPassword, VALIDATION_RULES.PASSWORD.MIN_LENGTH);
    maxLength(path.newPassword, 200);
    disabled(path.newPassword, { when: (ctx) => !ctx.valueOf(path.token) });
  });

  protected readonly message = signal('');
  protected readonly errorMessage = signal('');
  protected readonly isLoading = signal(false);
  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  constructor() {
    this.route.queryParams.subscribe((params) => {
      const token = params['token'] || '';
      if (token) {
        this.model.update((model) => ({ ...model, token }));
      } else {
        this.errorMessage.set(
          this.translate.translate('FEATURES.AUTH.RESET_PASSWORD.INVALID_LINK'),
        );
      }
    });
  }

  onSubmit() {
    if (this.form().invalid()) return;

    this.isLoading.set(true);
    this.message.set('');
    this.errorMessage.set('');

    this.authService.resetPassword(this.form().value()).subscribe({
      next: (res) => {
        this.message.set(res.message);
        this.isLoading.set(false);
        setTimeout(() => {
          this.router.navigate(['/login']);
        }, 3000);
      },
      error: (err) => {
        this.errorMessage.set(
          this.translate.translate('FEATURES.AUTH.RESET_PASSWORD.FAILED_EXPIRED'),
        );
        this.isLoading.set(false);
        console.error(err);
      },
    });
  }
}
