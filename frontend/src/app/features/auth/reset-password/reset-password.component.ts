import { Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
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
    FormsModule,
    RouterModule,
    TranslatePipe,
    AuthCardComponent,
    InputComponent,
    ButtonComponent,
    AlertComponent,
  ],
  templateUrl: './reset-password.component.html',
  styleUrls: ['./reset-password.component.scss'],
})
export class ResetPasswordComponent implements OnInit {
  private readonly authService = inject(AuthService);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly translate = inject(TranslateService);

  token = '';
  newPassword = '';
  message = '';
  errorMessage = '';
  isLoading = false;
  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  ngOnInit() {
    this.route.queryParams.subscribe((params) => {
      this.token = params['token'] || '';
      if (!this.token) {
        this.errorMessage = this.translate.instant('FEATURES.AUTH.RESET_PASSWORD.INVALID_LINK');
      }
    });
  }

  onSubmit() {
    if (!this.token) {
      this.errorMessage = this.translate.instant('FEATURES.AUTH.RESET_PASSWORD.MISSING_TOKEN');
      return;
    }

    this.isLoading = true;
    this.message = '';
    this.errorMessage = '';

    this.authService.resetPassword(this.token, this.newPassword).subscribe({
      next: (res) => {
        this.message = res.message;
        this.isLoading = false;
        setTimeout(() => {
          this.router.navigate(['/login']);
        }, 3000);
      },
      error: (err) => {
        this.errorMessage = this.translate.instant('FEATURES.AUTH.RESET_PASSWORD.FAILED_EXPIRED');
        this.isLoading = false;
        console.error(err);
      },
    });
  }
}
