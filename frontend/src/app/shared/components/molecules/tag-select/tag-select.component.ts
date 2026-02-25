import { Component, Input, forwardRef, inject, ElementRef, ViewChild, OnInit } from '@angular/core';
import {
  ControlValueAccessor,
  NG_VALUE_ACCESSOR,
  FormControl,
  ReactiveFormsModule,
} from '@angular/forms';
import { Tag, RecipeService } from '../../../../core/services/recipe.service';
import { CommonModule } from '@angular/common';
import { MatChipsModule, MatChipInputEvent } from '@angular/material/chips';
import {
  MatAutocompleteModule,
  MatAutocompleteSelectedEvent,
} from '@angular/material/autocomplete';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { COMMA, ENTER } from '@angular/cdk/keycodes';
import { Observable, startWith, map } from 'rxjs';
import { TranslocoPipe } from '@jsverse/transloco';

@Component({
  selector: 'app-tag-select',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatChipsModule,
    MatAutocompleteModule,
    MatIconModule,
    MatFormFieldModule,
    TranslocoPipe,
  ],
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => TagSelectComponent),
      multi: true,
    },
  ],
  templateUrl: './tag-select.component.html',
  styleUrl: './tag-select.component.scss',
})
export class TagSelectComponent implements ControlValueAccessor, OnInit {
  @Input() tags: Tag[] = [];

  separatorKeysCodes: number[] = [ENTER, COMMA];
  tagCtrl = new FormControl('');
  filteredTags: Observable<Tag[]>;
  selectedTagIds: number[] = [];
  isDisabled = false;

  @ViewChild('tagInput') tagInput!: ElementRef<HTMLInputElement>;

  private readonly recipeService = inject(RecipeService);

  // eslint-disable-next-line @typescript-eslint/no-empty-function
  onChange: (value: number[]) => void = () => {};
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  onTouched: () => void = () => {};

  constructor() {
    this.filteredTags = this.tagCtrl.valueChanges.pipe(
      startWith(null),
      map((tag: string | null) => (tag ? this.filter(tag) : this.getUnselectedTags())),
    );
  }

  ngOnInit() {
    // Ensure initial filtering works correctly
    this.tagCtrl.setValue(null);
  }

  getTagName(id: number): string {
    return this.tags.find((t) => t.id === id)?.name || '';
  }

  private getUnselectedTags(): Tag[] {
    return this.tags.filter((tag) => !this.selectedTagIds.includes(tag.id));
  }

  add(event: MatChipInputEvent): void {
    const value = (event.value || '').trim();

    if (value) {
      // Check if tag already exists in the full list
      const existingTag = this.tags.find((t) => t.name === value);

      if (existingTag) {
        if (!this.selectedTagIds.includes(existingTag.id)) {
          this.selectedTagIds = [...this.selectedTagIds, existingTag.id];
          this.triggerChange();
        }
      } else {
        // Create new tag
        this.recipeService.createTag(value).subscribe({
          next: (newTag) => {
            this.tags = [...this.tags, newTag];
            this.selectedTagIds = [...this.selectedTagIds, newTag.id];
            this.triggerChange();

            // Clear input manually since we handled the event thoroughly
            event.chipInput!.clear();
            this.tagCtrl.setValue(null);
          },
          error: (err) => console.error('Failed to create tag', err),
        });
        return; // Return early to avoid clearing input before async op completes (though we clear it in subscribe)
      }
    }

    // Clear the input value
    event.chipInput!.clear();
    this.tagCtrl.setValue(null);
  }

  remove(tagId: number): void {
    const index = this.selectedTagIds.indexOf(tagId);

    if (index >= 0) {
      this.selectedTagIds.splice(index, 1);
      this.triggerChange();
    }
  }

  selected(event: MatAutocompleteSelectedEvent): void {
    const tag: Tag = event.option.value;
    if (!this.selectedTagIds.includes(tag.id)) {
      this.selectedTagIds.push(tag.id);
      this.triggerChange();
    }
    this.tagInput.nativeElement.value = '';
    this.tagCtrl.setValue(null);
  }

  private filter(value: string | Tag | null): Tag[] {
    // If value is a Tag object (from autocomplete selection), use its name, otherwise use the string
    const filterValue = (typeof value === 'string' ? value : value?.name || '').toLowerCase();

    return this.tags.filter(
      (tag) =>
        tag.name.toLowerCase().includes(filterValue) && !this.selectedTagIds.includes(tag.id),
    );
  }

  triggerChange() {
    this.onChange(this.selectedTagIds);
    this.onTouched();
    // Force refresh filtered list
    this.tagCtrl.setValue(null);
  }

  // ControlValueAccessor implementation
  writeValue(value: number[]): void {
    this.selectedTagIds = value || [];
  }

  registerOnChange(fn: (value: number[]) => void): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: () => void): void {
    this.onTouched = fn;
  }

  setDisabledState?(isDisabled: boolean): void {
    this.isDisabled = isDisabled;
    if (isDisabled) {
      this.tagCtrl.disable();
    } else {
      this.tagCtrl.enable();
    }
  }
}
