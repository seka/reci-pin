import { Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  FieldTree,
  form,
  FormField,
  FormRoot,
  maxLength,
  minLength,
  required,
  validate,
} from '@angular/forms/signals';
import { Router, RouterModule } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';
import { InputComponent } from '../../../shared/components/atoms/input/input.component';
import { ButtonComponent } from '../../../shared/components/atoms/button/button.component';
import { HeadlineComponent } from '../../../shared/components/atoms/headline/headline.component';
import { VALIDATION_RULES } from '../../../core/constants/validation.constants';
import { TranslocoPipe, TranslocoService } from '@jsverse/transloco';
import { LinkComponent } from '../../../shared/components/atoms/link/link.component';
import { AlertComponent } from '../../../shared/components/atoms/alert/alert.component';

interface ChangePasswordFormValue {
  currentPassword: string;
  newPassword: string;
  confirmNewPassword: string;
}

@Component({
  selector: 'app-change-password',
  standalone: true,
  imports: [
    CommonModule,
    FormField,
    FormRoot,
    RouterModule,
    InputComponent,
    ButtonComponent,
    HeadlineComponent,
    TranslocoPipe,
    LinkComponent,
    AlertComponent,
  ],
  templateUrl: './change-password.component.html',
  styleUrl: './change-password.component.scss',
})
export class ChangePasswordComponent {
  private authService = inject(AuthService);
  private router = inject(Router);
  private translate = inject(TranslocoService);

  isProcessing = false;
  errorMessage = '';

  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  private readonly model = signal<ChangePasswordFormValue>({
    currentPassword: '',
    newPassword: '',
    confirmNewPassword: '',
  });

  protected readonly form = form(this.model, (path) => {
    required(path.currentPassword);
    maxLength(path.currentPassword, VALIDATION_RULES.PASSWORD.MAX_LENGTH);

    required(path.newPassword);
    minLength(path.newPassword, VALIDATION_RULES.PASSWORD.MIN_LENGTH);
    maxLength(path.newPassword, VALIDATION_RULES.PASSWORD.MAX_LENGTH);

    required(path.confirmNewPassword);
    maxLength(path.confirmNewPassword, VALIDATION_RULES.PASSWORD.MAX_LENGTH);
    validate(path.confirmNewPassword, (ctx) =>
      ctx.valueOf(path.newPassword) === ctx.value() ? null : { kind: 'mismatch' },
    );
  });

  getErrorMessage(field: FieldTree<string>): string | null {
    const state = field();
    if (!state.touched() || state.errors().length === 0) {
      return null;
    }
    for (const err of state.errors()) {
      switch (err.kind) {
        case 'required':
          return this.translate.translate('VALIDATION.REQUIRED');
        case 'minLength':
          return this.translate.translate('VALIDATION.MIN_LENGTH', {
            min: VALIDATION_RULES.PASSWORD.MIN_LENGTH,
          });
        case 'maxLength':
          return this.translate.translate('VALIDATION.MAX_LENGTH', {
            max: VALIDATION_RULES.PASSWORD.MAX_LENGTH,
          });
        case 'mismatch':
          return this.translate.translate('VALIDATION.PASSWORD_MISMATCH');
      }
    }
    return null;
  }

  onSubmit() {
    if (this.form().invalid()) return;

    this.isProcessing = true;
    this.errorMessage = '';

    this.authService.changePassword(this.form().value()).subscribe({
      next: () => {
        alert(this.translate.translate('FEATURES.SETTINGS.CHANGE_PASSWORD.SUCCESS'));
        this.router.navigate(['/settings']);
      },
      error: (err) => {
        this.isProcessing = false;
        // Backend returns bad request for incorrect password or validation errors
        if (err.status === 400 || err.status === 401) {
          this.errorMessage = this.translate.translate(
            'FEATURES.SETTINGS.CHANGE_PASSWORD.FAILED_INVALID',
          );
        } else {
          this.errorMessage = this.translate.translate(
            'FEATURES.SETTINGS.CHANGE_PASSWORD.FAILED_ERROR',
          );
        }
        console.error(err);
      },
    });
  }
}
