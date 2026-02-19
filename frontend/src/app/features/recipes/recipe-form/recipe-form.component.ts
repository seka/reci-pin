import { Component, inject, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { RecipeService, Tag } from '../../../core/services/recipe.service';
import { TagSelectComponent } from '../../../shared/components/molecules/tag-select/tag-select.component';
import { ButtonComponent } from '../../../shared/components/atoms/button/button.component';
import { HeadlineComponent } from '../../../shared/components/atoms/headline/headline.component';
import { InputComponent } from '../../../shared/components/atoms/input/input.component';
import { TextareaComponent } from '../../../shared/components/atoms/textarea/textarea.component';
import { VALIDATION_RULES } from '../../../core/constants/validation.constants';

@Component({
  selector: 'app-recipe-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    MatCardModule,
    MatIconModule,
    TranslatePipe,
    TagSelectComponent,
    ButtonComponent,
    HeadlineComponent,
    InputComponent,
    TextareaComponent,
  ],
  template: `
    <div class="form-container">
      <mat-card>
        <mat-card-header>
          <mat-card-title>
            <app-headline variant="h2">{{ 'RECIPE.NEW_TITLE' | translate }}</app-headline>
          </mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <form [formGroup]="recipeForm" (ngSubmit)="onSubmit()">
            <div style="margin-bottom: var(--spacing-2);">
              <app-input
                [label]="'RECIPE.NAME' | translate"
                formControlName="name"
                [placeholder]="'RECIPE.NAME_PLACEHOLDER' | translate"
                [required]="true"
                [maxLength]="VALIDATION_RULES.RECIPE.NAME_MAX_LENGTH"
                [showCounter]="true"
                [errorMessage]="fieldErrors['name']"
              ></app-input>
            </div>

            <div style="margin-bottom: var(--spacing-2);">
              <app-input
                label="URL"
                formControlName="url"
                placeholder="https://example.com/recipe"
                [required]="true"
                (blur)="onUrlBlur()"
                [errorMessage]="fieldErrors['url']"
              ></app-input>
            </div>

            <div style="margin-bottom: var(--spacing-2);">
              <app-textarea
                [label]="'RECIPE.MEMO' | translate"
                formControlName="memo"
                [placeholder]="'RECIPE.MEMO_PLACEHOLDER' | translate"
                [rows]="4"
                [maxLength]="VALIDATION_RULES.RECIPE.MEMO_MAX_LENGTH"
                [showCounter]="true"
                [errorMessage]="fieldErrors['memo']"
              ></app-textarea>
            </div>

            <app-tag-select
              [tags]="tags"
              formControlName="tag_ids"
            >
            </app-tag-select>

            <!-- Image Upload -->
            <div class="image-upload-section">
              <label class="upload-label">{{ 'RECIPE.IMAGE' | translate }}</label>
              <div
                class="dropzone"
                [class.dragover]="isDragover"
                (dragover)="onDragOver($event)"
                (dragleave)="onDragLeave($event)"
                (drop)="onDrop($event)"
                (click)="fileInput.click()"
              >
                @if (imagePreview) {
                  <img
                    [src]="imagePreview"
                    alt="Preview"
                    class="preview-image"
                  />
                  <button
                    type="button"
                    class="remove-btn"
                    (click)="removeImage($event)"
                  >
                    <mat-icon>close</mat-icon>
                  </button>
                } @else {
                  <mat-icon class="upload-icon">cloud_upload</mat-icon>
                  <p class="upload-text">{{ 'RECIPE.IMAGE_DROP_HINT' | translate }}</p>
                  <p class="upload-subtext">JPEG, PNG, WebP (最大 50MB)</p>
                }
              </div>
              <input
                #fileInput
                type="file"
                [accept]="VALIDATION_RULES.IMAGE.ACCEPT"
                (change)="onFileSelected($event)"
                hidden
              />
              @if (imageError) {
                <p class="image-error">{{ imageError }}</p>
              }
            </div>

            <div class="actions">
              <app-button type="button" routerLink="/recipes" variant="warn" class="action-btn"
                >{{ 'COMMON.CANCEL' | translate }}</app-button
              >
              <app-button
                type="submit"
                variant="primary"
                [disabled]="recipeForm.invalid || isSubmitting"
                class="action-btn"
              >
                {{ isSubmitting ? ('RECIPE.SAVING' | translate) : ('RECIPE.SAVE' | translate) }}
              </app-button>
            </div>
          </form>
        </mat-card-content>
      </mat-card>
    </div>
  `,
  styles: [
    `
      .recipe-form-container {
        padding: var(--spacing-3);
        max-width: 600px;
        margin: 0 auto;
      }
      mat-card {
        padding: var(--spacing-3);
      }
      mat-card-title {
        margin-bottom: var(--spacing-3);
      }
      .full-width {
        width: 100%;
        margin-bottom: var(--spacing-1);
      }
      .actions {
        display: flex;
        justify-content: flex-end;
        gap: var(--spacing-2);
        margin-top: var(--spacing-3);
      }
      .action-btn {
        width: auto;
        min-width: 100px;
      }
      .image-upload-section {
        margin-top: var(--spacing-2);
        margin-bottom: var(--spacing-2);
      }
      .upload-label {
        display: block;
        font-size: var(--font-size-sm);
        color: var(--color-text-secondary);
        margin-bottom: var(--spacing-1);
      }
      .dropzone {
        position: relative;
        border: 2px dashed var(--color-border);
        border-radius: var(--radius-2);
        padding: var(--spacing-3);
        text-align: center;
        cursor: pointer;
        transition: border-color 0.2s, background-color 0.2s;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        min-height: 160px;
      }
      .dropzone:hover,
      .dropzone.dragover {
        border-color: var(--color-primary);
        background-color: rgba(var(--color-primary-rgb, 99, 102, 241), 0.05);
      }
      .upload-icon {
        font-size: 40px;
        width: 40px;
        height: 40px;
        color: var(--color-text-secondary);
        margin-bottom: var(--spacing-1);
      }
      .upload-text {
        color: var(--color-text-secondary);
        margin: 0;
        font-size: var(--font-size-sm);
      }
      .upload-subtext {
        color: var(--color-text-tertiary, #999);
        margin: var(--spacing-1) 0 0;
        font-size: var(--font-size-xs, 12px);
      }
      .preview-image {
        max-width: 100%;
        max-height: 200px;
        border-radius: var(--radius-1);
        object-fit: contain;
      }
      .remove-btn {
        position: absolute;
        top: 8px;
        right: 8px;
        background: rgba(0, 0, 0, 0.5);
        border: none;
        border-radius: 50%;
        color: white;
        width: 32px;
        height: 32px;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        padding: 0;
        transition: background 0.2s;
      }
      .remove-btn:hover {
        background: rgba(0, 0, 0, 0.7);
      }
      .remove-btn mat-icon {
        font-size: 18px;
        width: 18px;
        height: 18px;
      }
      .image-error {
        color: var(--color-error, #f44336);
        font-size: var(--font-size-xs, 12px);
        margin-top: var(--spacing-1);
      }
    `,
  ],
})
export class RecipeFormComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly recipeService = inject(RecipeService);
  private readonly router = inject(Router);
  private readonly translate = inject(TranslateService);

  recipeForm: FormGroup;
  tags: Tag[] = [];
  fieldErrors: Record<string, string[]> = {};
  isSubmitting = false;

  selectedFile: File | null = null;
  imagePreview: string | null = null;
  imageError: string | null = null;
  isDragover = false;

  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  constructor() {
    this.recipeForm = this.fb.group({
      name: ['', [Validators.required, Validators.maxLength(VALIDATION_RULES.RECIPE.NAME_MAX_LENGTH)]],
      url: ['', [Validators.required, Validators.pattern(/^https?:\/\/[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}.*|^https?:\/\/localhost.*/)]],
      memo: ['', Validators.maxLength(VALIDATION_RULES.RECIPE.MEMO_MAX_LENGTH)],
      tag_ids: [[]],
    });
  }

  onUrlBlur() {
    const urlControl = this.recipeForm.get('url');
    if (urlControl?.value) {
      let url = urlControl.value.trim();
      if (url && !url.includes('://')) {
        url = 'https://' + url;
        urlControl.setValue(url);
      }
    }
  }

  ngOnInit() {
    this.recipeService.getAllTags().subscribe({
      next: (tags) => (this.tags = tags),
      error: (err) => console.error('Failed to load tags', err),
    });
  }

  onFileSelected(event: Event) {
    const input = event.target as HTMLInputElement;
    if (input.files?.[0]) {
      this.handleFile(input.files[0]);
    }
  }

  onDragOver(event: DragEvent) {
    event.preventDefault();
    event.stopPropagation();
    this.isDragover = true;
  }

  onDragLeave(event: DragEvent) {
    event.preventDefault();
    event.stopPropagation();
    this.isDragover = false;
  }

  onDrop(event: DragEvent) {
    event.preventDefault();
    event.stopPropagation();
    this.isDragover = false;
    if (event.dataTransfer?.files?.[0]) {
      this.handleFile(event.dataTransfer.files[0]);
    }
  }

  removeImage(event: Event) {
    event.stopPropagation();
    this.selectedFile = null;
    this.imagePreview = null;
    this.imageError = null;
  }

  private handleFile(file: File) {
    this.imageError = null;

    if (!(VALIDATION_RULES.IMAGE.ALLOWED_TYPES as readonly string[]).includes(file.type)) {
      this.imageError = 'JPEG, PNG, WebP のみアップロードできます';
      return;
    }

    if (file.size > VALIDATION_RULES.IMAGE.MAX_FILE_SIZE) {
      this.imageError = 'ファイルサイズは50MB以下にしてください';
      return;
    }

    this.selectedFile = file;

    const reader = new FileReader();
    reader.onload = (e) => {
      this.imagePreview = e.target?.result as string;
    };
    reader.readAsDataURL(file);
  }

  onSubmit() {
    this.fieldErrors = {};
    if (this.recipeForm.valid) {
      this.isSubmitting = true;
      const formData = this.recipeForm.value;

      this.recipeService.createRecipe(formData).subscribe({
        next: (recipe) => {
          if (this.selectedFile) {
            this.recipeService.uploadImage(recipe.id, this.selectedFile).subscribe({
              next: () => this.router.navigate(['/recipes']),
              error: (err) => {
                console.error('Failed to upload image', err);
                // レシピは作成済みなので一覧に戻す
                this.router.navigate(['/recipes']);
              },
            });
          } else {
            this.router.navigate(['/recipes']);
          }
        },
        error: (err) => {
          console.error('Failed to create recipe', err);
          this.isSubmitting = false;

          if (err.error?.error?.details) {
            const details = err.error.error.details;
            Object.keys(details).forEach((field) => {
              const messages = (details as any)[field].map((d: any) => {
                switch (d.code) {
                  case 'REQUIRED':
                    return this.translate.instant('VALIDATION.REQUIRED');
                  case 'TEXT_TOO_LONG':
                    return this.translate.instant('VALIDATION.MAX_LENGTH', { max: d.params?.max });
                  case 'URL_INVALID_FORMAT':
                    return this.translate.instant('VALIDATION.INVALID_URL');
                  default:
                    return this.translate.instant('VALIDATION.INVALID_INPUT');
                }
              });
              this.fieldErrors[field] = messages;
            });
          }
        },
      });
    }
  }
}

