import { Component, signal, inject, ChangeDetectionStrategy } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { TranslocoService } from '@jsverse/transloco';
import { HeaderComponent } from './shared/components/organisms/header/header.component';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, HeaderComponent],
  templateUrl: './app.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './app.scss',
})
export class App {
  protected readonly title = signal('frontend');

  private translate = inject(TranslocoService);

  constructor() {
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
  }
}
