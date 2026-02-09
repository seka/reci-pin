import { Component, inject } from '@angular/core';
import { AsyncPipe } from '@angular/common';
import { RouterModule } from '@angular/router';
import { AuthService } from '../../../../core/services/auth.service';
import { LogoComponent } from '../../atoms/logo/logo.component';
import { ButtonComponent } from '../../atoms/button/button.component';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [AsyncPipe, RouterModule, LogoComponent, ButtonComponent],
  template: `
    <header class="app-header">
      <div class="header-content">
        <app-logo size="medium"></app-logo>
        @if (currentUser$ | async; as user) {
          <div class="user-actions">
            <span class="username">Hello, {{ user.name }}</span>
            <app-button variant="text" (click)="onLogout()">ログアウト</app-button>
          </div>
        }
      </div>
    </header>
  `,
  styles: [
    `
      .app-header {
        background: #fff;
        border-bottom: 1px solid #eee;
        padding: 0 24px;
        height: 64px;
        display: flex;
        align-items: center;
        position: sticky;
        top: 0;
        z-index: 1000;
      }
      .header-content {
        width: 100%;
        max-width: 1200px;
        margin: 0 auto;
        display: flex;
        justify-content: space-between;
        align-items: center;
      }
      .user-actions {
        display: flex;
        align-items: center;
        gap: 16px;
      }
      .username {
        color: #555;
        font-weight: 500;
        font-size: 1rem;
      }
    `,
  ],
})
export class HeaderComponent {
  private readonly authService = inject(AuthService);

  currentUser$ = this.authService.currentUser$;

  onLogout() {
    this.authService.logout();
  }
}
