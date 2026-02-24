import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-link',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './link.component.html',
  styleUrl: './link.component.scss',
})
export class LinkComponent {
  @Input() routerLink?: string | unknown[];
  @Input() href?: string;
  @Input() title?: string;
  @Input() variant: 'primary' | 'secondary' = 'primary';
}
