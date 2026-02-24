import { Component, inject, ViewChild } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { RecipeService } from '../../../core/services/recipe.service';
import {
  RecipeFormComponent,
  RecipeFormSubmitEvent,
} from '../../../shared/components/organisms/recipe-form/recipe-form.component';

@Component({
  selector: 'app-recipe-create',
  standalone: true,
  imports: [CommonModule, RecipeFormComponent],
  template: `
    <app-recipe-form
      titleKey="RECIPE.NEW_TITLE"
      submitLabelKey="RECIPE.SAVE"
      submittingLabelKey="RECIPE.SAVING"
      [isSubmitting]="isSubmitting"
      (save)="onSave($event)"
    ></app-recipe-form>
  `,
})
export class RecipeCreateComponent {
  private readonly recipeService = inject(RecipeService);
  private readonly router = inject(Router);

  @ViewChild(RecipeFormComponent) recipeFormComponent!: RecipeFormComponent;

  isSubmitting = false;

  onSave(event: RecipeFormSubmitEvent) {
    this.isSubmitting = true;

    const requestData = {
      name: event.formData.name,
      url: event.formData.url,
      memo: event.formData.memo,
      tagIds: event.formData.tagIds || [],
    };

    this.recipeService.createRecipe(requestData).subscribe({
      next: (recipe) => {
        if (event.file) {
          this.recipeService.uploadImage(recipe.id, event.file).subscribe({
            next: () => this.router.navigate(['/recipes']),
            error: (err) => {
              console.error('Failed to upload image', err);
              this.router.navigate(['/recipes']);
            },
          });
        } else {
          this.router.navigate(['/recipes']);
        }
      },
      error: (err) => {
        console.error('Failed to create recipe', err);
        this.isSubmitting = false;
        this.recipeFormComponent.handleServerErrors(err);
      },
    });
  }
}
