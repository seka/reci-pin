import { Component, inject, OnInit } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { TranslocoPipe, TranslocoService } from '@jsverse/transloco';
import { AuthService } from '../../../core/services/auth.service';
import { AuthCardComponent } from '../../../shared/components/organisms/auth-card/auth-card.component';
import { InputComponent } from '../../../shared/components/atoms/input/input.component';
import { ButtonComponent } from '../../../shared/components/atoms/button/button.component';
import { VALIDATION_RULES } from '../../../core/constants/validation.constants';
import { AlertComponent } from '../../../shared/components/atoms/alert/alert.component';
import { ApiError } from '../../../core/models/api-error.model';

@Component({
  selector: 'app-signup',
  standalone: true,
  imports: [
    FormsModule,
    RouterModule,
    ReactiveFormsModule,
    TranslocoPipe,
    AuthCardComponent,
    InputComponent,
    ButtonComponent,
    AlertComponent,
  ],
  templateUrl: './signup.component.html',
  styleUrl: './signup.component.scss',
})
export class SignupComponent implements OnInit {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);
  private readonly translate = inject(TranslocoService);

  signupForm!: FormGroup;
  fieldErrors: Record<string, string[]> = {};
  errorMessage = '';

  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  ngOnInit() {
    this.signupForm = this.fb.group({
      name: [
        '',
        [Validators.required, Validators.maxLength(VALIDATION_RULES.RECIPE.NAME_MAX_LENGTH)],
      ],
      email: [
        '',
        [
          Validators.required,
          Validators.email,
          Validators.maxLength(VALIDATION_RULES.EMAIL.MAX_LENGTH),
        ],
      ],
      password: [
        '',
        [
          Validators.required,
          Validators.minLength(VALIDATION_RULES.PASSWORD.MIN_LENGTH),
          Validators.maxLength(VALIDATION_RULES.PASSWORD.MAX_LENGTH),
          Validators.pattern(/[a-zA-Z]/),
          Validators.pattern(/[0-9]/),
        ],
      ],
    });
  }

  onSubmit() {
    if (this.signupForm.invalid) {
      return;
    }

    this.fieldErrors = {}; // Reset errors
    this.errorMessage = '';

    this.authService.signup(this.signupForm.value).subscribe({
      next: () => {
        this.router.navigate(['/recipes']);
      },
      error: (err: { error?: ApiError }) => {
        // バックエンドからのエラーレスポンスを処理
        if (err.error?.error?.details) {
          const details = err.error.error.details;

          // 各フィールドのエラーをマッピング
          Object.keys(details).forEach((field) => {
            const fieldDetails = details[field];
            const messages = fieldDetails.map((d) => {
              switch (d.code) {
                case 'PASSWORD_TOO_SHORT':
                  return this.translate.translate('VALIDATION.MIN_LENGTH', {
                    min: d.params?.['min'] || 8,
                  });
                case 'PASSWORD_TOO_LONG':
                  return this.translate.translate('VALIDATION.MAX_LENGTH', {
                    max: d.params?.['max'] || VALIDATION_RULES.PASSWORD.MAX_LENGTH,
                  });
                case 'PASSWORD_NO_ALPHA':
                  return this.translate.translate('VALIDATION.PASSWORD_NO_ALPHA');
                case 'PASSWORD_NO_NUMERIC':
                  return this.translate.translate('VALIDATION.PASSWORD_NO_NUMERIC');
                case 'EMAIL_INVALID_FORMAT':
                  return this.translate.translate('VALIDATION.INVALID_EMAIL');
                case 'EMAIL_TOO_LONG':
                  return this.translate.translate('VALIDATION.MAX_LENGTH', {
                    max: d.params?.['max'] || VALIDATION_RULES.EMAIL.MAX_LENGTH,
                  });
                case 'REQUIRED':
                  return this.translate.translate('VALIDATION.REQUIRED');
                default:
                  return this.translate.translate('VALIDATION.INVALID_INPUT');
              }
            });
            this.fieldErrors[field] = messages;
          });

          if (Object.keys(this.fieldErrors).length > 0) {
            return;
          }
        }
        this.errorMessage = this.translate.translate('FEATURES.AUTH.SIGNUP.FAILED');
      },
    });
  }
}
