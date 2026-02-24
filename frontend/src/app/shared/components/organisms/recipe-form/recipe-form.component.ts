import {
  Component,
  EventEmitter,
  inject,
  Input,
  OnChanges,
  OnInit,
  Output,
  SimpleChanges,
} from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { TranslatePipe, TranslateService } from '@ngx-translate/core';
import { RecipeService, Tag } from '../../../../core/services/recipe.service';
import { TagSelectComponent } from '../../molecules/tag-select/tag-select.component';
import { ButtonComponent } from '../../atoms/button/button.component';
import { HeadlineComponent } from '../../atoms/headline/headline.component';
import { InputComponent } from '../../atoms/input/input.component';
import { TextareaComponent } from '../../atoms/textarea/textarea.component';
import { VALIDATION_RULES } from '../../../../core/constants/validation.constants';

export interface RecipeFormData {
  name: string;
  url: string;
  memo: string;
  tagIds: number[];
}

export interface RecipeFormSubmitEvent {
  formData: RecipeFormData;
  file: File | null;
}

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
  templateUrl: './recipe-form.component.html',
  styleUrl: './recipe-form.component.scss',
})
export class RecipeFormComponent implements OnInit, OnChanges {
  private readonly fb = inject(FormBuilder);
  private readonly recipeService = inject(RecipeService);
  private readonly translate = inject(TranslateService);

  @Input() titleKey: string = 'RECIPE.NEW_TITLE';
  @Input() submitLabelKey: string = 'RECIPE.SAVE';
  @Input() submittingLabelKey: string = 'RECIPE.SAVING';

  @Input() isSubmitting = false;
  @Input() initialData: Partial<RecipeFormData> = {};
  @Input() initialImagePreview: string | null = null;

  @Output() save = new EventEmitter<RecipeFormSubmitEvent>();

  recipeForm: FormGroup;
  tags: Tag[] = [];
  fieldErrors: Record<string, string[]> = {};
  selectedFile: File | null = null;
  imagePreview: string | null = null;
  imageError: string | null = null;
  isDragover = false;

  protected readonly VALIDATION_RULES = VALIDATION_RULES;

  constructor() {
    this.recipeForm = this.fb.group({
      name: [
        '',
        [Validators.required, Validators.maxLength(VALIDATION_RULES.RECIPE.NAME_MAX_LENGTH)],
      ],
      url: [
        '',
        [
          Validators.required,
          Validators.pattern(/^https?:\/\/[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}.*|^https?:\/\/localhost.*/),
        ],
      ],
      memo: ['', Validators.maxLength(VALIDATION_RULES.RECIPE.MEMO_MAX_LENGTH)],
      tagIds: [[]],
    });
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['initialData'] && this.initialData) {
      this.recipeForm.patchValue({
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
      this.save.emit({
        formData: this.recipeForm.value as RecipeFormData,
        file: this.selectedFile,
      });
    }
  }

  handleServerErrors(err: any): boolean {
    let hasValidationErrors = false;
    if (err.error?.error?.details) {
      const details = err.error.error.details;
      Object.keys(details).forEach((field) => {
        hasValidationErrors = true;
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
    return hasValidationErrors;
  }
}
