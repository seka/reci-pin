import { Component, forwardRef, Input, OnInit, Injector } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ControlValueAccessor, NG_VALUE_ACCESSOR, NgControl, FormsModule, ReactiveFormsModule, FormControl } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';

@Component({
  selector: 'app-input',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule, MatFormFieldModule, MatInputModule],
  templateUrl: './input.component.html',
  styleUrl: './input.component.scss',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => InputComponent),
      multi: true
    }
  ]
})
export class InputComponent implements ControlValueAccessor, OnInit {
  @Input() label = '';
  @Input() placeholder = '';
  @Input() type: 'text' | 'password' | 'email' | 'number' = 'text';
  @Input() required = false;
  @Input() errorMessage: string | null = null;

  control: FormControl | null = null;

  // ControlValueAccessor internal value
  value: any = '';
  disabled = false;

  onChange: (value: any) => void = () => { };
  onTouched: () => void = () => { };

  constructor(private injector: Injector) { }

  ngOnInit() {
    // Attempt to get the NgControl associated with this component to show/hide errors based on touch state
    try {
      const ngControl = this.injector.get(NgControl);
      if (ngControl) {
        ngControl.valueAccessor = this;
        // Wait for next tick/lifecycle to get the control instance
        setTimeout(() => {
          if (ngControl.control instanceof FormControl) {
            this.control = ngControl.control;
          }
        });
      }
    } catch (e) {
      // Standalone usage without form control
    }
  }

  // Value Accessor Methods
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

  // Handle Input Event
  onInput(event: Event) {
    const target = event.target as HTMLInputElement;
    this.value = target.value;
    this.onChange(this.value);
  }

  onBlur() {
    this.onTouched();
  }
}
