export interface UserData {
    id: number;
    email: string;
    name: string;
    createdAt: string;
    updatedAt: string;
}

export class User {
    id: number;
    email: string;
    name: string;
    createdAt: string;
    updatedAt: string;

    constructor(data: UserData) {
        this.id = data.id;
        this.email = data.email;
        this.name = data.name;
        this.createdAt = data.createdAt;
        this.updatedAt = data.updatedAt;
    }
}
