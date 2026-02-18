import { Component, inject, OnInit } from '@angular/core';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { TranslatePipe } from '@ngx-translate/core';
import { RecipeService, Recipe, Tag } from '../../core/services/recipe.service';
import { RecipeCardComponent } from '../../shared/components/organisms/recipe-card/recipe-card.component';
import { HeadlineComponent } from '../../shared/components/atoms/headline/headline.component';
import { ButtonComponent } from '../../shared/components/atoms/button/button.component';
import { TagSelectComponent } from '../../shared/components/molecules/tag-select/tag-select.component';
import { InputComponent } from '../../shared/components/atoms/input/input.component';

@Component({
  selector: 'app-recipes',
  standalone: true,
  imports: [
    RouterModule,
    FormsModule,
    MatIconModule,
    TranslatePipe,
    RecipeCardComponent,
    HeadlineComponent,
    ButtonComponent,
    TagSelectComponent,
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
        <div class="search-row">
          <app-input
            [(ngModel)]="searchQuery"
            placeholder="キーワードで検索..."
            class="search-input"
          ></app-input>
          <app-button (click)="search()" variant="secondary" class="search-btn">
            <mat-icon>search</mat-icon>
            検索
          </app-button>
        </div>
        <app-tag-select
          [(ngModel)]="selectedTagIds"
          [tags]="availableTags"
        ></app-tag-select>
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
      }
      .search-row {
        display: flex;
        gap: var(--spacing-2);
        align-items: flex-start;
        margin-bottom: var(--spacing-2);
      }
      .search-input {
        flex: 1;
      }
      .search-btn {
        height: 56px; /* Match standard input height */
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
    this.recipeService.searchRecipes({
      query: this.searchQuery,
      tag_ids: this.selectedTagIds
    }).subscribe({
      next: (recipes) => (this.recipes = recipes),
      error: (err: Error) => console.error('Failed to search recipes', err),
    });
  }
}
