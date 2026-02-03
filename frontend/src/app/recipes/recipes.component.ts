import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { RecipeService, Recipe } from '../core/services/recipe.service';

@Component({
  selector: 'app-recipes',
  standalone: true,
  imports: [CommonModule, RouterModule],
  template: `
    <div class="recipes-container">
      <h2>マイレシピ</h2>
      <button routerLink="/recipes/new">新規レシピ追加</button>
      <div class="recipes-list">
        <div *ngFor="let recipe of recipes" class="recipe-card">
          <h3>{{ recipe.name }}</h3>
          <p>{{ recipe.memo }}</p>
          <a *ngIf="recipe.url" [href]="getExternalUrl(recipe.url)" target="_blank" rel="noopener noreferrer">レシピを見る</a>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .recipes-container { padding: 20px; }
    button { padding: 10px 20px; background: #007bff; color: white; border: none; cursor: pointer; margin-bottom: 20px; }
    .recipe-card { border: 1px solid #ddd; padding: 15px; margin: 10px 0; border-radius: 5px; }
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

  getExternalUrl(url: string): string {
    if (!url) return '';
    if (/^https?:\/\//i.test(url)) {
      return url;
    }
    return 'https://' + url;
  }
}
