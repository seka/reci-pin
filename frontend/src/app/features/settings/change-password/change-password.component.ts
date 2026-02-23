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
  ],
  template: `
    <div class="change-password-container">
      <div class="back-link">
        <app-link routerLink="/settings">{{ 'SETTINGS.BACK_TO_SETTINGS' | translate }}</app-link>
      </div>

      <app-headline level="1">{{ 'SETTINGS.CHANGE_PASSWORD_TITLE' | translate }}</app-headline>

      <form [formGroup]="form" (ngSubmit)="onSubmit()" class="password-form">
        <app-input
          [label]="'SETTINGS.CURRENT_PASSWORD' | translate"
          type="password"
          formControlName="currentPassword"
          [errorMessage]="getErrorMessage('currentPassword')"
          [required]="true"
          [maxLength]="VALIDATION_RULES.PASSWORD.MAX_LENGTH"
        ></app-input>
        <p class="forgot-password-link">
          <app-link routerLink="/password-reset/request">{{ 'AUTH.FORGOT_PASSWORD' | translate }}</app-link>
        </p>

        <app-input
          [label]="'SETTINGS.NEW_PASSWORD' | translate"
          type="password"
          formControlName="newPassword"
          [errorMessage]="getErrorMessage('newPassword')"
          [required]="true"
          [maxLength]="VALIDATION_RULES.PASSWORD.MAX_LENGTH"
          [showCounter]="true"
        ></app-input>
        <p class="hint">{{ 'SETTINGS.PASSWORD_HINT' | translate: { min: VALIDATION_RULES.PASSWORD.MIN_LENGTH } }}</p>

        <app-input
          [label]="'SETTINGS.NEW_PASSWORD_CONFIRM' | translate"
          type="password"
          formControlName="confirmNewPassword"
          [errorMessage]="getErrorMessage('confirmNewPassword')"
          [required]="true"
          [maxLength]="VALIDATION_RULES.PASSWORD.MAX_LENGTH"
        ></app-input>

        <div class="actions">
          <app-button type="submit" [disabled]="form.invalid || isProcessing">
            {{ isProcessing ? ('SETTINGS.CHANGING_PASSWORD' | translate) : ('SETTINGS.CHANGE_PASSWORD_BUTTON' | translate) }}
          </app-button>
        </div>

        @if (errorMessage) {
          <p class="error-global">{{ errorMessage }}</p>
        }
      </form>
    </div>
  `,
  styles: [
    `
      .change-password-container {
        max-width: 600px;
        margin: 0 auto;
        padding: var(--spacing-3);
      }

      .back-link {
        margin-bottom: var(--spacing-2);
      }

      .password-form {
        margin-top: var(--spacing-2);
        display: flex;
        flex-direction: column;
        gap: var(--spacing-1_5);
        background: var(--color-surface);
        padding: var(--spacing-3);
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      }

      .hint {
        font-size: var(--font-size-2);
        color: var(--color-text-secondary);
        margin-top: -8px;
        margin-left: 4px;
      }

      .forgot-password-link {
        font-size: var(--font-size-2);
        margin-top: -8px;
        margin-left: 4px;
        text-align: left;
      }

      .actions {
        margin-top: var(--spacing-1);
      }

      .error-global {
        color: var(--color-error);
        margin-top: var(--spacing-1_5);
        font-weight: bold;
      }
    `,
  ],
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
      currentPassword: ['', [Validators.required, Validators.maxLength(VALIDATION_RULES.PASSWORD.MAX_LENGTH)]],
      newPassword: [
        '',
        [
          Validators.required,
          Validators.minLength(VALIDATION_RULES.PASSWORD.MIN_LENGTH),
          Validators.maxLength(VALIDATION_RULES.PASSWORD.MAX_LENGTH),
        ],
      ],
      confirmNewPassword: ['', [Validators.required, Validators.maxLength(VALIDATION_RULES.PASSWORD.MAX_LENGTH)]],
    },
    { validators: this.passwordMatchValidator }
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
        return this.translate.instant('VALIDATION.MIN_LENGTH', { min: VALIDATION_RULES.PASSWORD.MIN_LENGTH });
      if (control.errors['maxlength'])
        return this.translate.instant('VALIDATION.MAX_LENGTH', { max: VALIDATION_RULES.PASSWORD.MAX_LENGTH });
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
