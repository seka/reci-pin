export interface CreateRecipeRequest {
    name: string;
    url: string;
    memo: string;
    tagIds: number[];
}

export interface UpdateRecipeRequest {
    name: string;
    url: string;
    memo: string;
}

export interface SearchRecipeRequest {
    query: string;
    tagIds: number[];
}

export interface CreateRecipeImageRequest {
    filename: string;
    contentType: string;
    size: number;
}
