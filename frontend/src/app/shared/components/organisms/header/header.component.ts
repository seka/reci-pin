import { Component, inject } from '@angular/core';
import { AsyncPipe } from '@angular/common';
import { RouterModule } from '@angular/router';
import { TranslocoPipe, TranslocoService } from '@jsverse/transloco';
import { AuthService } from '../../../../core/services/auth.service';
import { LogoComponent } from '../../atoms/logo/logo.component';
import { ButtonComponent } from '../../atoms/button/button.component';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [AsyncPipe, RouterModule, TranslocoPipe, LogoComponent, ButtonComponent],
  templateUrl: './header.component.html',
  styleUrl: './header.component.scss',
})
export class HeaderComponent {
  private readonly authService = inject(AuthService);
  private readonly translate = inject(TranslocoService);

  currentUser$ = this.authService.currentUser$;

  get currentLang() {
    return this.translate.getActiveLang();
  }

  onLogout() {
    this.authService.logout();
  }

  switchLang(lang: string) {
    this.translate.setActiveLang(lang);
    localStorage.setItem('lang', lang);
  }
}
