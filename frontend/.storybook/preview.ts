import 'zone.js';
import '../src/styles.scss';
import type { Preview } from '@storybook/angular-vite';
import { setCompodocJson } from '@storybook/addon-docs/angular';
import { applicationConfig } from '@storybook/angular-vite';
import docJson from '../documentation.json';
import { provideHttpClient, withFetch } from '@angular/common/http';
import { provideTransloco } from '@jsverse/transloco';
import { TranslocoHttpLoader } from '../src/app/core/i18n/transloco.loader';
import { isDevMode } from '@angular/core';

setCompodocJson(docJson);

const preview: Preview = {
  decorators: [
    applicationConfig({
      providers: [
        provideHttpClient(withFetch()),
        provideTransloco({
          config: {
            availableLangs: ['ja', 'en'],
            defaultLang: 'ja',
            reRenderOnLangChange: true,
            prodMode: !isDevMode(),
          },
          loader: TranslocoHttpLoader,
        }),
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

    a11y: {
      // 'todo' - show a11y violations in the test UI only
      // 'error' - fail CI on a11y violations
      // 'off' - skip a11y checks entirely
      test: 'todo',
    },
  },
};

export default preview;
