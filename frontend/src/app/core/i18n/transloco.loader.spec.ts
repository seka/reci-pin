/**
 * @vitest-environment jsdom
 */
import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting, HttpTestingController } from '@angular/common/http/testing';
import { TranslocoHttpLoader } from './transloco.loader';
import { describe, it, expect, beforeEach, afterEach } from 'vitest';

describe('TranslocoHttpLoader', () => {
    let loader: TranslocoHttpLoader;
    let httpMock: HttpTestingController;

    beforeEach(() => {
        TestBed.configureTestingModule({
            providers: [
                TranslocoHttpLoader,
                provideHttpClient(),
                provideHttpClientTesting()
            ]
        });
        loader = TestBed.inject(TranslocoHttpLoader);
        httpMock = TestBed.inject(HttpTestingController);
    });

    afterEach(() => {
        httpMock.verify();
    });

    it('should get translation from assets', () => {
        const mockTranslation = { test: 'value' };
        loader.getTranslation('en').subscribe((translation) => {
            expect(translation).toEqual(mockTranslation);
        });

        const req = httpMock.expectOne('/assets/i18n/en.json');
        expect(req.request.method).toBe('GET');
        req.flush(mockTranslation);
    });
});
