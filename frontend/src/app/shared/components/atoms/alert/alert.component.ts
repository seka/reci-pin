import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-alert',
  standalone: true,
  imports: [CommonModule],
  template: `
    @if (message) {
      <div [ngClass]="['alert', type]">
        <ng-content></ng-content>
        {{ message }}
      </div>
    }
  `,
  styles: [
    `
      :host {
        display: block;
        width: 100%;
      }
      .alert {
        margin-top: var(--spacing-2);
        padding: var(--spacing-2);
        border-radius: var(--radius-1);
        text-align: center;
        font-size: var(--font-size-2);
      }
      .alert.error {
        color: var(--color-error);
      }
      .alert.success {
        color: var(--color-success);
      }
      .alert.info {
        color: var(--color-primary);
      }
    `,
  ],
})
export class AlertComponent {
  @Input() type: 'error' | 'success' | 'info' = 'info';
  @Input() message?: string | null;
}
