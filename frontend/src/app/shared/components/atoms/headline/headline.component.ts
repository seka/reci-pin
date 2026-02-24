import { Component, Input } from '@angular/core';
import { NgTemplateOutlet } from '@angular/common';

@Component({
  selector: 'app-headline',
  standalone: true,
  imports: [NgTemplateOutlet],
  templateUrl: './headline.component.html',
  styleUrl: './headline.component.scss',
})
export class HeadlineComponent {
  @Input() variant: 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6' = 'h2';
}
