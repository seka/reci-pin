import { User } from '../../models/user.model';

export interface AuthResponse {
    token: string;
    user: User;
}

/**
 * AuthResponse を User モデルに変換します
 */
export function toUserModel(res: AuthResponse): User {
    return res.user;
}

export interface MessageResponse {
    message: string;
}
