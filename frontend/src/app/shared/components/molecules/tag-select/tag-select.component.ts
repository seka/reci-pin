import { Component, Input, Output, EventEmitter } from '@angular/core';
import { Tag } from '../../../../core/services/recipe.service';

@Component({
  selector: 'app-tag-select',
  standalone: true,
  imports: [],
  template: `
    <div class="tags-section">
      <span class="section-label" id="tag-select-label">タグ選択</span>
      <p class="helper-text">このレシピに関連するタグを選択してください。</p>

      @if (tags.length > 0) {
        <div class="tags-container" role="group" aria-labelledby="tag-select-label">
          @for (tag of tags; track tag.id) {
            <div class="tag-item">
              <input
                type="checkbox"
                [id]="'tag-' + tag.id"
                [value]="tag.id"
                [checked]="isSelected(tag.id)"
                (change)="onChange($event, tag.id)"
              />
              <label [for]="'tag-' + tag.id">
                {{ tag.name }}
              </label>
            </div>
          }
        </div>
      } @else {
        <p class="no-tags-message">登録されているタグがありません。</p>
      }
    </div>
  `,
  styles: [
    `
      .tags-section {
        margin-top: 16px;
        margin-bottom: 24px;
      }
      .section-label {
        font-size: 1rem;
        font-weight: 500;
        color: rgba(0, 0, 0, 0.6);
        display: block;
        margin-bottom: 8px;
      }
      .helper-text {
        font-size: 0.85em;
        color: #666;
        margin-bottom: 12px;
        margin-top: 0;
      }

      .tags-container {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
      }

      .tag-item {
        position: relative;
      }
      .tag-item input[type='checkbox'] {
        position: absolute;
        opacity: 0;
        width: 0;
        height: 0;
      }
      .tag-item label {
        display: inline-flex;
        align-items: center;
        padding: 6px 16px;
        border: 1px solid #ddd;
        border-radius: 999px; /* Pill shape */
        background: #fff;
        cursor: pointer;
        font-weight: 500;
        transition: all 0.2s ease;
        color: #555;
        user-select: none;
        font-size: 0.9rem;
      }
      .tag-item label:hover {
        background: #f5f5f5;
      }
      .tag-item input[type='checkbox']:checked + label {
        background: #e0f7fa; /* Light Cyan */
        border-color: #00bcd4;
        color: #006064;
        font-weight: bold;
      }
      .tag-item input[type='checkbox']:checked + label::before {
        content: '✓';
        margin-right: 6px;
        font-weight: bold;
      }
      .no-tags-message {
        color: #888;
        font-style: italic;
      }
    `,
  ],
})
export class TagSelectComponent {
  @Input() tags: Tag[] = [];
  @Input() selectedTagIds: number[] = [];
  @Output() selectionChange = new EventEmitter<number[]>();

  isSelected(tagId: number): boolean {
    return this.selectedTagIds.includes(tagId);
  }

  onChange(event: Event, tagId: number) {
    const target = event.target as HTMLInputElement;
    let newSelection = [...this.selectedTagIds];
    if (target.checked) {
      newSelection.push(tagId);
    } else {
      newSelection = newSelection.filter((id) => id !== tagId);
    }
    this.selectionChange.emit(newSelection);
  }
}
