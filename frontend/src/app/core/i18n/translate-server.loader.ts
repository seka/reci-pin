import { TranslateLoader } from '@ngx-translate/core';
import { Observable } from 'rxjs';
import { join } from 'path';
import { readFileSync, existsSync } from 'fs';

export class TranslateServerLoader implements TranslateLoader {
    getTranslation(lang: string): Observable<any> {
        return new Observable((observer) => {
            const assetPath = process.env['I18N_ASSET_PATH'];

            if (!assetPath) {
                console.error('I18N_ASSET_PATH environment variable is not set.');
                observer.next({});
                observer.complete();
                return;
            }

            const path = join(assetPath, `${lang}.json`);

            if (existsSync(path)) {
                try {
                    const jsonData = JSON.parse(readFileSync(path, 'utf8'));
                    observer.next(jsonData);
                    observer.complete();
                    return;
                } catch (e) {
                    console.error(`Failed to read/parse translation file: ${path}`, e);
                }
            } else {
                console.warn(`Translation file not found at: ${path} (lang: ${lang})`);
            }

            observer.next({});
            observer.complete();
        });
    }
}
