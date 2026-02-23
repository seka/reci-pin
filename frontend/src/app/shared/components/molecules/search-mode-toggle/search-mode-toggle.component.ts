import { Component, Input, Output, EventEmitter } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonToggleModule } from '@angular/material/button-toggle';

@Component({
  selector: 'app-search-mode-toggle',
  standalone: true,
  imports: [FormsModule, MatButtonToggleModule],
  template: `
    <mat-button-toggle-group [ngModel]="value" (ngModelChange)="onValueChange($event)">
      <mat-button-toggle value="keyword">キーワード</mat-button-toggle>
      <mat-button-toggle value="tag">タグ</mat-button-toggle>
    </mat-button-toggle-group>
  `,
  styles: [
    `
      :host {
        display: block;
      }
      mat-button-toggle-group {
        width: 100%;
        max-width: fit-content;
      }
    `,
  ],
})
export class SearchModeToggleComponent {
  @Input() value: 'keyword' | 'tag' = 'keyword';
  @Output() valueChange = new EventEmitter<'keyword' | 'tag'>();

  onValueChange(newValue: 'keyword' | 'tag') {
    this.valueChange.emit(newValue);
  }
}
