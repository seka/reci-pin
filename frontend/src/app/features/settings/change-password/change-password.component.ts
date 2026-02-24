import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  FormBuilder,
  ReactiveFormsModule,
  Validators,
  AbstractControl,
  ValidationErrors,
} from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';
import { InputComponent } from '../../../shared/components/atoms/input/input.component';
import { ButtonComponent } from '../../../shared/components/atoms/button/button.component';
import { HeadlineComponent } from '../../../shared/components/atoms/headline/headline.component';
import { VALIDATION_RULES } from '../../../core/constants/validation.constants';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { LinkComponent } from '../../../shared/components/atoms/link/link.component';
import { AlertComponent } from '../../../shared/components/atoms/alert/alert.component';

@Component({
  selector: 'app-change-password',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    InputComponent,
    ButtonComponent,
    HeadlineComponent,
    TranslatePipe,
    LinkComponent,
    AlertComponent,
  ],
  templateUrl: './change-password.component.html',
  styleUrl: './change-password.component.scss',
})
export class ChangePasswordComponent {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);
  private translate = inject(TranslateService);

  isProcessing = false;
  errorMessage = '';

  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  form = this.fb.group(
    {
      currentPassword: [
        '',
        [Validators.required, Validators.maxLength(VALIDATION_RULES.PASSWORD.MAX_LENGTH)],
      ],
      newPassword: [
        '',
        [
          Validators.required,
          Validators.minLength(VALIDATION_RULES.PASSWORD.MIN_LENGTH),
          Validators.maxLength(VALIDATION_RULES.PASSWORD.MAX_LENGTH),
        ],
      ],
      confirmNewPassword: [
        '',
        [Validators.required, Validators.maxLength(VALIDATION_RULES.PASSWORD.MAX_LENGTH)],
      ],
    },
    { validators: this.passwordMatchValidator },
  );

  private passwordMatchValidator(control: AbstractControl): ValidationErrors | null {
    const newPassword = control.get('newPassword')?.value;
    const confirmPassword = control.get('confirmNewPassword')?.value;
    return newPassword === confirmPassword ? null : { mismatch: true };
  }

  getErrorMessage(controlName: string): string | null {
    const control = this.form.get(controlName);
    if (control?.touched && control?.errors) {
      if (control.errors['required']) return this.translate.instant('VALIDATION.REQUIRED');
      if (control.errors['minlength'])
        return this.translate.instant('VALIDATION.MIN_LENGTH', {
          min: VALIDATION_RULES.PASSWORD.MIN_LENGTH,
        });
      if (control.errors['maxlength'])
        return this.translate.instant('VALIDATION.MAX_LENGTH', {
          max: VALIDATION_RULES.PASSWORD.MAX_LENGTH,
        });
      if (controlName === 'confirmNewPassword' && this.form.errors?.['mismatch']) {
        return this.translate.instant('VALIDATION.PASSWORD_MISMATCH');
      }
    }
    return null;
  }

  onSubmit() {
    if (this.form.invalid) return;

    this.isProcessing = true;
    this.errorMessage = '';

    const { currentPassword, newPassword } = this.form.value;

    if (!currentPassword || !newPassword) return;

    this.authService
      .changePassword({
        currentPassword: currentPassword,
        newPassword: newPassword,
      })
      .subscribe({
        next: () => {
          alert(this.translate.instant('SETTINGS.PASSWORD_CHANGED'));
          this.router.navigate(['/settings']);
        },
        error: (err) => {
          this.isProcessing = false;
          // Backend returns bad request for incorrect password or validation errors
          if (err.status === 400 || err.status === 401) {
            this.errorMessage = this.translate.instant('SETTINGS.CHANGE_FAILED_INVALID');
          } else {
            this.errorMessage = this.translate.instant('SETTINGS.CHANGE_FAILED_ERROR');
          }
          console.error(err);
        },
      });
  }
}
