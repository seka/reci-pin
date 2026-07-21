import { Component, Input, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-icon',
  standalone: true,
  imports: [CommonModule, MatIconModule],
  templateUrl: './icon.component.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './icon.component.scss',
})
export class IconComponent {
  @Input() size: 'sm' | 'md' | 'lg' = 'md';
  @Input() color: 'inherit' | 'primary' | 'secondary' | 'warn' = 'inherit';
}
