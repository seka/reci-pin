import {
  Component,
  computed,
  effect,
  inject,
  Injector,
  ViewChild,
  ChangeDetectionStrategy,
} from '@angular/core';
import { rxResource } from '@angular/core/rxjs-interop';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { forkJoin, Observable, of, switchMap } from 'rxjs';
import { RecipeService } from '../../../core/services/recipe.service';
import { RecipeImage } from '../../../core/models/recipe.model';
import {
  RecipeFormComponent,
  RecipeFormSubmitEvent,
} from '../../../shared/components/organisms/recipe-form/recipe-form.component';
import { RecipeFormModel } from '../../../core/models/recipe.model';

@Component({
  selector: 'app-recipe-edit',
  standalone: true,
  imports: [CommonModule, RecipeFormComponent],
  templateUrl: './recipe-edit.component.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './recipe-edit.component.scss',
})
export class RecipeEditComponent {
  private readonly recipeService = inject(RecipeService);
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);
  private readonly injector = inject(Injector);

  @ViewChild(RecipeFormComponent) recipeFormComponent!: RecipeFormComponent;

  isSubmitting = false;

  // Resolved once from the route snapshot; null means the route was missing
  // the `id` param, in which case we redirect away without ever calling
  // getRecipe().
  private readonly recipeId: number | null;

  protected readonly recipeResource = rxResource({
    injector: this.injector,
    params: () => this.recipeId,
    stream: ({ params }) =>
      params === null ? of(undefined) : this.recipeService.getRecipe(params),
  });

  // resource.value() throws (ResourceValueError) once the resource has
  // settled into the 'error' status, so every derived computed() here must
  // check hasValue() first rather than relying on optional chaining alone.

  protected readonly initialData = computed<Partial<RecipeFormModel>>(() => {
    if (!this.recipeResource.hasValue()) return {};
    const recipe = this.recipeResource.value();
    if (!recipe) return {};
    return {
      name: recipe.name,
      url: recipe.url,
      memo: recipe.memo,
      tagIds: recipe.tags?.map((t) => t.id) || [],
    };
  });

  private readonly originalTagIds = computed<number[]>(() => {
    if (!this.recipeResource.hasValue()) return [];
    return this.recipeResource.value()?.tags?.map((t) => t.id) || [];
  });

  protected readonly initialImagePreview = computed<string | null>(() => {
    if (!this.recipeResource.hasValue()) return null;
    const recipe = this.recipeResource.value();
    return recipe?.images && recipe.images.length > 0 ? recipe.images[0].imageUrl : null;
  });

  constructor() {
    const idParam = this.route.snapshot.paramMap.get('id');
    if (!idParam) {
      this.recipeId = null;
      this.router.navigate(['/recipes']);
    } else {
      this.recipeId = Number(idParam);
    }

    effect(() => {
      const error = this.recipeResource.error();
      if (error) {
        console.error('Failed to load recipe for edit', error);
        this.router.navigate(['/recipes']);
      }
    });
  }

  onSave(event: RecipeFormSubmitEvent) {
    this.isSubmitting = true;
    const formData = event.formData;
    const originalTagIds = this.originalTagIds();

    this.recipeService
      .updateRecipe(this.recipeId!, formData)
      .pipe(
        switchMap(() => {
          let tagUpdates$: Observable<void[] | RecipeImage | null> = of(null);
          const currentTagIds: number[] = formData.tagIds || [];
          const tagsToAdd = currentTagIds.filter((id) => !originalTagIds.includes(id));
          const tagsToRemove = originalTagIds.filter((id) => !currentTagIds.includes(id));

          const ops = [];
          if (tagsToAdd.length > 0)
            ops.push(this.recipeService.addTags(this.recipeId!, tagsToAdd));
          if (tagsToRemove.length > 0)
            ops.push(this.recipeService.removeTags(this.recipeId!, tagsToRemove));
          if (ops.length > 0) tagUpdates$ = forkJoin(ops);

          return tagUpdates$.pipe(
            switchMap(() => {
              if (event.file) {
                return this.recipeService.uploadImage(this.recipeId!, event.file);
              }
              return of(null);
            }),
          );
        }),
      )
      .subscribe({
        next: () => this.router.navigate(['/recipes']),
        error: (err) => {
          console.error('Failed to save recipe', err);
          this.isSubmitting = false;
          const hasValidationErrors = this.recipeFormComponent.handleServerErrors(err);

          if (!hasValidationErrors) {
            this.router.navigate(['/recipes']);
          }
        },
      });
  }
}
