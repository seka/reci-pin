import { Component, inject, signal } from '@angular/core';
import { email, form, FormField, FormRoot, maxLength, required } from '@angular/forms/signals';
import { Router, RouterModule } from '@angular/router';
import { TranslocoPipe, TranslocoService } from '@jsverse/transloco';
import { AuthService } from '../../../core/services/auth.service';
import { LoginFormModel } from '../../../core/models/auth.model';
import { AuthCardComponent } from '../../../shared/components/organisms/auth-card/auth-card.component';
import { InputComponent } from '../../../shared/components/atoms/input/input.component';
import { ButtonComponent } from '../../../shared/components/atoms/button/button.component';
import { VALIDATION_RULES } from '../../../core/constants/validation.constants';
import { LinkComponent } from '../../../shared/components/atoms/link/link.component';
import { AlertComponent } from '../../../shared/components/atoms/alert/alert.component';

@Component({
  selector: 'app-login',
  imports: [
    FormField,
    FormRoot,
    RouterModule,
    TranslocoPipe,
    AuthCardComponent,
    InputComponent,
    ButtonComponent,
    LinkComponent,
    AlertComponent,
  ],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss',
})
export class LoginComponent {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  private readonly translate = inject(TranslocoService);

  protected readonly errorMessage = signal('');

  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  private readonly model = signal<LoginFormModel>({ email: '', password: '' });

  protected readonly loginForm = form(this.model, (path) => {
    required(path.email);
    email(path.email);
    maxLength(path.email, VALIDATION_RULES.EMAIL.MAX_LENGTH);

    required(path.password);
    maxLength(path.password, VALIDATION_RULES.PASSWORD.MAX_LENGTH);
  });

  onSubmit() {
    if (this.loginForm().invalid()) {
      return;
    }

    this.authService.login(this.loginForm().value()).subscribe({
      next: () => {
        this.router.navigate(['/recipes']);
      },
      error: () => {
        this.errorMessage.set(this.translate.translate('FEATURES.AUTH.LOGIN.FAILED'));
      },
    });
  }
}
