import { Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';
import { AuthCardComponent } from '../../../shared/components/organisms/auth-card/auth-card.component';
import { InputComponent } from '../../../shared/components/atoms/input/input.component';
import { ButtonComponent } from '../../../shared/components/atoms/button/button.component';

@Component({
    selector: 'app-reset-password',
    standalone: true,
    imports: [
        CommonModule,
        FormsModule,
        RouterModule,
        AuthCardComponent,
        InputComponent,
        ButtonComponent,
    ],
    templateUrl: './reset-password.component.html',
    styleUrls: ['./reset-password.component.scss'],
})
export class ResetPasswordComponent implements OnInit {
    private readonly authService = inject(AuthService);
    private readonly route = inject(ActivatedRoute);
    private readonly router = inject(Router);

    token = '';
    newPassword = '';
    message = '';
    errorMessage = '';
    isLoading = false;

    ngOnInit() {
        this.route.queryParams.subscribe((params) => {
            this.token = params['token'] || '';
            if (!this.token) {
                this.errorMessage = '無効なリンクです。';
            }
        });
    }

    onSubmit() {
        if (!this.token) {
            this.errorMessage = 'トークンが不足しています。';
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
                this.errorMessage = 'パスワードの再設定に失敗しました。リンクの有効期限が切れている可能性があります。';
                this.isLoading = false;
                console.error(err);
            },
        });
    }
}
