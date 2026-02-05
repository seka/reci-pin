import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { RecipeService, Tag } from '../../../core/services/recipe.service';
import { TagSelectComponent } from '../../../shared/components/molecules/tag-select/tag-select.component';
import { ButtonComponent } from '../../../shared/components/atoms/button/button.component';
import { HeadlineComponent } from '../../../shared/components/atoms/headline/headline.component';
import { TextComponent } from '../../../shared/components/atoms/text/text.component';

@Component({
  selector: 'app-recipe-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    TagSelectComponent,
    ButtonComponent,
    HeadlineComponent,
    TextComponent
  ],
  template: `
    <div class="form-container">
      <mat-card>
        <mat-card-header>
          <mat-card-title>
            <app-headline variant="h2">新規レシピ作成</app-headline>
          </mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <form [formGroup]="recipeForm" (ngSubmit)="onSubmit()">
            
            <mat-form-field appearance="outline" class="full-width">
              <mat-label>レシピ名</mat-label>
              <input matInput type="text" formControlName="name" placeholder="例: オムライス">
              <mat-error *ngIf="recipeForm.get('name')?.invalid">レシピ名は必須です</mat-error>
            </mat-form-field>

            <mat-form-field appearance="outline" class="full-width">
              <mat-label>URL</mat-label>
              <input matInput type="text" formControlName="url" placeholder="https://example.com/recipe" (blur)="onUrlBlur()">
              <mat-error *ngIf="recipeForm.get('url')?.invalid">有効なURLを入力してください</mat-error>
            </mat-form-field>

            <mat-form-field appearance="outline" class="full-width">
              <mat-label>メモ</mat-label>
              <textarea matInput formControlName="memo" rows="4" placeholder="メモを入力"></textarea>
            </mat-form-field>

            <app-tag-select 
              [tags]="tags" 
              [selectedTagIds]="selectedTagIds"
              (selectionChange)="selectedTagIds = $event">
            </app-tag-select>

            <div class="actions">
              <app-button type="button" routerLink="/recipes" variant="warn" class="action-btn">キャンセル</app-button>
              <app-button type="submit" variant="primary" [disabled]="recipeForm.invalid || isSubmitting" class="action-btn">
                 {{ isSubmitting ? '保存中...' : '保存' }}
              </app-button>
            </div>
          </form>
        </mat-card-content>
      </mat-card>
    </div>
  `,
  styles: [`
    .recipe-form-container { padding: 24px; max-width: 600px; margin: 0 auto; }
    mat-card { padding: 24px; }
    mat-card-title { margin-bottom: 24px; }
    .full-width { width: 100%; margin-bottom: 8px; }
    .actions { display: flex; justify-content: flex-end; gap: 16px; margin-top: 24px; }
    .action-btn { width: auto; min-width: 100px; }
  `]
})
export class RecipeFormComponent implements OnInit {
  recipeForm: FormGroup;
  tags: Tag[] = [];
  selectedTagIds: number[] = [];
  isSubmitting = false;

  constructor(
    private fb: FormBuilder,
    private recipeService: RecipeService,
    private router: Router
  ) {
    this.recipeForm = this.fb.group({
      name: ['', Validators.required],
      url: ['', [Validators.required]],
      memo: ['']
    });
  }

  onUrlBlur() {
    const urlControl = this.recipeForm.get('url');
    if (urlControl?.value) {
      let url = urlControl.value.trim();
      if (url && !/^https?:\/\//i.test(url)) {
        url = 'https://' + url;
        urlControl.setValue(url);
      }
    }
  }

  ngOnInit() {
    this.recipeService.getAllTags().subscribe({
      next: (tags) => this.tags = tags,
      error: (err) => console.error('Failed to load tags', err)
    });
  }

  onSubmit() {
    if (this.recipeForm.valid) {
      this.isSubmitting = true;
      const formData = {
        ...this.recipeForm.value,
        tag_ids: this.selectedTagIds
      };

      this.recipeService.createRecipe(formData).subscribe({
        next: () => {
          this.router.navigate(['/recipes']);
        },
        error: (err) => {
          console.error('Failed to create recipe', err);
          this.isSubmitting = false;
        }
      });
    }
  }
}
