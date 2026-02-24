import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-icon',
  standalone: true,
  imports: [CommonModule, MatIconModule],
  templateUrl: './icon.component.html',
  styleUrl: './icon.component.scss',
})
export class IconComponent {
  @Input() size: 'sm' | 'md' | 'lg' = 'md';
  @Input() color: 'inherit' | 'primary' | 'secondary' | 'warn' = 'inherit';
}
