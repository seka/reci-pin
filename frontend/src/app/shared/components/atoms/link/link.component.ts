import { Component, Input, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-link',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './link.component.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './link.component.scss',
})
export class LinkComponent {
  @Input() routerLink?: string | (string | number)[];
  @Input() href?: string;
  @Input() title?: string;
  @Input() variant: 'primary' | 'secondary' = 'primary';
}
