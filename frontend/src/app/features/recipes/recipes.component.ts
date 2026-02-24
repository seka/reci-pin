import { Component, inject, OnInit, ElementRef, ViewChild } from '@angular/core';
import { RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule, FormControl } from '@angular/forms';
import { CommonModule, AsyncPipe } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { SearchModeToggleComponent } from '../../shared/components/molecules/search-mode-toggle/search-mode-toggle.component';
import { MatChipsModule } from '@angular/material/chips';
import { MatAutocompleteModule, MatAutocompleteSelectedEvent } from '@angular/material/autocomplete';
import { MatFormFieldModule } from '@angular/material/form-field';
import { TranslatePipe } from '@ngx-translate/core';
import { Observable, startWith, map } from 'rxjs';
import { COMMA, ENTER } from '@angular/cdk/keycodes';

import { RecipeService, Recipe, Tag } from '../../core/services/recipe.service';
import { RecipeCardComponent } from '../../shared/components/organisms/recipe-card/recipe-card.component';
import { HeadlineComponent } from '../../shared/components/atoms/headline/headline.component';
import { ButtonComponent } from '../../shared/components/atoms/button/button.component';
import { InputComponent } from '../../shared/components/atoms/input/input.component';
import { SearchBarComponent } from '../../shared/components/molecules/search-bar/search-bar.component';
import { EmptyStateComponent } from '../../shared/components/molecules/empty-state/empty-state.component';

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
    TranslatePipe,
    AsyncPipe,
    RecipeCardComponent,
    HeadlineComponent,
    ButtonComponent,
    InputComponent,
    SearchBarComponent,
    EmptyStateComponent
  ],
  templateUrl: './recipes.component.html',
  styleUrl: './recipes.component.scss',
})
export class RecipesComponent implements OnInit {
  private readonly recipeService = inject(RecipeService);

  recipes: Recipe[] = [];
  availableTags: Tag[] = [];
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
        const filterValue = typeof tag === 'string' ? tag : (tag?.name || '');
        return filterValue ? this.filterTags(filterValue) : this.getUnselectedTags();
      })
    );
  }

  ngOnInit() {
    this.loadRecipes();
    this.loadTags();
  }

  loadRecipes() {
    this.recipeService.getUserRecipes().subscribe({
      next: (recipes) => (this.recipes = recipes),
      error: (err: Error) => console.error('Failed to load recipes', err),
    });
  }

  loadTags() {
    this.recipeService.getAllTags().subscribe({
      next: (tags) => (this.availableTags = tags),
      error: (err: Error) => console.error('Failed to load tags', err),
    });
  }



  search() {
    const query = this.searchMode === 'keyword' ? this.searchQuery : '';
    const tagIds = this.searchMode === 'tag' ? this.selectedTagIds : [];

    this.recipeService.searchRecipes({
      query: query || '',
      tagIds: tagIds,
    }).subscribe({
      next: (recipes) => (this.recipes = recipes),
      error: (err: Error) => console.error('Failed to search recipes', err),
    });
  }

  onDeleteRecipe(id: number): void {
    this.recipeService.deleteRecipe(id).subscribe({
      next: () => {
        this.recipes = this.recipes.filter(r => r.id !== id);
      },
      error: (err: Error) => console.error('Failed to delete recipe', err),
    });
  }

  // --- Tag Logic ---

  appendTag(event: any): void {
    const value = (event.value || '').trim();

    // If matches an existing tag, select it
    if (value) {
      const existingTag = this.availableTags.find(
        tag => tag.name.toLowerCase() === value.toLowerCase()
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
    return this.availableTags.find((t) => t.id === id)?.name || '';
  }

  private filterTags(value: string): Tag[] {
    const filterValue = value.toLowerCase();
    return this.availableTags.filter((tag) =>
      tag.name.toLowerCase().includes(filterValue) &&
      !this.selectedTagIds.includes(tag.id)
    );
  }

  private getUnselectedTags(): Tag[] {
    return this.availableTags.filter(tag => !this.selectedTagIds.includes(tag.id));
  }
}
