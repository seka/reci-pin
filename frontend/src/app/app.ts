import { Component, Inject, PLATFORM_ID, signal } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { isPlatformBrowser } from '@angular/common';
import { TranslateService } from '@ngx-translate/core';
import { HeaderComponent } from './shared/components/organisms/header/header.component';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, HeaderComponent],
  template: `
    <app-header></app-header>
    <router-outlet></router-outlet>
  `,
  styleUrl: './app.scss',
})
export class App {
  protected readonly title = signal('frontend');

  constructor(
    private translate: TranslateService,
    @Inject(PLATFORM_ID) private platformId: Object,
  ) {
    this.translate.addLangs(['ja', 'en']);
    this.translate.setDefaultLang('ja');

    if (isPlatformBrowser(this.platformId)) {
      const savedLang = localStorage.getItem('lang');
      const browserLang = this.translate.getBrowserLang();

      // Normalize language code (e.g. 'en-US' -> 'en')
      let langToUse = 'ja';
      if (savedLang) {
        langToUse = savedLang;
      } else if (browserLang?.match(/en/)) {
        langToUse = 'en';
      } else if (browserLang?.match(/ja/)) {
        langToUse = 'ja';
      }

      this.translate.use(langToUse);
    } else {
      this.translate.use('ja');
    }
  }
}
