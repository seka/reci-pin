/**
 * @vitest-environment jsdom
 */
import { expect, describe, it, beforeEach, afterEach } from 'vitest';
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
        provideHttpClientTesting(),
      ],
    });

    httpTestingController = TestBed.inject(HttpTestingController);
    httpClient = TestBed.inject(HttpClient);
  });

  afterEach(() => {
    httpTestingController.verify();
  });

  it('should convert request body keys from camelCase to snake_case', () => {
    const testData = { firstName: 'John', lastName: 'Doe' };
    httpClient.post('/api/test', testData).subscribe();

    const req = httpTestingController.expectOne('/api/test');
    expect(req.request.body).toEqual({ first_name: 'John', last_name: 'Doe' });
    req.flush({});
  });

  it('should not convert request body if it is FormData', () => {
    const formData = new FormData();
    formData.append('firstName', 'John');
    httpClient.post('/api/test', formData).subscribe();

    const req = httpTestingController.expectOne('/api/test');
    expect(req.request.body).toBeInstanceOf(FormData);
    req.flush({});
  });

  it('should not convert request body if it is a Blob', () => {
    const blob = new Blob(['test'], { type: 'text/plain' });
    httpClient.post('/api/test', blob).subscribe();

    const req = httpTestingController.expectOne('/api/test');
    expect(req.request.body).toBeInstanceOf(Blob);
    req.flush({});
  });

  it('should convert response body keys from snake_case to camelCase', () => {
    const mockResponse = { first_name: 'John', last_name: 'Doe' };
    httpClient.get('/api/test').subscribe((data) => {
      const typedData = data as { firstName: string; lastName: string };
      expect(typedData).toEqual({ firstName: 'John', lastName: 'Doe' });
    });

    const req = httpTestingController.expectOne('/api/test');
    req.flush(mockResponse);
  });

  it('should not convert response body if it is a Blob', () => {
    const blob = new Blob(['test'], { type: 'application/pdf' });
    httpClient.get('/api/test', { responseType: 'blob' }).subscribe((data) => {
      expect(data).toBeInstanceOf(Blob);
    });

    const req = httpTestingController.expectOne('/api/test');
    req.flush(blob);
  });

  it('should skip conversion for non-API requests', () => {
    const testData = { firstName: 'John' };
    httpClient.post('/other-api/test', testData).subscribe();

    const req = httpTestingController.expectOne('/other-api/test');
    expect(req.request.body).toEqual({ firstName: 'John' });
    req.flush({});
  });
});
