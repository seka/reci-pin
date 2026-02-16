import { Component, inject } from '@angular/core';
import { AsyncPipe } from '@angular/common';
import { RouterModule } from '@angular/router';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { AuthService } from '../../../../core/services/auth.service';
import { LogoComponent } from '../../atoms/logo/logo.component';
import { ButtonComponent } from '../../atoms/button/button.component';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [AsyncPipe, RouterModule, TranslatePipe, LogoComponent, ButtonComponent],
  template: `
    <header class="app-header">
      <div class="header-content">
        <app-logo size="medium"></app-logo>
        <div class="actions">
          <div class="lang-switcher">
            <button class="lang-btn" [class.active]="currentLang === 'ja'" (click)="switchLang('ja')">JP</button>
            <span class="divider">/</span>
            <button class="lang-btn" [class.active]="currentLang === 'en'" (click)="switchLang('en')">EN</button>
          </div>
          @if (currentUser$ | async; as user) {
            <div class="user-actions">
              <span class="username">{{ 'HELLO' | translate }}, {{ user.name }}</span>
              <app-button variant="text" (click)="onLogout()">{{ 'LOGOUT' | translate }}</app-button>
            </div>
          }
        </div>
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
      .actions {
        display: flex;
        align-items: center;
        gap: 24px;
      }
      .lang-switcher {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 0.9rem;
        color: #888;
      }
      .lang-btn {
        background: none;
        border: none;
        cursor: pointer;
        padding: 4px;
        color: inherit;
        font-weight: 500;
        transition: color 0.2s;
        
        &.active, &:hover {
          color: #333;
          font-weight: 700;
        }
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
  private readonly translate = inject(TranslateService);

  currentUser$ = this.authService.currentUser$;

  get currentLang() {
    return this.translate.currentLang;
  }

  onLogout() {
    this.authService.logout();
  }

  switchLang(lang: string) {
    this.translate.use(lang);
    localStorage.setItem('lang', lang);
  }
}
