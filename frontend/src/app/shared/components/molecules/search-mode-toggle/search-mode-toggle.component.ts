import { Component, Input, Output, EventEmitter, ChangeDetectionStrategy } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { TranslocoPipe } from '@jsverse/transloco';

@Component({
  selector: 'app-search-mode-toggle',
  standalone: true,
  imports: [FormsModule, MatButtonToggleModule, TranslocoPipe],
  templateUrl: './search-mode-toggle.component.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './search-mode-toggle.component.scss',
})
export class SearchModeToggleComponent {
  @Input() value: 'keyword' | 'tag' = 'keyword';
  @Output() valueChange = new EventEmitter<'keyword' | 'tag'>();

  onValueChange(newValue: 'keyword' | 'tag') {
    this.valueChange.emit(newValue);
  }
}
