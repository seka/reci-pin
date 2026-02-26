import {
    Recipe,
    RecipeData,
    RecipeImage,
    RecipeImageData,
} from '../../models/recipe.model';

export type RecipeResponse = RecipeData;

export interface CreateRecipeImageResponse {
    image: RecipeImageData;
    uploadUrl: string;
}

/**
 * RecipeResponse を Recipe モデルに変換します
 */
export function toRecipeModel(res: RecipeResponse): Recipe {
    return new Recipe(res);
}

/**
 * RecipeResponse の配列を Recipe モデルの配列に変換します
 */
export function toRecipeModels(res: RecipeResponse[]): Recipe[] {
    return res.map((item) => new Recipe(item));
}

/**
 * CreateRecipeImageResponse を RecipeImage モデルに変換します
 */
export function toRecipeImageModel(res: CreateRecipeImageResponse): RecipeImage {
    return new RecipeImage(res.image);
}
