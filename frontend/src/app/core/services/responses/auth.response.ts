import { User, UserData } from '../../models/user.model';

export interface AuthResponse {
    token: string;
    user: UserData;
}

/**
 * AuthResponse を User モデルに変換します
 */
export function toUserModel(res: AuthResponse): User {
    return new User(res.user);
}

export interface MessageResponse {
    message: string;
}
