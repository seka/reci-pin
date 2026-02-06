import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { HeadlineComponent } from '../../atoms/headline/headline.component';

@Component({
  selector: 'app-auth-card',
  standalone: true,
  imports: [CommonModule, MatCardModule, HeadlineComponent],
  templateUrl: './auth-card.component.html',
  styleUrl: './auth-card.component.scss'
})
export class AuthCardComponent {
  @Input() title = '';
}
