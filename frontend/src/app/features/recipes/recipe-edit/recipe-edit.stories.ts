import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { applicationConfig } from '@storybook/angular';
import { Meta, StoryObj, moduleMetadata } from '@storybook/angular';
import { RecipeEditComponent } from './recipe-edit.component';
import { TranslocoModule } from '@jsverse/transloco';
import { RouterModule } from '@angular/router';
import {} from '@angular/platform-browser/animations';
import { RecipeService } from '../../../core/services/recipe.service';
import { of } from 'rxjs';
import { ActivatedRoute } from '@angular/router';

const mockRecipeService = {
  getRecipe: () =>
    of({ id: 1, name: 'Sample Recipe', url: 'https://example.com', memo: 'Sample memo', tags: [] }),
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
