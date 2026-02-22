import { Component, inject, OnInit, ElementRef, ViewChild } from '@angular/core';
import { RouterModule } from '@angular/router';
import { FormsModule, ReactiveFormsModule, FormControl } from '@angular/forms';
import { CommonModule, AsyncPipe } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
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

@Component({
  selector: 'app-recipes',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    FormsModule,
    ReactiveFormsModule,
    MatIconModule,
    MatButtonToggleModule,
    MatChipsModule,
    MatAutocompleteModule,
    MatFormFieldModule,
    TranslatePipe,
    AsyncPipe,
    RecipeCardComponent,
    HeadlineComponent,
    ButtonComponent,
    InputComponent
  ],
  template: `
    <div class="recipes-container">
      <div class="header">
        <app-headline variant="h2">{{ 'RECIPE.MY_RECIPES' | translate }}</app-headline>
        <div class="header-actions">
          <app-button routerLink="/recipes/new" variant="primary" class="add-btn">
            <mat-icon style="vertical-align: middle; margin-right: 4px;">add</mat-icon>
            {{ 'RECIPE.ADD_NEW' | translate }}
          </app-button>
          <a routerLink="/settings" class="settings-link" title="設定">
            <mat-icon>settings</mat-icon>
          </a>
        </div>
      </div>

      <div class="search-section">
        <div class="search-mode-toggle">
          <mat-button-toggle-group [(ngModel)]="searchMode">
            <mat-button-toggle value="keyword">キーワード</mat-button-toggle>
            <mat-button-toggle value="tag">タグ</mat-button-toggle>
          </mat-button-toggle-group>
        </div>

        <div class="search-row" *ngIf="searchMode === 'keyword'">
          <app-input
            [(ngModel)]="searchQuery"
            (keyup.enter)="search()"
            label="キーワードで絞り込み"
            floatLabel="always"
            placeholder="キーワードで検索..."
            class="search-input"
          ></app-input>
          <app-button (click)="search()" variant="secondary" class="search-btn">
            <mat-icon style="font-size: 18px; width: 18px; height: 18px; vertical-align: middle; margin-right: 4px;">search</mat-icon>
            検索
          </app-button>
        </div>

        <div class="search-row" *ngIf="searchMode === 'tag'">
           <mat-form-field class="tag-chip-list" appearance="outline" floatLabel="always">
            <mat-label>タグで絞り込み</mat-label>
            <mat-chip-grid #chipGrid aria-label="Tag selection">
              @for (tagId of selectedTagIds; track tagId) {
                <mat-chip-row (removed)="removeTag(tagId)">
                  {{ getTagName(tagId) }}
                  <button matChipRemove [attr.aria-label]="'remove ' + getTagName(tagId)">
                    <mat-icon>cancel</mat-icon>
                  </button>
                </mat-chip-row>
              }
              <input
                placeholder="タグを選択..."
                #tagInput
                [formControl]="tagCtrl"
                [matChipInputFor]="chipGrid"
                [matAutocomplete]="auto"
                [matChipInputSeparatorKeyCodes]="separatorKeysCodes"
                (matChipInputTokenEnd)="appendTag($event)"
              />
            </mat-chip-grid>
            <mat-autocomplete #auto="matAutocomplete" (optionSelected)="selectedTag($event)">
              @for (tag of filteredTags | async; track tag.id) {
                <mat-option [value]="tag">
                  {{ tag.name }}
                </mat-option>
              }
            </mat-autocomplete>
          </mat-form-field>
          <app-button (click)="search()" variant="secondary" class="search-btn">
            <mat-icon style="font-size: 18px; width: 18px; height: 18px; vertical-align: middle; margin-right: 4px;">search</mat-icon>
            検索
          </app-button>
        </div>
      </div>

      <div class="recipes-grid">
        @for (recipe of recipes; track recipe.id) {
          <app-recipe-card [recipe]="recipe" class="recipe-card-item"></app-recipe-card>
        }
      </div>
    </div>
  `,
  styles: [
    `
      .recipes-container {
        padding: var(--spacing-3);
        max-width: 1200px;
        margin: 0 auto;
      }
      .header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: var(--spacing-3);
      }
      .search-section {
        background-color: var(--color-surface);
        padding: var(--spacing-3);
        border-radius: var(--radius-2);
        margin-bottom: var(--spacing-3);
        box-shadow: var(--shadow-1);
        display: flex;
        flex-direction: column;
        gap: var(--spacing-2);
      }
      .search-mode-toggle {
        display: flex;
        justify-content: flex-start;
      }
      .search-row {
        display: flex;
        gap: var(--spacing-2);
        align-items: center;
        width: 100%;
      }
      .search-input {
        flex: 1;
      }
      .tag-chip-list {
        flex: 1;
      }
      /* Override mat-form-field bottom margin for both inputs to prevent layout shift */
      .tag-chip-list ::ng-deep .mat-mdc-form-field-subscript-wrapper,
      .search-input ::ng-deep .mat-mdc-form-field-subscript-wrapper {
        display: none;
      }
      .header-actions {
        display: flex;
        align-items: center;
        gap: var(--spacing-2);
      }
      .add-btn {
        width: auto;
      }
      .settings-link {
        color: var(--color-text-secondary);
        display: flex;
        align-items: center;
        transition: color 0.2s;
      }
      .settings-link:hover {
        color: var(--color-primary);
      }
      .recipes-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: var(--spacing-3);
      }
      .recipe-card-item {
        height: 100%;
        display: block;
      }
    `,
  ],
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
        return filterValue ? this._filter(filterValue) : this.getUnselectedTags();
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

  private _filter(value: string): Tag[] {
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
