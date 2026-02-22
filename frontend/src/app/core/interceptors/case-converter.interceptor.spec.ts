import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { HttpClient, provideHttpClient, withInterceptors } from '@angular/common/http';
import { caseConverterInterceptor } from './case-converter.interceptor';

describe('caseConverterInterceptor', () => {
    let httpTestingController: HttpTestingController;
    let httpClient: HttpClient;

    beforeEach(() => {
        TestBed.configureTestingModule({
            providers: [
                provideHttpClient(withInterceptors([caseConverterInterceptor])),
                provideHttpClientTesting()
            ]
        });

        httpTestingController = TestBed.inject(HttpTestingController);
        httpClient = TestBed.inject(HttpClient);
    });

    afterEach(() => {
        httpTestingController.verify();
    });

    it('should convert request body keys from camelCase to snake_case', () => {
        const payload = {
            firstName: 'John',
            lastName: 'Doe',
            addressInfo: {
                zipCode: '12345'
            }
        };

        httpClient.post('/api/test', payload).subscribe();

        const req = httpTestingController.expectOne('/api/test');
        expect(req.request.method).toEqual('POST');
        expect(req.request.body).toEqual({
            first_name: 'John',
            last_name: 'Doe',
            address_info: {
                zip_code: '12345'
            }
        });

        // Resolve request
        req.flush({});
    });

    it('should not convert request body if it is FormData', () => {
        const formData = new FormData();
        formData.append('firstName', 'John');

        httpClient.post('/api/test-formdata', formData).subscribe();

        const req = httpTestingController.expectOne('/api/test-formdata');
        expect(req.request.body instanceof FormData).toBe(true);
        expect((req.request.body as FormData).get('firstName')).toBe('John');

        req.flush({});
    });

    it('should not convert request body if it is a Blob', () => {
        const blob = new Blob(['test content'], { type: 'text/plain' });

        httpClient.post('/api/test-blob', blob).subscribe();

        const req = httpTestingController.expectOne('/api/test-blob');
        expect(req.request.body instanceof Blob).toBe(true);

        req.flush({});
    });

    it('should convert response body keys from snake_case to camelCase', () => {
        const mockResponse = {
            user_id: 1,
            created_at: '2023-01-01',
            nested_data: {
                inner_property: true
            }
        };

        httpClient.get('/api/test-response').subscribe((res: any) => {
            expect(res).toEqual({
                userId: 1,
                createdAt: '2023-01-01',
                nestedData: {
                    innerProperty: true
                }
            });
        });

        const req = httpTestingController.expectOne('/api/test-response');
        req.flush(mockResponse);
    });

    it('should not convert response body if it is a Blob', () => {
        const blobResponse = new Blob(['response content'], { type: 'text/plain' });

        httpClient.get('/api/test-blob-response', { responseType: 'blob' }).subscribe((res: any) => {
            expect(res instanceof Blob).toBe(true);
        });

        const req = httpTestingController.expectOne('/api/test-blob-response');
        req.flush(blobResponse);
    });
});
