import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-text',
  standalone: true,
  imports: [CommonModule],
  template: `
    <ng-template #content><ng-content></ng-content></ng-template>

    <ng-container [ngSwitch]="variant">
      <p *ngSwitchCase="'body'" class="body"><ng-container *ngTemplateOutlet="content"></ng-container></p>
      <span *ngSwitchCase="'caption'" class="caption"><ng-container *ngTemplateOutlet="content"></ng-container></span>
      <label *ngSwitchCase="'label'" class="label"><ng-container *ngTemplateOutlet="content"></ng-container></label>
    </ng-container>
  `,
  styles: [`
    /* Typography System */
    .body { font-size: 1rem; line-height: 1.6; color: #666; margin: 0 0 16px; }
    .caption { font-size: 0.875rem; color: #888; }
    .label { font-size: 1rem; font-weight: 700; color: #333; display: block; margin-bottom: 8px; }
  `]
})
export class TextComponent {
  @Input() variant: 'body' | 'caption' | 'label' = 'body';
}
