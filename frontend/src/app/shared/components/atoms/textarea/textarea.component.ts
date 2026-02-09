import { Component, forwardRef, inject, Input, OnInit, Injector } from '@angular/core';
import {
  ControlValueAccessor,
  NG_VALUE_ACCESSOR,
  NgControl,
  FormsModule,
  ReactiveFormsModule,
  FormControl,
} from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';

@Component({
  selector: 'app-textarea',
  standalone: true,
  imports: [FormsModule, ReactiveFormsModule, MatFormFieldModule, MatInputModule],
  templateUrl: './textarea.component.html',
  styleUrl: './textarea.component.scss',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => TextareaComponent),
      multi: true,
    },
  ],
})
export class TextareaComponent implements ControlValueAccessor, OnInit {
  private readonly injector = inject(Injector);

  @Input() label = '';
  @Input() placeholder = '';
  @Input() rows = 4;
  @Input() required = false;
  @Input() errorMessage: string | null = null;

  control: FormControl | null = null;

  value: string = '';
  disabled = false;

  onChange: (value: string) => void = () => {
    // Placeholder for ControlValueAccessor - implemented in registerOnChange
  };
  onTouched: () => void = () => {
    // Placeholder for ControlValueAccessor - implemented in registerOnTouched
  };

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
    } catch {
      // Standalone usage without form control
    }
  }

  writeValue(obj: string): void {
    this.value = obj;
  }

  registerOnChange(fn: (value: string) => void): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: () => void): void {
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
