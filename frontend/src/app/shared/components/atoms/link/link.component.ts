import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-link',
  standalone: true,
  imports: [CommonModule, RouterModule],
  template: `
    @if (routerLink) {
      <a [routerLink]="routerLink" [ngClass]="['link', variant]" [attr.title]="title"><ng-content></ng-content></a>
    } @else {
      <a [href]="href" [ngClass]="['link', variant]" [attr.title]="title"><ng-content></ng-content></a>
    }
  `,
  styles: [
    `
      :host {
        display: inline;
      }
      .link {
        color: var(--color-primary);
        text-decoration: none;
        transition: text-decoration 0.2s, color 0.2s;
        cursor: pointer;
      }
      .link:hover {
        text-decoration: underline;
      }
      .link.secondary {
        color: var(--color-text-secondary);
        font-size: var(--font-size-2);
      }
      .link.secondary:hover {
        text-decoration: underline;
      }
    `,
  ],
})
export class LinkComponent {
  @Input() routerLink?: string | unknown[];
  @Input() href?: string;
  @Input() title?: string;
  @Input() variant: 'primary' | 'secondary' = 'primary';
}
