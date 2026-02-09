import { Component, Input } from '@angular/core';
import { NgTemplateOutlet } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';

@Component({
  selector: 'app-button',
  standalone: true,
  imports: [NgTemplateOutlet, MatButtonModule],
  template: `
    <ng-template #content><ng-content></ng-content></ng-template>

    @switch (variant) {
      @case ('primary') {
        <button
          mat-flat-button
          color="primary"
          [type]="type"
          [disabled]="disabled"
          (click)="handleClick($event)"
        >
          <ng-container *ngTemplateOutlet="content"></ng-container>
        </button>
      }
      @case ('secondary') {
        <button
          mat-flat-button
          color="accent"
          [type]="type"
          [disabled]="disabled"
          (click)="handleClick($event)"
        >
          <ng-container *ngTemplateOutlet="content"></ng-container>
        </button>
      }
      @case ('outline') {
        <button
          mat-stroked-button
          color="primary"
          [type]="type"
          [disabled]="disabled"
          (click)="handleClick($event)"
        >
          <ng-container *ngTemplateOutlet="content"></ng-container>
        </button>
      }
      @case ('text') {
        <button mat-button [type]="type" [disabled]="disabled" (click)="handleClick($event)">
          <ng-container *ngTemplateOutlet="content"></ng-container>
        </button>
      }
      @case ('warn') {
        <button
          mat-button
          color="warn"
          [type]="type"
          [disabled]="disabled"
          (click)="handleClick($event)"
        >
          <ng-container *ngTemplateOutlet="content"></ng-container>
        </button>
      }
      @case ('accent') {
        <button
          mat-button
          color="accent"
          [type]="type"
          [disabled]="disabled"
          (click)="handleClick($event)"
        >
          <ng-container *ngTemplateOutlet="content"></ng-container>
        </button>
      }
    }
  `,
  styles: [
    `
      :host {
        display: inline-block;
      }
      button {
        width: 100%;
        min-width: 120px;
        font-weight: bold;
      }

      /* Ensure white text for accent (secondary) filled buttons in MDC */
      :host ::ng-deep .mat-mdc-unelevated-button.mat-accent {
        --mdc-filled-button-label-text-color: #fff;
        color: #fff !important;
      }
      :host ::ng-deep .mat-mdc-unelevated-button.mat-accent .mdc-button__label {
        color: #fff !important;
      }
    `,
  ],
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
