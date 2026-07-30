/**
 * @vitest-environment jsdom
 */
import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { TranslocoTestingModule } from '@jsverse/transloco';
import { of } from 'rxjs';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { RecipeService } from '../../../../core/services/recipe.service';
import { RecipeFormComponent } from './recipe-form.component';

describe('RecipeFormComponent', () => {
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        RecipeFormComponent,
        TranslocoTestingModule.forRoot({
          langs: { ja: {} },
          translocoConfig: { availableLangs: ['ja'], defaultLang: 'ja' },
        }),
      ],
      providers: [
        provideRouter([]),
        {
          provide: RecipeService,
          useValue: { getAllTags: () => of([]) },
        },
      ],
    }).compileComponents();
  });

  it('有効なフォームをsubmitするとsaveイベントをemitする', async () => {
    const fixture = TestBed.createComponent(RecipeFormComponent);
    fixture.autoDetectChanges();
    const save = vi.fn();
    fixture.componentInstance.save.subscribe(save);

    const nameInput = fixture.nativeElement.querySelector('input[type="text"]') as HTMLInputElement;
    nameInput.value = 'テストレシピ';
    nameInput.dispatchEvent(new Event('input', { bubbles: true }));
    await fixture.whenStable();

    const form = fixture.nativeElement.querySelector('form') as HTMLFormElement;
    form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));

    expect(save).toHaveBeenCalledWith({
      formData: { name: 'テストレシピ', url: '', memo: '', tagIds: [] },
      file: null,
    });
  });

  it('フォームが無効な場合はsaveイベントをemitしない', () => {
    const fixture = TestBed.createComponent(RecipeFormComponent);
    fixture.autoDetectChanges();
    const save = vi.fn();
    fixture.componentInstance.save.subscribe(save);

    const form = fixture.nativeElement.querySelector('form') as HTMLFormElement;
    form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));

    expect(save).not.toHaveBeenCalled();
  });

  it('未対応形式の画像を選択するとエラーを表示する', async () => {
    const fixture = TestBed.createComponent(RecipeFormComponent);
    fixture.autoDetectChanges();
    const fileInput = fixture.nativeElement.querySelector('input[type="file"]') as HTMLInputElement;
    const file = new File(['image'], 'image.gif', { type: 'image/gif' });
    Object.defineProperty(fileInput, 'files', { value: [file] });

    fileInput.dispatchEvent(new Event('change', { bubbles: true }));
    await fixture.whenStable();

    expect(fixture.nativeElement.textContent).toContain('JPEG, PNG, WebP のみアップロードできます');
  });
});
