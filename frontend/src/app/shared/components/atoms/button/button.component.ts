import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';

@Component({
  selector: 'app-button',
  standalone: true,
  imports: [CommonModule, MatButtonModule],
  template: `
    <ng-template #content><ng-content></ng-content></ng-template>

    <!-- Primary (Filled) -->
    <button mat-flat-button color="primary" *ngIf="variant === 'primary'" [type]="type" [disabled]="disabled" (click)="handleClick($event)">
      <ng-container *ngTemplateOutlet="content"></ng-container>
    </button>

    <!-- Secondary (Accent/Teal Filled) -->
    <button mat-flat-button color="accent" *ngIf="variant === 'secondary'" [type]="type" [disabled]="disabled" (click)="handleClick($event)">
      <ng-container *ngTemplateOutlet="content"></ng-container>
    </button>

    <!-- Outline (Stroked) -->
    <button mat-stroked-button color="primary" *ngIf="variant === 'outline'" [type]="type" [disabled]="disabled" (click)="handleClick($event)">
      <ng-container *ngTemplateOutlet="content"></ng-container>
    </button>

    <!-- Text (Basic) -->
    <button mat-button *ngIf="variant === 'text'" [type]="type" [disabled]="disabled" (click)="handleClick($event)">
      <ng-container *ngTemplateOutlet="content"></ng-container>
    </button>

    <!-- Warn (Error) -->
    <button mat-button color="warn" *ngIf="variant === 'warn'" [type]="type" [disabled]="disabled" (click)="handleClick($event)">
      <ng-container *ngTemplateOutlet="content"></ng-container>
    </button>

    <!-- Accent (Text Link) -->
    <button mat-button color="accent" *ngIf="variant === 'accent'" [type]="type" [disabled]="disabled" (click)="handleClick($event)">
      <ng-container *ngTemplateOutlet="content"></ng-container>
    </button>
  `,
  styles: [`
    :host { display: inline-block; }
    button { width: 100%; min-width: 120px; font-weight: bold; }
  `]
})
export class ButtonComponent {
  @Input() variant: 'primary' | 'secondary' | 'outline' | 'text' | 'warn' | 'accent' = 'primary';
  @Input() type: 'button' | 'submit' = 'button';
  @Input() disabled = false;

  handleClick(event: Event) {
    if (this.disabled) {
      event.preventDefault();
      event.stopPropagation();
    }
  }
}
