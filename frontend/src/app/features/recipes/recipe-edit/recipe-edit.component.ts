import { Component, inject, OnInit, ViewChild } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { forkJoin, of, switchMap } from 'rxjs';
import { RecipeService } from '../../../core/services/recipe.service';
import { RecipeFormComponent, RecipeFormSubmitEvent, RecipeFormData } from '../../../shared/components/organisms/recipe-form/recipe-form.component';

@Component({
    selector: 'app-recipe-edit',
    standalone: true,
    imports: [CommonModule, RecipeFormComponent],
  templateUrl: './recipe-edit.component.html',
  styleUrl: './recipe-edit.component.scss',
})
export class RecipeEditComponent implements OnInit {
    private readonly recipeService = inject(RecipeService);
    private readonly router = inject(Router);
    private readonly route = inject(ActivatedRoute);

    @ViewChild(RecipeFormComponent) recipeFormComponent!: RecipeFormComponent;

    isLoading = true;
    isSubmitting = false;
    recipeId!: number;
    originalTagIds: number[] = [];

    initialData: Partial<RecipeFormData> = {};
    initialImagePreview: string | null = null;

    ngOnInit() {
        const idParam = this.route.snapshot.paramMap.get('id');
        if (!idParam) {
            this.router.navigate(['/recipes']);
            return;
        }

        this.recipeId = Number(idParam);
        this.recipeService.getRecipe(this.recipeId).subscribe({
            next: (recipe) => {
                this.initialData = {
                    name: recipe.name,
                    url: recipe.url,
                    memo: recipe.memo,
                    tagIds: recipe.tags?.map((t) => t.id) || [],
                };
                this.originalTagIds = recipe.tags?.map((t) => t.id) || [];
                if (recipe.images && recipe.images.length > 0) {
                    this.initialImagePreview = recipe.images[0].imageUrl;
                }
                this.isLoading = false;
            },
            error: (err) => {
                console.error('Failed to load recipe for edit', err);
                this.router.navigate(['/recipes']);
            },
        });
    }

    onSave(event: RecipeFormSubmitEvent) {
        this.isSubmitting = true;
        const formData = event.formData;

        const requestData = {
            name: formData.name,
            url: formData.url,
            memo: formData.memo,
        };

        this.recipeService.updateRecipe(this.recipeId, requestData).pipe(
            switchMap(() => {
                let tagUpdates$: any = of(null);
                const currentTagIds: number[] = formData.tagIds || [];
                const tagsToAdd = currentTagIds.filter((id) => !this.originalTagIds.includes(id));
                const tagsToRemove = this.originalTagIds.filter((id) => !currentTagIds.includes(id));

                const ops = [];
                if (tagsToAdd.length > 0) ops.push(this.recipeService.addTags(this.recipeId, tagsToAdd));
                if (tagsToRemove.length > 0) ops.push(this.recipeService.removeTags(this.recipeId, tagsToRemove));
                if (ops.length > 0) tagUpdates$ = forkJoin(ops);

                return tagUpdates$.pipe(
                    switchMap(() => {
                        if (event.file) {
                            return this.recipeService.uploadImage(this.recipeId, event.file);
                        }
                        return of(null);
                    })
                );
            })
        ).subscribe({
            next: () => this.router.navigate(['/recipes']),
            error: (err) => {
                console.error('Failed to save recipe', err);
                this.isSubmitting = false;
                let hasValidationErrors = this.recipeFormComponent.handleServerErrors(err);

                if (!hasValidationErrors) {
                    this.router.navigate(['/recipes']);
                }
            },
        });
    }
}
