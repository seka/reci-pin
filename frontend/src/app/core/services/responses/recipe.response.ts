import { Recipe, RecipeImage, Tag } from '../../models/recipe.model';

export interface RecipeResponse {
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

export interface CreateRecipeImageResponse {
    image: RecipeImage;
    uploadUrl: string;
}

/**
 * RecipeResponse を Recipe モデルに変換します
 */
export function toRecipeModel(res: RecipeResponse): Recipe {
    return {
        ...res,
    };
}

/**
 * RecipeResponse の配列を Recipe モデルの配列に変換します
 */
export function toRecipeModels(res: RecipeResponse[]): Recipe[] {
    return res.map(toRecipeModel);
}

/**
 * CreateRecipeImageResponse を RecipeImage モデルに変換します
 */
export function toRecipeImageModel(res: CreateRecipeImageResponse): RecipeImage {
    return res.image;
}
