import { Component, Input, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

@Component({
  selector: 'app-loading-state',
  standalone: true,
  imports: [CommonModule, MatProgressSpinnerModule],
  templateUrl: './loading-state.component.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './loading-state.component.scss',
})
export class LoadingStateComponent {
  @Input() message?: string;
}
