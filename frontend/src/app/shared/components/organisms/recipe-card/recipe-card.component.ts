import { Component, Input } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { TranslatePipe } from '@ngx-translate/core';
import { Recipe } from '../../../../core/services/recipe.service';
import { VALIDATION_RULES } from '../../../../core/constants/validation.constants';

@Component({
  selector: 'app-recipe-card',
  standalone: true,
  imports: [MatCardModule, MatButtonModule, MatIconModule, MatChipsModule, TranslatePipe],
  template: `
    <mat-card class="recipe-card">
      @if (thumbnailUrl) {
        <div class="card-image">
          <img [src]="thumbnailUrl" [alt]="recipe.name" />
        </div>
      }
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
            {{ 'RECIPE.VIEW_DETAILS' | translate }}
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
      .card-image {
        width: 100%;
        height: 180px;
        overflow: hidden;
        border-radius: var(--radius-2) var(--radius-2) 0 0;
      }
      .card-image img {
        width: 100%;
        height: 100%;
        object-fit: cover;
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

  get thumbnailUrl(): string | null {
    if (!this.recipe.images?.length) return null;

    const imageUrl = this.recipe.images[0].image_url;
    if (!imageUrl) return null;

    try {
      const url = new URL(imageUrl, window.location.origin);
      const allowedDomains = [...VALIDATION_RULES.IMAGE.ALLOWED_DOMAINS, window.location.hostname];

      if (allowedDomains.includes(url.hostname)) {
        return imageUrl;
      }
      console.warn('Blocked image from untrusted domain:', url.hostname);
      return null;
    } catch (e) {
      console.error('Invalid image URL:', imageUrl);
      return null;
    }
  }

  getExternalUrl(url: string): string {
    if (!url) return '';
    if (/^https?:\/\//i.test(url)) {
      return url;
    }
    return 'https://' + url;
  }
}

