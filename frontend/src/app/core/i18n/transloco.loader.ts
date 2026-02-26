import { inject, Injectable, PLATFORM_ID } from '@angular/core';
import { isPlatformServer } from '@angular/common';
import { Translation, TranslocoLoader } from '@jsverse/transloco';
import { HttpClient } from '@angular/common/http';

@Injectable({ providedIn: 'root' })
export class TranslocoHttpLoader implements TranslocoLoader {
  private http = inject(HttpClient);
  private platformId = inject(PLATFORM_ID);

  getTranslation(lang: string) {
    const baseUrl = isPlatformServer(this.platformId) ? 'http://localhost:4200' : '';
    return this.http.get<Translation>(`${baseUrl}/assets/i18n/${lang}.json`);
  }
}
