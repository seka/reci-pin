import { Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
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
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
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
  private readonly fb = inject(FormBuilder).nonNullable;

  form: FormGroup = this.fb.group({
    email: [
      '',
      [
        Validators.required,
        Validators.email,
        Validators.maxLength(VALIDATION_RULES.EMAIL.MAX_LENGTH),
      ],
    ],
  });

  message = '';
  errorMessage = '';
  isLoading = false;
  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  onSubmit() {
    if (this.form.invalid) return;

    this.isLoading = true;
    this.message = '';
    this.errorMessage = '';

    this.authService.requestPasswordReset(this.form.getRawValue()).subscribe({
      next: (res) => {
        this.message = res.message;
        this.isLoading = false;
      },
      error: (err) => {
        this.errorMessage = this.translate.translate('FEATURES.AUTH.REQUEST_PASSWORD_RESET.FAILED');
        this.isLoading = false;
        console.error(err);
      },
    });
  }
}
