import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Recipe {
    id: number;
    user_id: number;
    name: string;
    url: string;
    memo: string;
    created_at: string;
    updated_at: string;
    tags?: Tag[];
    images?: RecipeImage[];
}

export interface Tag {
    id: number;
    name: string;
}

export interface RecipeImage {
    id: number;
    recipe_id: number;
    image_path: string;
    created_at: string;
}

export interface CreateRecipeRequest {
    name: string;
    url: string;
    memo: string;
    tag_ids: number[];
}

export interface UpdateRecipeRequest {
    name: string;
    url: string;
    memo: string;
}

export interface SearchRecipeRequest {
    query: string;
    tag_ids: number[];
}

@Injectable({
    providedIn: 'root'
})
export class RecipeService {
    private readonly API_URL = '/api/recipes';
    private readonly TAG_API_URL = '/api/tags';

    constructor(private http: HttpClient) { }

    createRecipe(data: CreateRecipeRequest): Observable<Recipe> {
        return this.http.post<Recipe>(this.API_URL, data);
    }

    getRecipe(id: number): Observable<Recipe> {
        return this.http.get<Recipe>(`${this.API_URL}/${id}`);
    }

    getUserRecipes(): Observable<Recipe[]> {
        return this.http.get<Recipe[]>(this.API_URL);
    }

    searchRecipes(data: SearchRecipeRequest): Observable<Recipe[]> {
        return this.http.post<Recipe[]>(`${this.API_URL}/search`, data);
    }

    updateRecipe(id: number, data: UpdateRecipeRequest): Observable<Recipe> {
        return this.http.put<Recipe>(`${this.API_URL}/${id}`, data);
    }

    deleteRecipe(id: number): Observable<void> {
        return this.http.delete<void>(`${this.API_URL}/${id}`);
    }

    addTags(recipeId: number, tagIds: number[]): Observable<void> {
        return this.http.post<void>(`${this.API_URL}/${recipeId}/tags`, { tag_ids: tagIds });
    }

    removeTags(recipeId: number, tagIds: number[]): Observable<void> {
        return this.http.request<void>('delete', `${this.API_URL}/${recipeId}/tags`, { body: { tag_ids: tagIds } });
    }

    addImage(recipeId: number, imagePath: string): Observable<RecipeImage> {
        return this.http.post<RecipeImage>(`${this.API_URL}/${recipeId}/images`, { image_path: imagePath });
    }

    // Tag management
    createTag(name: string): Observable<Tag> {
        return this.http.post<Tag>(this.TAG_API_URL, { name });
    }

    getAllTags(): Observable<Tag[]> {
        return this.http.get<Tag[]>(this.TAG_API_URL);
    }

    deleteTag(id: number): Observable<void> {
        return this.http.delete<void>(`${this.TAG_API_URL}/${id}`);
    }
}
