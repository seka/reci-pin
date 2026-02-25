export interface ValidationErrorDetail {
    code: string;
    params?: Record<string, string | number>;
}

export interface ApiError {
    error: {
        message: string;
        code?: string;
        details?: Record<string, ValidationErrorDetail[]>;
    };
}
