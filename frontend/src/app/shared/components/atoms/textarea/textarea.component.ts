import { Component, forwardRef, Input, OnInit, Injector } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ControlValueAccessor, NG_VALUE_ACCESSOR, NgControl, FormsModule, ReactiveFormsModule, FormControl } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';

@Component({
  selector: 'app-textarea',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule, MatFormFieldModule, MatInputModule],
  templateUrl: './textarea.component.html',
  styleUrl: './textarea.component.scss',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => TextareaComponent),
      multi: true
    }
  ]
})
export class TextareaComponent implements ControlValueAccessor, OnInit {
  @Input() label = '';
  @Input() placeholder = '';
  @Input() rows = 4;
  @Input() required = false;
  @Input() errorMessage: string | null = null;

  control: FormControl | null = null;

  value: any = '';
  disabled = false;

  onChange: (value: any) => void = () => { };
  onTouched: () => void = () => { };

  constructor(private injector: Injector) { }

  ngOnInit() {
    try {
      const ngControl = this.injector.get(NgControl);
      if (ngControl) {
        ngControl.valueAccessor = this;
        setTimeout(() => {
          if (ngControl.control instanceof FormControl) {
            this.control = ngControl.control;
          }
        });
      }
    } catch (e) { }
  }

  writeValue(obj: any): void {
    this.value = obj;
  }

  registerOnChange(fn: any): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: any): void {
    this.onTouched = fn;
  }

  setDisabledState?(isDisabled: boolean): void {
    this.disabled = isDisabled;
  }

  onInput(event: Event) {
    const target = event.target as HTMLTextAreaElement;
    this.value = target.value;
    this.onChange(this.value);
  }

  onBlur() {
    this.onTouched();
  }
}
