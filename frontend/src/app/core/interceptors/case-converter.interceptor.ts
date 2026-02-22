import { HttpInterceptorFn, HttpRequest, HttpResponse } from '@angular/common/http';
import { map } from 'rxjs/operators';
import camelcaseKeys from 'camelcase-keys';
import decamelizeKeys from 'decamelize-keys';

export const caseConverterInterceptor: HttpInterceptorFn = (req, next) => {
    let modifiedReq: HttpRequest<any> = req;

    // Convert request body to snake_case if it exists and is not FormData/Blob
    if (req.body && !(req.body instanceof FormData) && !(req.body instanceof Blob)) {
        try {
            modifiedReq = req.clone({
                body: decamelizeKeys(req.body, { deep: true }),
            });
        } catch (e) {
            console.error('Error converting request body keys to snake_case', e);
        }
    }

    return next(modifiedReq).pipe(
        map((event) => {
            if (event instanceof HttpResponse && event.body && typeof event.body === 'object' && !(event.body instanceof Blob)) {
                try {
                    // Convert response body to camelCase
                    const camelCaseBody = camelcaseKeys(event.body, { deep: true });
                    return event.clone({ body: camelCaseBody });
                } catch (e) {
                    console.error('Error converting response body keys to camelCase', e);
                }
            }
            return event;
        }),
    );
};
