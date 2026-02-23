import { Component, Input, Output, EventEmitter, inject } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { TranslatePipe } from '@ngx-translate/core';
import { RouterModule } from '@angular/router';
import { MatDialog } from '@angular/material/dialog';
import { Recipe } from '../../../../core/services/recipe.service';
import { VALIDATION_RULES } from '../../../../core/constants/validation.constants';
import { ConfirmDialogComponent } from '../../molecules/confirm-dialog/confirm-dialog.component';

@Component({
  selector: 'app-recipe-card',
  standalone: true,
  imports: [MatCardModule, MatButtonModule, MatIconModule, MatChipsModule, TranslatePipe, RouterModule],
  template: `
    <mat-card class="recipe-card">
      <div class="action-buttons">
        <a
          mat-icon-button
          class="action-btn"
          [routerLink]="['/recipes', recipe.id, 'edit']"
          [attr.aria-label]="'COMMON.EDIT' | translate"
        >
          <mat-icon>edit</mat-icon>
        </a>
        <button
          mat-icon-button
          color="warn"
          class="action-btn"
          (click)="openDeleteDialog()"
          [attr.aria-label]="'COMMON.DELETE' | translate"
        >
          <mat-icon>delete</mat-icon>
        </button>
      </div>
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
      .action-buttons {
        position: absolute;
        top: var(--spacing-1);
        right: var(--spacing-1);
        z-index: 10;
        display: flex;
        gap: var(--spacing-1);
      }
      .action-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        background: rgba(255, 255, 255, 0.7);
        backdrop-filter: blur(4px);
        color: var(--color-primary);
        transition: all 0.2s ease;
      }
      .action-btn[color="warn"] {
        color: var(--color-error);
      }
      .action-btn:hover {
        background: rgba(255, 255, 255, 0.9);
        transform: scale(1.1);
      }
    `,
  ],
})
export class RecipeCardComponent {
  @Input() recipe!: Recipe;
  @Output() delete = new EventEmitter<number>();

  private readonly dialog = inject(MatDialog);

  openDeleteDialog(): void {
    const dialogRef = this.dialog.open(ConfirmDialogComponent, {
      width: '400px',
      data: {
        title: 'COMMON.CONFIRM_DELETE',
        message: 'RECIPE.DELETE_CONFIRMATION',
      },
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        this.delete.emit(this.recipe.id);
      }
    });
  }

  get thumbnailUrl(): string | null {
    if (!this.recipe.images?.length) return null;

    const imageUrl = this.recipe.images[0].imageUrl;
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

