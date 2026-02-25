import type { Preview } from '@storybook/angular'
import { setCompodocJson } from "@storybook/addon-docs/angular";
import { applicationConfig } from "@storybook/angular";
import docJson from "../documentation.json";
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
  },
};

export default preview;