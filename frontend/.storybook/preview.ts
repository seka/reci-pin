import type { Preview } from '@storybook/angular'
import { setCompodocJson } from "@storybook/addon-docs/angular";
import { applicationConfig } from "@storybook/angular";
import docJson from "../documentation.json";
import { importProvidersFrom, APP_INITIALIZER } from '@angular/core';
import { provideHttpClient, withFetch } from '@angular/common/http';
import { TranslateModule, TranslateLoader, TranslateService } from '@ngx-translate/core';
import { TranslateHttpLoader, TRANSLATE_HTTP_LOADER_CONFIG } from '@ngx-translate/http-loader';
import { HttpClient } from '@angular/common/http';

export function HttpLoaderFactory() {
  return new TranslateHttpLoader();
}

export function setupTranslateFactory(service: TranslateService): () => void {
  return () => {
    service.setDefaultLang('ja');
    service.use('ja');
  };
}

setCompodocJson(docJson);

const preview: Preview = {
  decorators: [
    applicationConfig({
      providers: [
        provideHttpClient(withFetch()),
        {
          provide: TRANSLATE_HTTP_LOADER_CONFIG,
          useValue: {
            prefix: '/assets/i18n/',
            suffix: '.json',
          },
        },
        importProvidersFrom(
          TranslateModule.forRoot({
            loader: {
              provide: TranslateLoader,
              useFactory: HttpLoaderFactory,
              deps: [HttpClient],
            },
            defaultLanguage: 'ja',
          })
        ),
        {
          provide: APP_INITIALIZER,
          useFactory: setupTranslateFactory,
          deps: [TranslateService],
          multi: true,
        },
      ],
    }),
  ],
  parameters: {
    controls: {
      matchers: {
       color: /(background|color)$/i,
       date: /Date$/i,
      },
    },
  },
};

export default preview;