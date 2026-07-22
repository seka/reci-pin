import {
  Component,
  EventEmitter,
  inject,
  Input,
  OnChanges,
  OnInit,
  Output,
  SimpleChanges,
  signal,
} from '@angular/core';
import { form, FormField, FormRoot, maxLength, required } from '@angular/forms/signals';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { TranslocoPipe, TranslocoService } from '@jsverse/transloco';
import { RecipeService } from '../../../../core/services/recipe.service';
import { Tag, RecipeFormModel } from '../../../../core/models/recipe.model';
import { TagSelectComponent } from '../../molecules/tag-select/tag-select.component';
import { ButtonComponent } from '../../atoms/button/button.component';
import { HeadlineComponent } from '../../atoms/headline/headline.component';
import { InputComponent } from '../../atoms/input/input.component';
import { TextareaComponent } from '../../atoms/textarea/textarea.component';
import { VALIDATION_RULES } from '../../../../core/constants/validation.constants';
import { ApiError } from '../../../../core/models/api-error.model';

export interface RecipeFormSubmitEvent {
  formData: RecipeFormModel;
  file: File | null;
}

@Component({
  selector: 'app-recipe-form',
  imports: [
    CommonModule,
    FormField,
    FormRoot,
    RouterModule,
    MatCardModule,
    MatIconModule,
    TranslocoPipe,
    TagSelectComponent,
    ButtonComponent,
    HeadlineComponent,
    InputComponent,
    TextareaComponent,
  ],
  templateUrl: './recipe-form.component.html',
  styleUrl: './recipe-form.component.scss',
})
export class RecipeFormComponent implements OnInit, OnChanges {
  private readonly recipeService = inject(RecipeService);
  private readonly translate = inject(TranslocoService);

  @Input() titleKey: string = 'FEATURES.RECIPES.RECIPE_CREATE.TITLE';
  @Input() submitLabelKey: string = 'COMPONENTS.ORGANISMS.RECIPE_FORM.SAVE';
  @Input() submittingLabelKey: string = 'COMPONENTS.ORGANISMS.RECIPE_FORM.SAVING';

  @Input() isSubmitting = false;
  @Input() initialData: Partial<RecipeFormModel> = {};
  @Input() initialImagePreview: string | null = null;

  @Output() save = new EventEmitter<RecipeFormSubmitEvent>();

  tags: Tag[] = [];
  fieldErrors: Record<string, string[]> = {};
  selectedFile: File | null = null;
  imagePreview: string | null = null;
  imageError: string | null = null;
  isDragover = false;

  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  private readonly model = signal<RecipeFormModel>({ name: '', url: '', memo: '', tagIds: [] });

  protected readonly recipeForm = form(this.model, (path) => {
    required(path.name);
    maxLength(path.name, VALIDATION_RULES.RECIPE.NAME_MAX_LENGTH);

    maxLength(path.url, VALIDATION_RULES.RECIPE.URL_MAX_LENGTH);

    maxLength(path.memo, VALIDATION_RULES.RECIPE.MEMO_MAX_LENGTH);
  });

  ngOnChanges(changes: SimpleChanges) {
    if (changes['initialData'] && this.initialData) {
      this.model.set({
        name: this.initialData.name || '',
        url: this.initialData.url || '',
        memo: this.initialData.memo || '',
        tagIds: this.initialData.tagIds || [],
      });
    }
    if (changes['initialImagePreview'] && this.initialImagePreview) {
      this.imagePreview = this.initialImagePreview;
    }
  }

  onUrlBlur() {
    const urlField = this.recipeForm.url;
    const currentValue = urlField().value();
    if (currentValue) {
      let url = currentValue.trim();
      if (url && !url.includes('://')) {
        url = 'https://' + url;
        urlField().value.set(url);
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
    if (this.recipeForm().valid()) {
      this.save.emit({
        formData: this.recipeForm().value(),
        file: this.selectedFile,
      });
    }
  }

  handleServerErrors(err: { error?: ApiError }): boolean {
    let hasValidationErrors = false;
    if (err.error?.error?.details) {
      const details = err.error.error.details;
      Object.keys(details).forEach((field) => {
        hasValidationErrors = true;
        const fieldDetails = details[field];
        const messages = fieldDetails.map((d) => {
          switch (d.code) {
            case 'REQUIRED':
              return this.translate.translate('VALIDATION.REQUIRED');
            case 'TEXT_TOO_LONG':
              return this.translate.translate('VALIDATION.MAX_LENGTH', { max: d.params?.['max'] });
            case 'URL_INVALID_FORMAT':
              return this.translate.translate('VALIDATION.INVALID_URL');
            default:
              return this.translate.translate('VALIDATION.INVALID_INPUT');
          }
        });
        this.fieldErrors[field] = messages;
      });
    }
    return hasValidationErrors;
  }
}
