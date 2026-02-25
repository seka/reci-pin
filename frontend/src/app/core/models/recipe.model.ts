export interface Tag {
    id: number;
    name: string;
}

export interface RecipeImage {
    id: number;
    recipeId: number;
    imageUrl: string;
    createdAt: string;
}

export interface Recipe {
    id: number;
    userId: number;
    name: string;
    url: string;
    memo: string;
    createdAt: string;
    updatedAt: string;
    tags?: Tag[];
    images?: RecipeImage[];
}
