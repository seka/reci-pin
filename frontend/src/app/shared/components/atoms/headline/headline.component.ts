import { Component, Input, ChangeDetectionStrategy } from '@angular/core';
import { NgTemplateOutlet } from '@angular/common';

@Component({
  selector: 'app-headline',
  standalone: true,
  imports: [NgTemplateOutlet],
  templateUrl: './headline.component.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './headline.component.scss',
})
export class HeadlineComponent {
  @Input() variant: 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6' = 'h2';
}
