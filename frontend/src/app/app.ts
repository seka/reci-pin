import { Component, PLATFORM_ID, signal, inject } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { isPlatformBrowser } from '@angular/common';
import { TranslocoService } from '@jsverse/transloco';
import { HeaderComponent } from './shared/components/organisms/header/header.component';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, HeaderComponent],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class App {
  protected readonly title = signal('frontend');

  private translate = inject(TranslocoService);
  private platformId = inject(PLATFORM_ID);

  constructor() {
    if (isPlatformBrowser(this.platformId as object)) {
      const savedLang = localStorage.getItem('lang');
      const browserLang = typeof window !== 'undefined' ? window.navigator.language : 'ja';

      // Normalize language code (e.g. 'en-US' -> 'en')
      let langToUse = 'ja';
      if (savedLang) {
        langToUse = savedLang;
      } else if (browserLang?.match(/en/)) {
        langToUse = 'en';
      } else if (browserLang?.match(/ja/)) {
        langToUse = 'ja';
      }

      this.translate.setActiveLang(langToUse);
    } else {
      this.translate.setActiveLang('ja');
    }
  }
}
