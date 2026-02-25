import { RecipeImage } from '../../models/recipe.model';

export interface CreateRecipeImageResponse {
    image: RecipeImage;
    uploadUrl: string;
}
