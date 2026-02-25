import { Component, EventEmitter, Output } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { ButtonComponent } from '../../atoms/button/button.component';
import { TranslocoPipe } from '@jsverse/transloco';

@Component({
  selector: 'app-search-bar',
  standalone: true,
  imports: [CommonModule, MatIconModule, ButtonComponent, TranslocoPipe],
  templateUrl: './search-bar.component.html',
  styleUrl: './search-bar.component.scss',
})
export class SearchBarComponent {
  @Output() searchSubmit = new EventEmitter<void>();

  onSearch() {
    this.searchSubmit.emit();
  }
}
