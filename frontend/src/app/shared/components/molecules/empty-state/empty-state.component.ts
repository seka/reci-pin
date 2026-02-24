import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-empty-state',
  standalone: true,
  imports: [CommonModule, MatIconModule],
  template: `
    <div class="empty-state">
      @if (icon) {
        <div class="icon-container">
          <mat-icon class="empty-icon">{{ icon }}</mat-icon>
        </div>
      }
      @if (title) {
        <h3 class="empty-title">{{ title }}</h3>
      }
      @if (message) {
        <p class="empty-message">{{ message }}</p>
      }
      <div class="empty-actions">
        <ng-content></ng-content>
      </div>
    </div>
  `,
  styles: [
    `
      .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        text-align: center;
        padding: var(--spacing-4) var(--spacing-2);
        color: var(--color-text-secondary);
        width: 100%;
        height: 100%;
        min-height: 200px;
      }
      .icon-container {
        margin-bottom: var(--spacing-2);
        color: var(--color-border);
      }
      .empty-icon {
        font-size: 64px;
        width: 64px;
        height: 64px;
      }
      .empty-title {
        font-size: var(--font-size-lg);
        font-weight: 500;
        margin: 0 0 var(--spacing-1) 0;
        color: var(--color-text);
      }
      .empty-message {
        font-size: var(--font-size-md);
        margin: 0 0 var(--spacing-3) 0;
        max-width: 400px;
        line-height: 1.5;
      }
      .empty-actions {
        display: flex;
        justify-content: center;
        gap: var(--spacing-2);
      }
    `,
  ],
})
export class EmptyStateComponent {
  @Input() icon?: string;
  @Input() title?: string;
  @Input() message?: string;
}
