import {
  Component,
  forwardRef,
  inject,
  Input,
  OnInit,
  Injector,
  Output,
  EventEmitter,
  ChangeDetectionStrategy,
} from '@angular/core';
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
import { TranslocoPipe } from '@jsverse/transloco';

@Component({
  selector: 'app-input',
  standalone: true,
  imports: [FormsModule, ReactiveFormsModule, MatFormFieldModule, MatInputModule, TranslocoPipe],
  templateUrl: './input.component.html',
  styleUrl: './input.component.scss',
  changeDetection: ChangeDetectionStrategy.Eager,
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => InputComponent),
      multi: true,
    },
  ],
})
export class InputComponent implements ControlValueAccessor, OnInit {
  private readonly injector = inject(Injector);

  @Input() label = '';
  @Input() placeholder = '';
  @Input() type: 'text' | 'password' | 'email' | 'number' = 'text';
  @Input() required = false;
  // Signal Forms の maxLength() バリデータは metadata を number | undefined として公開し、
  // [formField] 経由でこの Input に自動バインドされるため undefined も許容する。
  @Input() maxLength: number | null | undefined = null;
  @Input() showCounter = false;
  @Input() errorMessage: string | string[] | null = null;
  @Input() focus = false;
  @Input() floatLabel: 'always' | 'auto' = 'auto';
  @Input() disabled = false;

  get currentLength(): number {
    return String(this.value || '').length;
  }

  @Output() inputBlur = new EventEmitter<void>();
  @Output() valueChange = new EventEmitter<string | number>();

  get errorMessages(): string[] {
    if (!this.errorMessage) return [];
    if (Array.isArray(this.errorMessage)) return this.errorMessage;
    return [this.errorMessage];
  }

  control: FormControl | null = null;

  value: string | number = '';

  onChange: (value: string | number) => void = () => {
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

  writeValue(obj: string | number): void {
    this.value = obj;
  }

  registerOnChange(fn: (value: string | number) => void): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: () => void): void {
    this.onTouched = fn;
  }

  setDisabledState?(isDisabled: boolean): void {
    this.disabled = isDisabled;
  }

  onInput(event: Event) {
    const target = event.target as HTMLInputElement;
    this.value = target.value;
    this.onChange(this.value);
    this.valueChange.emit(this.value);
  }

  onBlur() {
    this.onTouched();
    this.inputBlur.emit();
  }
}
