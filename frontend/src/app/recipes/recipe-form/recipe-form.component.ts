import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { RecipeService, Tag } from '../../core/services/recipe.service';

@Component({
  selector: 'app-recipe-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  template: `
    <div class="recipe-form-container">
      <h2>新規レシピ作成</h2>
      <form [formGroup]="recipeForm" (ngSubmit)="onSubmit()">
        <div class="form-group">
          <label for="name">レシピ名</label>
          <input id="name" type="text" formControlName="name" placeholder="レシピ名を入力">
          <div *ngIf="recipeForm.get('name')?.invalid && recipeForm.get('name')?.touched" class="error">
            レシピ名は必須です
          </div>
        </div>

        <div class="form-group">
          <label for="url">URL</label>
          <input id="url" type="text" formControlName="url" placeholder="https://example.com/recipe">
          <div *ngIf="recipeForm.get('url')?.invalid && recipeForm.get('url')?.touched" class="error">
            有効なURLを入力してください
          </div>
        </div>

        <div class="form-group">
          <label for="memo">メモ</label>
          <textarea id="memo" formControlName="memo" rows="4" placeholder="メモを入力"></textarea>
        </div>

        <div class="form-group">
          <label class="section-label">タグ選択</label>
          <p class="helper-text">このレシピに関連するタグを選択してください。</p>
          
          <div class="tags-container" *ngIf="tags.length > 0; else noTags">
             <div *ngFor="let tag of tags" class="tag-item">
               <input type="checkbox" [id]="'tag-' + tag.id" [value]="tag.id" (change)="onTagChange($event, tag.id)">
               <label [for]="'tag-' + tag.id">
                 {{ tag.name }}
               </label>
             </div>
          </div>
          <ng-template #noTags>
            <p class="no-tags-message">登録されているタグがありません。</p>
          </ng-template>
        </div>

        <div class="actions">
          <button type="button" routerLink="/recipes" class="cancel-btn">キャンセル</button>
          <button type="submit" [disabled]="recipeForm.invalid || isSubmitting" class="submit-btn">
            {{ isSubmitting ? '保存中...' : '保存' }}
          </button>
        </div>
      </form>
    </div>
  `,
  styles: [`
    .recipe-form-container { padding: 20px; max-width: 600px; margin: 0 auto; }
    .form-group { margin-bottom: 15px; }
    label { display: block; margin-bottom: 5px; font-weight: bold; }
    input[type="text"], textarea { width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
    .error { color: red; font-size: 12px; margin-top: 5px; }
    .actions { display: flex; gap: 10px; margin-top: 20px; }
    button { padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; }
    .submit-btn { background: #007bff; color: white; }
    .submit-btn:disabled { background: #ccc; }
    .cancel-btn { background: #f0f0f0; }
    .section-label { font-size: 1.1em; margin-bottom: 4px; }
    .helper-text { font-size: 0.9em; color: #666; margin-bottom: 10px; margin-top: 0; }
    .tags-container { 
      display: flex; 
      flex-wrap: wrap; 
      gap: 10px; 
      padding: 15px; 
      background: #ffffff; 
      border: 1px solid #eee;
      border-radius: 8px; 
    }
    .tag-item {
      position: relative;
    }
    .tag-item input[type="checkbox"] {
      position: absolute;
      opacity: 0;
      width: 0;
      height: 0;
    }
    .tag-item label { 
      display: inline-flex; 
      align-items: center; 
      padding: 8px 16px;
      border: 1px solid #ddd;
      border-radius: 20px;
      background: #f8f9fa;
      cursor: pointer;
      font-weight: 500;
      transition: all 0.2s ease;
      color: #333;
      user-select: none;
    }
    .tag-item label:hover {
      background: #e9ecef;
      border-color: #adb5bd;
    }
    .tag-item input[type="checkbox"]:checked + label {
      background: #e7f5ff;
      border-color: #007bff;
      color: #007bff;
      font-weight: bold;
      box-shadow: 0 2px 4px rgba(0,123,255,0.1);
    }
    /* Add a checkmark icon when selected */
    .tag-item input[type="checkbox"]:checked + label::before {
      content: '✓';
      margin-right: 6px;
      font-weight: bold;
    }
    .no-tags-message { color: #888; font-style: italic; }
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

  ngOnInit() {
    this.recipeService.getAllTags().subscribe({
      next: (tags) => this.tags = tags,
      error: (err) => console.error('Failed to load tags', err)
    });
  }

  onTagChange(event: any, tagId: number) {
    if (event.target.checked) {
      this.selectedTagIds.push(tagId);
    } else {
      this.selectedTagIds = this.selectedTagIds.filter(id => id !== tagId);
    }
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
