import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { applicationConfig } from '@storybook/angular-vite';
import { Meta, StoryObj, moduleMetadata } from '@storybook/angular-vite';
import { RecipeEditComponent } from './recipe-edit.component';
import { TranslocoModule } from '@jsverse/transloco';
import { RouterModule } from '@angular/router';
import {} from '@angular/platform-browser/animations';
import { RecipeService } from '../../../core/services/recipe.service';
import { of, throwError } from 'rxjs';
import { ActivatedRoute } from '@angular/router';

const mockRecipeService = {
  getRecipe: () =>
    of({ id: 1, name: 'Sample Recipe', url: 'https://example.com', memo: 'Sample memo', tags: [] }),
  getAllTags: () => of([]),
};

const mockErrorRecipeService = {
  getRecipe: () => throwError(() => new Error('Failed to fetch recipe')),
  getAllTags: () => of([]),
};

const mockActivatedRoute = {
  snapshot: {
    paramMap: {
      get: () => '1',
    },
  },
};

const meta: Meta<RecipeEditComponent> = {
  title: 'Features/Recipes/RecipeEdit',
  component: RecipeEditComponent,
  tags: ['autodocs'],
  decorators: [
    applicationConfig({ providers: [provideAnimations(), provideRouter([])] }),
    moduleMetadata({
      imports: [TranslocoModule, RouterModule],
      providers: [
        { provide: RecipeService, useValue: mockRecipeService },
        { provide: ActivatedRoute, useValue: mockActivatedRoute },
      ],
    }),
  ],
};

export default meta;
type Story = StoryObj<RecipeEditComponent>;

export const Default: Story = {
  args: {},
};

// Verifies that a failed getRecipe() call redirects away via the effect()
// watching recipeResource.error(), same as the original subscribe(error)
// callback did.
export const LoadError: Story = {
  args: {},
  decorators: [
    moduleMetadata({
      providers: [{ provide: RecipeService, useValue: mockErrorRecipeService }],
    }),
  ],
};
