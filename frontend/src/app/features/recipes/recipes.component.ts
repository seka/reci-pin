import {
  Component,
  inject,
  Injector,
  signal,
  ElementRef,
  ViewChild,
  ChangeDetectionStrategy,
} from '@angular/core';
import { rxResource } from '@angular/core/rxjs-interop';
import { RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule, FormControl } from '@angular/forms';
import { CommonModule, AsyncPipe } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { SearchModeToggleComponent } from '../../shared/components/molecules/search-mode-toggle/search-mode-toggle.component';
import { MatChipsModule, MatChipInputEvent } from '@angular/material/chips';
import {
  MatAutocompleteModule,
  MatAutocompleteSelectedEvent,
} from '@angular/material/autocomplete';
import { MatFormFieldModule } from '@angular/material/form-field';
import { TranslocoPipe } from '@jsverse/transloco';
import { Observable, startWith, map } from 'rxjs';
import { COMMA, ENTER } from '@angular/cdk/keycodes';

import { RecipeService } from '../../core/services/recipe.service';
import { Tag } from '../../core/models/recipe.model';
import { RecipeCardComponent } from '../../shared/components/organisms/recipe-card/recipe-card.component';
import { HeadlineComponent } from '../../shared/components/atoms/headline/headline.component';
import { ButtonComponent } from '../../shared/components/atoms/button/button.component';
import { InputComponent } from '../../shared/components/atoms/input/input.component';
import { SearchBarComponent } from '../../shared/components/molecules/search-bar/search-bar.component';
import { EmptyStateComponent } from '../../shared/components/molecules/empty-state/empty-state.component';
import { LoadingStateComponent } from '../../shared/components/molecules/loading-state/loading-state.component';

@Component({
  selector: 'app-recipes',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    FormsModule,
    ReactiveFormsModule,
    MatIconModule,
    SearchModeToggleComponent,
    MatChipsModule,
    MatAutocompleteModule,
    MatFormFieldModule,
    TranslocoPipe,
    AsyncPipe,
    RecipeCardComponent,
    HeadlineComponent,
    ButtonComponent,
    InputComponent,
    SearchBarComponent,
    EmptyStateComponent,
    LoadingStateComponent,
  ],
  templateUrl: './recipes.component.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './recipes.component.scss',
})
export class RecipesComponent {
  private readonly recipeService = inject(RecipeService);
  private readonly injector = inject(Injector);

  // null means "no search criteria" -> load the unfiltered recipe list
  private readonly recipesParams = signal<{ query: string; tagIds: number[] } | null>(null);

  protected readonly recipesResource = rxResource({
    injector: this.injector,
    params: () => this.recipesParams(),
    defaultValue: [],
    stream: ({ params }) =>
      params === null
        ? this.recipeService.getUserRecipes()
        : this.recipeService.searchRecipes(params),
  });

  protected readonly tagsResource = rxResource({
    injector: this.injector,
    defaultValue: [],
    stream: () => this.recipeService.getAllTags(),
  });

  searchQuery = '';
  selectedTagIds: number[] = [];

  // Search Mode
  searchMode: 'keyword' | 'tag' = 'keyword';

  // Tag Chips Logic
  separatorKeysCodes: number[] = [ENTER, COMMA];
  tagCtrl = new FormControl('');
  filteredTags: Observable<Tag[]>;
  @ViewChild('tagInput') tagInput!: ElementRef<HTMLInputElement>;

  constructor() {
    this.filteredTags = this.tagCtrl.valueChanges.pipe(
      startWith(null),
      map((tag: string | null | Tag) => {
        // Handle both string input and Tag object selection
        const filterValue = typeof tag === 'string' ? tag : tag?.name || '';
        return filterValue ? this.filterTags(filterValue) : this.getUnselectedTags();
      }),
    );
  }

  search() {
    const query = this.searchMode === 'keyword' ? this.searchQuery : '';
    const tagIds = this.searchMode === 'tag' ? this.selectedTagIds : [];

    this.recipesParams.set({ query: query || '', tagIds });
  }

  onDeleteRecipe(id: number): void {
    this.recipeService.deleteRecipe(id).subscribe({
      next: () => {
        this.recipesResource.update((recipes) => recipes.filter((r) => r.id !== id));
      },
      error: (err: Error) => console.error('Failed to delete recipe', err),
    });
  }

  // --- Tag Logic ---

  appendTag(event: MatChipInputEvent): void {
    const value = (event.value || '').trim();

    // If matches an existing tag, select it
    if (value) {
      const availableTags = this.tagsResource.value();
      const existingTag = availableTags.find(
        (tag) => tag.name.toLowerCase() === value.toLowerCase(),
      );

      if (existingTag && !this.selectedTagIds.includes(existingTag.id)) {
        this.selectedTagIds.push(existingTag.id);
      }
    }

    // Always clear input
    if (event.chipInput) {
      event.chipInput.clear();
    }
    this.tagCtrl.setValue(null);
  }

  removeTag(tagId: number): void {
    const index = this.selectedTagIds.indexOf(tagId);
    if (index >= 0) {
      this.selectedTagIds.splice(index, 1);
    }
  }

  selectedTag(event: MatAutocompleteSelectedEvent): void {
    const tag: Tag = event.option.value;
    if (!this.selectedTagIds.includes(tag.id)) {
      this.selectedTagIds.push(tag.id);
    }
    this.tagInput.nativeElement.value = '';
    this.tagCtrl.setValue(null);
  }

  getTagName(id: number): string {
    const availableTags = this.tagsResource.value();
    return availableTags.find((t) => t.id === id)?.name || '';
  }

  private filterTags(value: string): Tag[] {
    const filterValue = value.toLowerCase();
    const availableTags = this.tagsResource.value();
    return availableTags.filter(
      (tag) =>
        tag.name.toLowerCase().includes(filterValue) && !this.selectedTagIds.includes(tag.id),
    );
  }

  private getUnselectedTags(): Tag[] {
    const availableTags = this.tagsResource.value();
    return availableTags.filter((tag) => !this.selectedTagIds.includes(tag.id));
  }
}
