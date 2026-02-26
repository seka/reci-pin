/**
 * @vitest-environment jsdom
 */
import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting, HttpTestingController } from '@angular/common/http/testing';
import { RecipeService } from './recipe.service';
import { Recipe, Tag, RecipeImage, RecipeFormModel } from '../models/recipe.model';
import { vi, expect, describe, it, beforeEach, afterEach } from 'vitest';

import { lastValueFrom } from 'rxjs';

describe('RecipeService', () => {
    let service: RecipeService;
    let httpMock: HttpTestingController;

    const mockRecipe = new Recipe({
        id: 1,
        userId: 1,
        name: 'Test Recipe',
        url: 'https://example.com',
        memo: 'Test memo',
        createdAt: '',
        updatedAt: '',
    });

    const mockTag = new Tag({ id: 1, name: 'Test Tag' });

    beforeEach(() => {
        TestBed.configureTestingModule({
            providers: [RecipeService, provideHttpClient(), provideHttpClientTesting()],
        });
        service = TestBed.inject(RecipeService);
        httpMock = TestBed.inject(HttpTestingController);
    });

    afterEach(() => {
        httpMock.verify();
    });

    it('should be created', () => {
        expect(service).toBeTruthy();
    });

    it('should create recipe', () => {
        const requestData: RecipeFormModel = { name: 'New', url: '', memo: '', tagIds: [] };
        service.createRecipe(requestData).subscribe((recipe) => {
            expect(recipe).toEqual(mockRecipe);
        });

        const req = httpMock.expectOne('/api/recipes');
        expect(req.request.method).toBe('POST');
        req.flush(mockRecipe);
    });

    it('should get recipe', () => {
        service.getRecipe(1).subscribe((recipe) => {
            expect(recipe).toEqual(mockRecipe);
        });

        const req = httpMock.expectOne('/api/recipes/1');
        expect(req.request.method).toBe('GET');
        req.flush(mockRecipe);
    });

    it('should get user recipes', () => {
        service.getUserRecipes().subscribe((recipes) => {
            expect(recipes).toEqual([mockRecipe]);
        });

        const req = httpMock.expectOne('/api/recipes');
        expect(req.request.method).toBe('GET');
        req.flush([mockRecipe]);
    });

    it('should search recipes', () => {
        const searchData = { query: 'test', tagIds: [] };
        service.searchRecipes(searchData).subscribe((recipes) => {
            expect(recipes).toEqual([mockRecipe]);
        });

        const req = httpMock.expectOne('/api/recipes/search');
        expect(req.request.method).toBe('POST');
        req.flush([mockRecipe]);
    });

    it('should update recipe', () => {
        const updateData: RecipeFormModel = { name: 'Updated', url: '', memo: '', tagIds: [] };
        service.updateRecipe(1, updateData).subscribe((recipe) => {
            expect(recipe).toEqual(mockRecipe);
        });

        const req = httpMock.expectOne('/api/recipes/1');
        expect(req.request.method).toBe('PUT');
        req.flush(mockRecipe);
    });

    it('should delete recipe', () => {
        service.deleteRecipe(1).subscribe();

        const req = httpMock.expectOne('/api/recipes/1');
        expect(req.request.method).toBe('DELETE');
        req.flush(null);
    });

    it('should add tags', () => {
        service.addTags(1, [1, 2]).subscribe();

        const req = httpMock.expectOne('/api/recipes/1/tags');
        expect(req.request.method).toBe('POST');
        expect(req.request.body).toEqual({ tagIds: [1, 2] });
        req.flush(null);
    });

    it('should remove tags', () => {
        service.removeTags(1, [1, 2]).subscribe();

        const req = httpMock.expectOne('/api/recipes/1/tags');
        expect(req.request.method).toBe('DELETE');
        expect(req.request.body).toEqual({ tagIds: [1, 2] });
        req.flush(null);
    });

    it('should create tag', () => {
        service.createTag({ name: 'New Tag' }).subscribe((tag) => {
            expect(tag).toEqual(mockTag);
        });

        const req = httpMock.expectOne('/api/tags');
        expect(req.request.method).toBe('POST');
        req.flush(mockTag);
    });

    it('should get all tags', () => {
        service.getAllTags().subscribe((tags) => {
            expect(tags).toEqual([mockTag]);
        });

        const req = httpMock.expectOne('/api/tags');
        expect(req.request.method).toBe('GET');
        req.flush([mockTag]);
    });

    it('should delete tag', () => {
        service.deleteTag(1).subscribe();

        const req = httpMock.expectOne('/api/tags/1');
        expect(req.request.method).toBe('DELETE');
        req.flush(null);
    });

    it('should upload image', async () => {
        const file = new File([''], 'test.jpg', { type: 'image/jpeg' });
        const mockImage: RecipeImage = { id: 1, recipeId: 1, imageUrl: 'url', createdAt: '' };
        const mockResponse = { image: mockImage, uploadUrl: 'http://upload' };

        const fetchSpy = vi.spyOn(globalThis, 'fetch').mockResolvedValue({
            ok: true,
        } as any);

        const uploadPromise = lastValueFrom(service.uploadImage(1, file));

        const req = httpMock.expectOne('/api/recipes/1/images');
        expect(req.request.method).toBe('POST');
        req.flush(mockResponse);

        const result = await uploadPromise;

        expect(result).toEqual(mockImage);
        expect(fetchSpy).toHaveBeenCalledWith('http://upload', expect.objectContaining({
            method: 'PUT',
            body: file
        }));

        fetchSpy.mockRestore();
    });
});
