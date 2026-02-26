import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, switchMap, from, map } from 'rxjs';
import { Recipe, Tag, RecipeImage, RecipeFormModel, RecipeSearchFormModel, TagFormModel } from '../models/recipe.model';
import {
  CreateRecipeRequest,
  UpdateRecipeRequest,
  SearchRecipeRequest,
  CreateRecipeImageRequest,
  toCreateRecipeRequest,
  toUpdateRecipeRequest,
  toSearchRecipeRequest,
  toTagRequest,
} from './requests/recipe.request';
import {
  RecipeResponse,
  CreateRecipeImageResponse,
  toRecipeModel,
  toRecipeModels,
  toRecipeImageModel,
} from './responses/recipe.response';

@Injectable({
  providedIn: 'root',
})
export class RecipeService {
  private readonly API_URL = '/api/recipes';
  private readonly TAG_API_URL = '/api/tags';
  private readonly http = inject(HttpClient);

  createRecipe(form: RecipeFormModel): Observable<Recipe> {
    const data = toCreateRecipeRequest(form);
    return this.http.post<RecipeResponse>(this.API_URL, data).pipe(map(toRecipeModel));
  }

  getRecipe(id: number): Observable<Recipe> {
    return this.http.get<RecipeResponse>(`${this.API_URL}/${id}`).pipe(map(toRecipeModel));
  }

  getUserRecipes(): Observable<Recipe[]> {
    return this.http.get<RecipeResponse[]>(this.API_URL).pipe(map(toRecipeModels));
  }

  searchRecipes(form: RecipeSearchFormModel): Observable<Recipe[]> {
    const data = toSearchRecipeRequest(form);
    return this.http.post<RecipeResponse[]>(`${this.API_URL}/search`, data).pipe(map(toRecipeModels));
  }

  updateRecipe(id: number, form: RecipeFormModel): Observable<Recipe> {
    const data = toUpdateRecipeRequest(form);
    return this.http.put<RecipeResponse>(`${this.API_URL}/${id}`, data).pipe(map(toRecipeModel));
  }

  deleteRecipe(id: number): Observable<void> {
    return this.http.delete<void>(`${this.API_URL}/${id}`);
  }

  addTags(recipeId: number, tagIds: number[]): Observable<void> {
    return this.http.post<void>(`${this.API_URL}/${recipeId}/tags`, { tagIds });
  }

  removeTags(recipeId: number, tagIds: number[]): Observable<void> {
    return this.http.request<void>('delete', `${this.API_URL}/${recipeId}/tags`, {
      body: { tagIds },
    });
  }

  uploadImage(recipeId: number, file: File): Observable<RecipeImage> {
    const body: CreateRecipeImageRequest = {
      filename: file.name,
      contentType: file.type,
      size: file.size,
    };

    return this.http
      .post<CreateRecipeImageResponse>(`${this.API_URL}/${recipeId}/images`, body)
      .pipe(
        switchMap((res) => {
          const uploadUrl = res.uploadUrl;
          return from(
            fetch(uploadUrl, {
              method: 'PUT',
              headers: { 'Content-Type': file.type },
              body: file,
            }).then((r) => {
              if (!r.ok) throw new Error(`S3 upload failed: ${r.status}`);
              return toRecipeImageModel(res);
            }),
          );
        }),
      );
  }

  // Tag management
  createTag(form: TagFormModel): Observable<Tag> {
    const data = toTagRequest(form);
    return this.http.post<Tag>(this.TAG_API_URL, data);
  }

  getAllTags(): Observable<Tag[]> {
    return this.http.get<Tag[]>(this.TAG_API_URL);
  }

  deleteTag(id: number): Observable<void> {
    return this.http.delete<void>(`${this.TAG_API_URL}/${id}`);
  }
}
