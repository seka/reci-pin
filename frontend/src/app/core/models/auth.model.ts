export interface LoginFormModel {
    email: string;
    password: string;
}

export interface SignupFormModel {
    email: string;
    password: string;
    name: string;
}

export interface ChangePasswordFormModel {
    currentPassword: string;
    newPassword: string;
}

export interface PasswordResetFormModel {
    email: string;
}

export interface PasswordResetConfirmFormModel {
    token: string;
    newPassword: string;
}
