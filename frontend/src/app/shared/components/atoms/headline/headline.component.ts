import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-headline',
  standalone: true,
  imports: [CommonModule],
  template: `
    <ng-template #content><ng-content></ng-content></ng-template>
    
    <ng-container [ngSwitch]="variant">
      <h1 *ngSwitchCase="'h1'" class="h1"><ng-container *ngTemplateOutlet="content"></ng-container></h1>
      <h2 *ngSwitchCase="'h2'" class="h2"><ng-container *ngTemplateOutlet="content"></ng-container></h2>
      <h3 *ngSwitchCase="'h3'" class="h3"><ng-container *ngTemplateOutlet="content"></ng-container></h3>
      <h4 *ngSwitchCase="'h4'" class="h4"><ng-container *ngTemplateOutlet="content"></ng-container></h4>
      <h5 *ngSwitchCase="'h5'" class="h5"><ng-container *ngTemplateOutlet="content"></ng-container></h5>
      <h6 *ngSwitchCase="'h6'" class="h6"><ng-container *ngTemplateOutlet="content"></ng-container></h6>
    </ng-container>
  `,
  styles: [`
    /* Margin reset and consistent color */
    h1, h2, h3, h4, h5, h6 { margin: 0 0 16px; font-weight: 700; line-height: 1.2; }
    
    /* Typography System - Pop Design */
    .h1 { font-size: 2.5rem; color: #333; }
    .h2 { font-size: 2rem; color: #e91e63; /* Main Title Pinkish */ }
    .h3 { font-size: 1.75rem; color: #333; }
    .h4 { font-size: 1.5rem; color: #333; }
    .h5 { font-size: 1.25rem; color: #333; }
    .h6 { font-size: 1rem; color: #333; text-transform: uppercase; letter-spacing: 0.05em; }
  `]
})
export class HeadlineComponent {
  @Input() variant: 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6' = 'h2';
}
