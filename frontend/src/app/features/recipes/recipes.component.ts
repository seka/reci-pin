import { Component, inject, OnInit } from '@angular/core';
import { RouterModule } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { TranslatePipe } from '@ngx-translate/core';
import { RecipeService, Recipe } from '../../core/services/recipe.service';
import { RecipeCardComponent } from '../../shared/components/organisms/recipe-card/recipe-card.component';
import { HeadlineComponent } from '../../shared/components/atoms/headline/headline.component';
import { ButtonComponent } from '../../shared/components/atoms/button/button.component';

@Component({
  selector: 'app-recipes',
  standalone: true,
  imports: [RouterModule, MatIconModule, TranslatePipe, RecipeCardComponent, HeadlineComponent, ButtonComponent],
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

  ngOnInit() {
    this.recipeService.getUserRecipes().subscribe({
      next: (recipes) => (this.recipes = recipes),
      error: (err: Error) => console.error('Failed to load recipes', err),
    });
  }
}
