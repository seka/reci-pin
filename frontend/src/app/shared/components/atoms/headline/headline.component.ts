import { Component, Input } from '@angular/core';
import { NgTemplateOutlet } from '@angular/common';

@Component({
  selector: 'app-headline',
  standalone: true,
  imports: [NgTemplateOutlet],
  template: `
    <ng-template #content><ng-content></ng-content></ng-template>

    @switch (variant) {
      @case ('h1') {
        <h1 class="h1"><ng-container *ngTemplateOutlet="content"></ng-container></h1>
      }
      @case ('h2') {
        <h2 class="h2"><ng-container *ngTemplateOutlet="content"></ng-container></h2>
      }
      @case ('h3') {
        <h3 class="h3"><ng-container *ngTemplateOutlet="content"></ng-container></h3>
      }
      @case ('h4') {
        <h4 class="h4"><ng-container *ngTemplateOutlet="content"></ng-container></h4>
      }
      @case ('h5') {
        <h5 class="h5"><ng-container *ngTemplateOutlet="content"></ng-container></h5>
      }
      @case ('h6') {
        <h6 class="h6"><ng-container *ngTemplateOutlet="content"></ng-container></h6>
      }
    }
  `,
  styles: [
    `
      /* Margin reset and consistent color */
      h1,
      h2,
      h3,
      h4,
      h5,
      h6 {
        margin: 0 0 var(--spacing-2);
        font-weight: 700;
        line-height: 1.2;
      }

      /* Typography System - Pop Design */
      .h1 {
        font-size: 2.5rem;
        color: var(--color-text-main);
      }
      .h2 {
        font-size: 2rem;
        color: var(--color-primary); /* Main Title Pinkish */
      }
      .h3 {
        font-size: 1.75rem;
        color: var(--color-text-main);
      }
      .h4 {
        font-size: 1.5rem;
        color: var(--color-text-main);
      }
      .h5 {
        font-size: 1.25rem;
        color: var(--color-text-main);
      }
      .h6 {
        font-size: 1rem;
        color: var(--color-text-main);
        text-transform: uppercase;
        letter-spacing: 0.05em;
      }
    `,
  ],
})
export class HeadlineComponent {
  @Input() variant: 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6' = 'h2';
}
