import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { RecipeService, Recipe } from '../../core/services/recipe.service';
import { RecipeCardComponent } from '../../shared/components/organisms/recipe-card/recipe-card.component';
import { HeadlineComponent } from '../../shared/components/atoms/headline/headline.component';
import { ButtonComponent } from '../../shared/components/atoms/button/button.component';

@Component({
  selector: 'app-recipes',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatIconModule,
    RecipeCardComponent,
    HeadlineComponent,
    ButtonComponent
  ],
  template: `
    <div class="recipes-container">
      <div class="header">
        <app-headline variant="h2">マイレシピ</app-headline>
        <app-button routerLink="/recipes/new" variant="primary" class="add-btn">
          <mat-icon style="vertical-align: middle; margin-right: 4px;">add</mat-icon>
          新規レシピ追加
        </app-button>
      </div>
      
      <div class="recipes-grid">
        <app-recipe-card *ngFor="let recipe of recipes" [recipe]="recipe" class="recipe-card-item"></app-recipe-card>
      </div>
    </div>
  `,
  styles: [`
    .recipes-container { padding: 24px; max-width: 1200px; margin: 0 auto; }
    .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
    /* app-headline handles the styling now */
    .add-btn { width: auto; } /* Override default full width */
    .recipes-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 24px; }
    .recipe-card-item { height: 100%; display: block; }
  `]
})
export class RecipesComponent implements OnInit {
  recipes: Recipe[] = [];

  constructor(private recipeService: RecipeService) { }

  ngOnInit() {
    this.recipeService.getUserRecipes().subscribe({
      next: (recipes) => this.recipes = recipes,
      error: (err: any) => console.error('Failed to load recipes', err)
    });
  }
}
