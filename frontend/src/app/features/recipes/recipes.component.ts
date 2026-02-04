import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { RecipeService, Recipe } from '../../core/services/recipe.service';
import { RecipeCardComponent } from '../../shared/components/organisms/recipe-card/recipe-card.component';

@Component({
  selector: 'app-recipes',
  standalone: true,
  imports: [CommonModule, RouterModule, MatButtonModule, MatIconModule, RecipeCardComponent],
  template: `
    <div class="recipes-container">
      <div class="header">
        <h2>マイレシピ</h2>
        <button mat-flat-button color="primary" routerLink="/recipes/new">
          <mat-icon>add</mat-icon>
          新規レシピ追加
        </button>
      </div>
      
      <div class="recipes-grid">
        <app-recipe-card *ngFor="let recipe of recipes" [recipe]="recipe" class="recipe-card-item"></app-recipe-card>
      </div>
    </div>
  `,
  styles: [`
    .recipes-container { padding: 24px; max-width: 1200px; margin: 0 auto; }
    .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px; }
    h2 { font-size: 2rem; color: #333; margin: 0; font-weight: 700; color: #e91e63; }
    .recipes-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 24px; }
    .recipe-card-item { height: 100%; display: block; } /* Ensure height usage */
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
