import { inject, Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, switchMap, from } from 'rxjs';

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

export interface Tag {
  id: number;
  name: string;
}

export interface RecipeImage {
  id: number;
  recipeId: number;
  imagePath: string;
  imageUrl: string;
  createdAt: string;
}

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

export interface CreateRecipeImageResponse {
  image: RecipeImage;
  uploadUrl: string;
}

@Injectable({
  providedIn: 'root',
})
export class RecipeService {
  private readonly API_URL = '/api/recipes';
  private readonly TAG_API_URL = '/api/tags';
  private readonly http = inject(HttpClient);

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
              return res.image;
            }),
          );
        }),
      );
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
