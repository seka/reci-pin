import {
    Recipe,
    RecipeFormModel,
    RecipeSearchFormModel,
    TagFormModel,
} from '../../models/recipe.model';

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

export interface TagRequest {
    name: string;
}

/**
 * RecipeFormModel から CreateRecipeRequest を作成します
 */
export function fromRecipeFormToCreateRequest(form: RecipeFormModel): CreateRecipeRequest {
    return {
        name: form.name,
        url: form.url,
        memo: form.memo,
        tagIds: form.tagIds,
    };
}

/**
 * RecipeFormModel から UpdateRecipeRequest を作成します
 */
export function fromRecipeFormToUpdateRequest(form: RecipeFormModel): UpdateRecipeRequest {
    return {
        name: form.name,
        url: form.url,
        memo: form.memo,
    };
}

/**
 * RecipeSearchFormModel から SearchRecipeRequest を作成します
 */
export function fromRecipeSearchFormToSearchRequest(form: RecipeSearchFormModel): SearchRecipeRequest {
    return {
        query: form.query,
        tagIds: form.tagIds,
    };
}

/**
 * TagFormModel から TagRequest を作成します
 */
export function toTagRequest(form: TagFormModel): TagRequest {
    return {
        name: form.name,
    };
}
