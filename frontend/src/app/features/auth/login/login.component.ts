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
import { LoginFormModel } from '../../../core/models/auth.model';
import { AuthCardComponent } from '../../../shared/components/organisms/auth-card/auth-card.component';
import { InputComponent } from '../../../shared/components/atoms/input/input.component';
import { ButtonComponent } from '../../../shared/components/atoms/button/button.component';
import { VALIDATION_RULES } from '../../../core/constants/validation.constants';
import { LinkComponent } from '../../../shared/components/atoms/link/link.component';
import { AlertComponent } from '../../../shared/components/atoms/alert/alert.component';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [
    FormsModule,
    RouterModule,
    ReactiveFormsModule,
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
export class LoginComponent implements OnInit {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);
  private readonly translate = inject(TranslocoService);

  loginForm!: FormGroup;
  errorMessage = '';

  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  ngOnInit() {
    this.loginForm = this.fb.group({
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
        [Validators.required, Validators.maxLength(VALIDATION_RULES.PASSWORD.MAX_LENGTH)],
      ],
    });
  }

  onSubmit() {
    if (this.loginForm.invalid) {
      return;
    }

    this.authService.login(this.loginForm.value as LoginFormModel).subscribe({
      next: () => {
        this.router.navigate(['/recipes']);
      },
      error: () => {
        this.errorMessage = this.translate.translate('FEATURES.AUTH.LOGIN.FAILED');
      },
    });
  }
}
