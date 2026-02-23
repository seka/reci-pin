import { Component, EventEmitter, Output } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { ButtonComponent } from '../../atoms/button/button.component';

@Component({
  selector: 'app-search-bar',
  standalone: true,
  imports: [CommonModule, MatIconModule, ButtonComponent],
  template: `
    <div class="search-bar">
      <div class="search-input-wrapper">
        <ng-content></ng-content>
      </div>
      <app-button (click)="onSearch()" variant="secondary" class="search-btn">
        <mat-icon class="search-icon">search</mat-icon>
        検索
      </app-button>
    </div>
  `,
  styles: [
    `
      .search-bar {
        display: flex;
        gap: var(--spacing-2);
        align-items: center;
        width: 100%;
      }
      .search-input-wrapper {
        flex: 1;
        width: 100%;
      }
      .search-icon {
        font-size: 18px;
        width: 18px;
        height: 18px;
        vertical-align: middle;
        margin-right: 4px;
      }
      /* Prevent layout shift for inputs inside search bar */
      ::ng-deep .search-input-wrapper .mat-mdc-form-field-subscript-wrapper {
        display: none;
      }
    `,
  ],
})
export class SearchBarComponent {
  @Output() searchSubmit = new EventEmitter<void>();

  onSearch() {
    this.searchSubmit.emit();
  }
}
