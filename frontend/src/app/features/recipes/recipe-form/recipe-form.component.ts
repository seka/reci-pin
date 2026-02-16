import { Component, inject, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
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
    ReactiveFormsModule,
    RouterModule,
    MatCardModule,
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
              [selectedTagIds]="selectedTagIds"
              (selectionChange)="selectedTagIds = $event"
            >
            </app-tag-select>

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
  selectedTagIds: number[] = [];
  fieldErrors: { [key: string]: string[] } = {};
  isSubmitting = false;

  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  constructor() {
    this.recipeForm = this.fb.group({
      name: ['', [Validators.required, Validators.maxLength(VALIDATION_RULES.RECIPE.NAME_MAX_LENGTH)]],
      url: ['', [Validators.required, Validators.pattern(/^https?:\/\/[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}.*|^https?:\/\/localhost.*/)]],
      memo: ['', Validators.maxLength(VALIDATION_RULES.RECIPE.MEMO_MAX_LENGTH)],
    });
  }

  onUrlBlur() {
    const urlControl = this.recipeForm.get('url');
    if (urlControl?.value) {
      let url = urlControl.value.trim();
      // If no scheme is present at all (no '://'), prepend https://
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

  onSubmit() {
    this.fieldErrors = {};
    if (this.recipeForm.valid) {
      this.isSubmitting = true;
      const formData = {
        ...this.recipeForm.value,
        tag_ids: this.selectedTagIds,
      };

      this.recipeService.createRecipe(formData).subscribe({
        next: () => {
          this.router.navigate(['/recipes']);
        },
        error: (err) => {
          console.error('Failed to create recipe', err);
          this.isSubmitting = false;

          if (err.error && err.error.error && err.error.error.details) {
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
