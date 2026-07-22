import { Component, Input, ChangeDetectionStrategy } from '@angular/core';
import { NgTemplateOutlet } from '@angular/common';

@Component({
  selector: 'app-text',
  standalone: true,
  imports: [NgTemplateOutlet],
  templateUrl: './text.component.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './text.component.scss',
})
export class TextComponent {
  @Input() variant: 'body' | 'caption' | 'label' = 'body';
}
