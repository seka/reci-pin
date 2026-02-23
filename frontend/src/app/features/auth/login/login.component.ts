import { Component, inject, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { AuthService } from '../../../core/services/auth.service';
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
    TranslatePipe,
    AuthCardComponent,
    InputComponent,
    ButtonComponent,
    LinkComponent,
    AlertComponent,
  ],
  template: `
    <app-auth-card [title]="'AUTH.LOGIN_BUTTON' | translate">
      <form [formGroup]="loginForm" (ngSubmit)="onSubmit()">
        <app-input
          [label]="'AUTH.EMAIL' | translate"
          type="email"
          formControlName="email"
          [required]="true"
          [maxLength]="VALIDATION_RULES.EMAIL.MAX_LENGTH"
        ></app-input>

        <div style="margin-top: var(--spacing-2);">
          <app-input
            [label]="'AUTH.PASSWORD' | translate"
            type="password"
            formControlName="password"
            [required]="true"
            [maxLength]="VALIDATION_RULES.PASSWORD.MAX_LENGTH"
          ></app-input>
        </div>

        <div class="actions">
          <app-button variant="primary" type="submit" class="submit-btn" [disabled]="loginForm.invalid">{{ 'AUTH.LOGIN_BUTTON' | translate }}</app-button>
        </div>

        <div class="forgot-password-link">
          <app-link routerLink="/password-reset/request" variant="secondary">
            {{ 'AUTH.FORGOT_PASSWORD' | translate }}
          </app-link>
        </div>

        <app-alert type="error" [message]="errorMessage"></app-alert>
      </form>

      <div footer>
        <a routerLink="/signup" class="no-text-decoration">
          <app-button variant="accent" type="button" class="full-width-btn"
            >{{ 'AUTH.CREATE_ACCOUNT' | translate }}</app-button
          >
        </a>
      </div>
    </app-auth-card>
  `,
  styles: [
    `
      .actions {
        margin-top: var(--spacing-3);
        margin-bottom: var(--spacing-2);
      }
      .submit-btn {
        width: 100%;
      }
      .submit-btn:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }
      .full-width-btn {
        width: 100%;
      }
    `,
  ],
})
export class LoginComponent implements OnInit {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);
  private readonly translate = inject(TranslateService);

  loginForm!: FormGroup;
  errorMessage = '';

  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  ngOnInit() {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email, Validators.maxLength(VALIDATION_RULES.EMAIL.MAX_LENGTH)]],
      password: ['', [Validators.required, Validators.maxLength(VALIDATION_RULES.PASSWORD.MAX_LENGTH)]],
    });
  }

  onSubmit() {
    if (this.loginForm.invalid) {
      return;
    }

    this.authService.login(this.loginForm.value).subscribe({
      next: () => {
        this.router.navigate(['/recipes']);
      },
      error: () => {
        this.errorMessage = this.translate.instant('AUTH.LOGIN_FAILED');
      },
    });
  }
}
