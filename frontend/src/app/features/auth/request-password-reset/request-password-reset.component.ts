import { Component, inject, signal } from '@angular/core';
import { email, form, FormField, FormRoot, maxLength, required } from '@angular/forms/signals';
import { RouterModule } from '@angular/router';
import { TranslocoPipe, TranslocoService } from '@jsverse/transloco';
import { AuthService } from '../../../core/services/auth.service';
import { PasswordResetFormModel } from '../../../core/models/auth.model';
import { AuthCardComponent } from '../../../shared/components/organisms/auth-card/auth-card.component';
import { InputComponent } from '../../../shared/components/atoms/input/input.component';
import { ButtonComponent } from '../../../shared/components/atoms/button/button.component';
import { VALIDATION_RULES } from '../../../core/constants/validation.constants';
import { AlertComponent } from '../../../shared/components/atoms/alert/alert.component';

@Component({
  selector: 'app-request-password-reset',
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
  templateUrl: './request-password-reset.component.html',
  styleUrls: ['./request-password-reset.component.scss'],
})
export class RequestPasswordResetComponent {
  private readonly authService = inject(AuthService);
  private readonly translate = inject(TranslocoService);

  private readonly model = signal<PasswordResetFormModel>({ email: '' });

  protected readonly form = form(this.model, (path) => {
    required(path.email);
    email(path.email);
    maxLength(path.email, 200);
  });

  protected readonly message = signal('');
  protected readonly errorMessage = signal('');
  protected readonly isLoading = signal(false);
  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  onSubmit() {
    if (this.form().invalid()) return;

    this.isLoading.set(true);
    this.message.set('');
    this.errorMessage.set('');

    this.authService.requestPasswordReset(this.form().value()).subscribe({
      next: (res) => {
        this.message.set(res.message);
        this.isLoading.set(false);
      },
      error: (err) => {
        this.errorMessage.set(
          this.translate.translate('FEATURES.AUTH.REQUEST_PASSWORD_RESET.FAILED'),
        );
        this.isLoading.set(false);
        console.error(err);
      },
    });
  }
}
