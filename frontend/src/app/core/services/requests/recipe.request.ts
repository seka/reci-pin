import { Recipe } from '../../models/recipe.model';

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

/**
 * Partial<Recipe> から CreateRecipeRequest を作成します
 */
export function fromRecipeModelToCreateRequest(model: Partial<Recipe>): CreateRecipeRequest {
    return {
        name: model.name || '',
        url: model.url || '',
        memo: model.memo || '',
        tagIds: model.tags?.map((t) => t.id) || [],
    };
}

/**
 * Partial<Recipe> から UpdateRecipeRequest を作成します
 */
export function fromRecipeModelToUpdateRequest(model: Partial<Recipe>): UpdateRecipeRequest {
    return {
        name: model.name || '',
        url: model.url || '',
        memo: model.memo || '',
    };
}
