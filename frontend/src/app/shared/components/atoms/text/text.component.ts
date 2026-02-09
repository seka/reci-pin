import { Component, Input } from '@angular/core';
import { NgTemplateOutlet } from '@angular/common';

@Component({
  selector: 'app-text',
  standalone: true,
  imports: [NgTemplateOutlet],
  template: `
    <ng-template #content><ng-content></ng-content></ng-template>

    @switch (variant) {
      @case ('body') {
        <p class="body"><ng-container *ngTemplateOutlet="content"></ng-container></p>
      }
      @case ('caption') {
        <span class="caption"><ng-container *ngTemplateOutlet="content"></ng-container></span>
      }
      @case ('label') {
        <span class="label"><ng-container *ngTemplateOutlet="content"></ng-container></span>
      }
    }
  `,
  styles: [
    `
      /* Typography System */
      .body {
        font-size: 1rem;
        line-height: 1.6;
        color: #666;
        margin: 0 0 16px;
      }
      .caption {
        font-size: 0.875rem;
        color: #888;
      }
      .label {
        font-size: 1rem;
        font-weight: 700;
        color: #333;
        display: block;
        margin-bottom: 8px;
      }
    `,
  ],
})
export class TextComponent {
  @Input() variant: 'body' | 'caption' | 'label' = 'body';
}
