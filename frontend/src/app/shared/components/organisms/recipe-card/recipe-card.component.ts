import { Component, Input } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { Recipe } from '../../../../core/services/recipe.service';

@Component({
  selector: 'app-recipe-card',
  standalone: true,
  imports: [MatCardModule, MatButtonModule, MatIconModule, MatChipsModule],
  template: `
    <mat-card class="recipe-card">
      <mat-card-header>
        <mat-card-title>{{ recipe.name }}</mat-card-title>
      </mat-card-header>
      <mat-card-content>
        <p>{{ recipe.memo }}</p>
        @if (recipe.tags && recipe.tags.length > 0) {
          <mat-chip-set>
            @for (tag of recipe.tags; track tag.id) {
              <mat-chip>{{ tag.name }}</mat-chip>
            }
          </mat-chip-set>
        }
      </mat-card-content>
      <mat-card-actions align="end">
        @if (recipe.url) {
          <a
            mat-button
            color="accent"
            [href]="getExternalUrl(recipe.url)"
            target="_blank"
            rel="noopener noreferrer"
          >
            レシピを見る
          </a>
        }
      </mat-card-actions>
    </mat-card>
  `,
  styles: [
    `
      .recipe-card {
        height: 100%;
        display: flex;
        flex-direction: column;
      }
      mat-card-title {
        color: var(--color-primary);
        font-weight: 700;
      }
      mat-card-content {
        flex-grow: 1;
        margin-top: var(--spacing-2);
        margin-bottom: var(--spacing-2);
      }
      p {
        color: var(--color-text-secondary);
        line-height: 1.6;
      }
      mat-chip-set {
        margin-top: var(--spacing-1_5);
      }
    `,
  ],
})
export class RecipeCardComponent {
  @Input() recipe!: Recipe;

  getExternalUrl(url: string): string {
    if (!url) return '';
    if (/^https?:\/\//i.test(url)) {
      return url;
    }
    return 'https://' + url;
  }
}
