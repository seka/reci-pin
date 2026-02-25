import { Component, Input, Output, EventEmitter, inject } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { TranslocoPipe } from '@jsverse/transloco';
import { RouterModule } from '@angular/router';
import { MatDialog } from '@angular/material/dialog';
import { Recipe } from '../../../../core/services/recipe.service';
import { VALIDATION_RULES } from '../../../../core/constants/validation.constants';
import { ConfirmDialogComponent } from '../../molecules/confirm-dialog/confirm-dialog.component';

@Component({
  selector: 'app-recipe-card',
  standalone: true,
  imports: [
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatChipsModule,
    TranslocoPipe,
    RouterModule,
  ],
  templateUrl: './recipe-card.component.html',
  styleUrl: './recipe-card.component.scss',
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
        message: 'COMPONENTS.ORGANISMS.RECIPE_CARD.DELETE_CONFIRMATION',
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
    } catch {
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
