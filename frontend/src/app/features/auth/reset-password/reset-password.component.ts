import { Component, inject, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { TranslocoPipe, TranslocoService } from '@jsverse/transloco';
import { AuthService } from '../../../core/services/auth.service';
import { AuthCardComponent } from '../../../shared/components/organisms/auth-card/auth-card.component';
import { InputComponent } from '../../../shared/components/atoms/input/input.component';
import { ButtonComponent } from '../../../shared/components/atoms/button/button.component';
import { VALIDATION_RULES } from '../../../core/constants/validation.constants';
import { AlertComponent } from '../../../shared/components/atoms/alert/alert.component';

@Component({
  selector: 'app-reset-password',
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
  templateUrl: './reset-password.component.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrls: ['./reset-password.component.scss'],
})
export class ResetPasswordComponent implements OnInit {
  private readonly authService = inject(AuthService);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly translate = inject(TranslocoService);
  private readonly fb = inject(FormBuilder).nonNullable;

  form: FormGroup = this.fb.group({
    token: ['', [Validators.required]],
    newPassword: [
      '',
      [
        Validators.required,
        Validators.minLength(VALIDATION_RULES.PASSWORD.MIN_LENGTH),
        Validators.maxLength(VALIDATION_RULES.PASSWORD.MAX_LENGTH),
      ],
    ],
  });

  token = '';
  message = '';
  errorMessage = '';
  isLoading = false;
  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  ngOnInit() {
    this.route.queryParams.subscribe((params) => {
      this.token = params['token'] || '';
      if (this.token) {
        this.form.patchValue({ token: this.token });
      } else {
        this.errorMessage = this.translate.translate('FEATURES.AUTH.RESET_PASSWORD.INVALID_LINK');
      }
    });
  }

  onSubmit() {
    if (this.form.invalid) return;

    this.isLoading = true;
    this.message = '';
    this.errorMessage = '';

    this.authService.resetPassword(this.form.getRawValue()).subscribe({
      next: (res) => {
        this.message = res.message;
        this.isLoading = false;
        setTimeout(() => {
          this.router.navigate(['/login']);
        }, 3000);
      },
      error: (err) => {
        this.errorMessage = this.translate.translate('FEATURES.AUTH.RESET_PASSWORD.FAILED_EXPIRED');
        this.isLoading = false;
        console.error(err);
      },
    });
  }
}
