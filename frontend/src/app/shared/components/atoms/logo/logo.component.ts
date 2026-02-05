import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

@Component({
    selector: 'app-logo',
    standalone: true,
    imports: [CommonModule, RouterModule],
    template: `
    <a routerLink="/" class="logo" [ngClass]="size">
      <span class="icon">📍</span>
      <span class="text">ReciPin</span>
    </a>
  `,
    styles: [`
    .logo {
      display: inline-flex;
      align-items: center;
      text-decoration: none;
      font-weight: 700;
      color: #333;
      transition: opacity 0.2s;
    }
    .logo:hover { opacity: 0.8; }
    
    .icon { margin-right: 8px; }
    
    /* Sizes */
    .small { font-size: 1.25rem; }
    .medium { font-size: 1.5rem; }
    .large { font-size: 2.5rem; }
    
    .text { color: #e91e63; /* Pinkish/Terracotta based */ }
  `]
})
export class LogoComponent {
    @Input() size: 'small' | 'medium' | 'large' = 'medium';
}
