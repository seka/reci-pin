import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-icon',
  standalone: true,
  imports: [CommonModule, MatIconModule],
  template: `
    <mat-icon [ngClass]="['icon-' + size, 'icon-' + color]"><ng-content></ng-content></mat-icon>
  `,
  styles: [
    `
      :host {
        display: inline-flex;
        align-items: center;
        vertical-align: middle;
      }

      // Sizes handled by font-size usually, Material Icons font adjustment
      .icon-sm {
        font-size: 18px;
        width: 18px;
        height: 18px;
        line-height: 18px;
      }
      .icon-md {
        font-size: 24px;
        width: 24px;
        height: 24px;
        line-height: 24px;
      } // Default
      .icon-lg {
        font-size: 36px;
        width: 36px;
        height: 36px;
        line-height: 36px;
      }

      // Colors - mapped to global tokens via classes or CSS variables if needed.
      // However, material icons usually inherit color.
      // Let's rely on standard color inheritance or helper classes.
      // For now, keep it simple.
    `,
  ],
})
export class IconComponent {
  @Input() size: 'sm' | 'md' | 'lg' = 'md';
  @Input() color: 'inherit' | 'primary' | 'secondary' | 'warn' = 'inherit';
}
