export interface TagData {
    id: number;
    name: string;
}

export class Tag {
    id: number;
    name: string;

    constructor(data: TagData) {
        this.id = data.id;
        this.name = data.name;
    }
}

export interface RecipeImageData {
    id: number;
    recipeId: number;
    imageUrl: string;
    createdAt: string;
}

export class RecipeImage {
    id: number;
    recipeId: number;
    imageUrl: string;
    createdAt: string;

    constructor(data: RecipeImageData) {
        this.id = data.id;
        this.recipeId = data.recipeId;
        this.imageUrl = data.imageUrl;
        this.createdAt = data.createdAt;
    }
}

export interface RecipeData {
    id: number;
    userId: number;
    name: string;
    url: string;
    memo: string;
    createdAt: string;
    updatedAt: string;
    tags?: TagData[];
    images?: RecipeImageData[];
}

export class Recipe {
    id: number;
    userId: number;
    name: string;
    url: string;
    memo: string;
    createdAt: string;
    updatedAt: string;
    tags: Tag[] = [];
    images: RecipeImage[] = [];

    constructor(data: RecipeData) {
        this.id = data.id;
        this.userId = data.userId;
        this.name = data.name;
        this.url = data.url;
        this.memo = data.memo;
        this.createdAt = data.createdAt;
        this.updatedAt = data.updatedAt;

        if (data.tags) {
            this.tags = data.tags.map(t => new Tag(t));
        }
        if (data.images) {
            this.images = data.images.map(i => new RecipeImage(i));
        }
    }
}

export interface RecipeFormModel {
    name: string;
    url: string;
    memo: string;
    tagIds: number[];
}
