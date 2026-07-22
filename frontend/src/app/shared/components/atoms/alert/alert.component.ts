import { Component, Input, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-alert',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './alert.component.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './alert.component.scss',
})
export class AlertComponent {
  @Input() type: 'error' | 'success' | 'info' = 'info';
  @Input() message?: string | null;
}
