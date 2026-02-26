import {
    LoginFormModel,
    SignupFormModel,
    ChangePasswordFormModel,
    PasswordResetFormModel,
    PasswordResetConfirmFormModel,
} from '../../models/auth.model';

export interface SignupRequest {
    email: string;
    password: string;
    name: string;
}

export interface LoginRequest {
    email: string;
    password: string;
}

export interface ChangePasswordRequest {
    currentPassword: string;
    newPassword: string;
}

export interface PasswordResetRequest {
    email: string;
}

export interface PasswordResetConfirmRequest {
    token: string;
    newPassword: string;
}

/**
 * LoginFormModel から LoginRequest を作成します
 */
export function toLoginRequest(form: LoginFormModel): LoginRequest {
    return {
        email: form.email,
        password: form.password,
    };
}

/**
 * SignupFormModel から SignupRequest を作成します
 */
export function toSignupRequest(form: SignupFormModel): SignupRequest {
    return {
        email: form.email,
        password: form.password,
        name: form.name,
    };
}

/**
 * ChangePasswordFormModel から ChangePasswordRequest を作成します
 */
export function toChangePasswordRequest(form: ChangePasswordFormModel): ChangePasswordRequest {
    return {
        currentPassword: form.currentPassword,
        newPassword: form.newPassword,
    };
}

/**
 * PasswordResetFormModel から PasswordResetRequest を作成します
 */
export function toPasswordResetRequest(form: PasswordResetFormModel): PasswordResetRequest {
    return {
        email: form.email,
    };
}

/**
 * PasswordResetConfirmFormModel から PasswordResetConfirmRequest を作成します
 */
export function toPasswordResetConfirmRequest(
    form: PasswordResetConfirmFormModel
): PasswordResetConfirmRequest {
    return {
        token: form.token,
        newPassword: form.newPassword,
    };
}
