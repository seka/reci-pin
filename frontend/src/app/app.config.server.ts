import { mergeApplicationConfig, ApplicationConfig } from '@angular/core';
import { provideServerRendering, withRoutes } from '@angular/ssr';
import { appConfig } from './app.config';
import { serverRoutes } from './app.routes.server';

import { TranslateLoader } from '@ngx-translate/core';
import { TranslateServerLoader } from './core/i18n/translate-server.loader';

const serverConfig: ApplicationConfig = {
  providers: [
    provideServerRendering(withRoutes(serverRoutes)),
    {
      provide: TranslateLoader,
      useClass: TranslateServerLoader
    }
  ]
};

export const config = mergeApplicationConfig(appConfig, serverConfig);
