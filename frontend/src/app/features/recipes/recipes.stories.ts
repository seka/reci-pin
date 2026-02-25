import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { applicationConfig } from '@storybook/angular';
import { Meta, StoryObj, moduleMetadata } from '@storybook/angular';
import { RecipesComponent } from './recipes.component';
import { TranslocoModule } from '@jsverse/transloco';
import { RouterModule } from '@angular/router';
import {} from '@angular/platform-browser/animations';
import { RecipeService } from '../../core/services/recipe.service';
import { of } from 'rxjs';

const mockRecipeService = {
  getUserRecipes: () =>
    of([
      {
        id: 1,
        name: 'Recipe 1',
        url: 'https://example.com/1',
        memo: 'Memo 1',
        tags: [{ id: 1, name: 'Tag1' }],
      },
      { id: 2, name: 'Recipe 2', url: 'https://example.com/2', memo: 'Memo 2', tags: [] },
    ]),
  getAllTags: () =>
    of([
      { id: 1, name: 'Tag1' },
      { id: 2, name: 'Tag2' },
    ]),
};

const meta: Meta<RecipesComponent> = {
  title: 'Features/Recipes/RecipesList',
  component: RecipesComponent,
  tags: ['autodocs'],
  decorators: [
    applicationConfig({ providers: [provideAnimations(), provideRouter([])] }),
    moduleMetadata({
      imports: [TranslocoModule, RouterModule],
      providers: [{ provide: RecipeService, useValue: mockRecipeService }],
    }),
  ],
};

export default meta;
type Story = StoryObj<RecipesComponent>;

export const Default: Story = {
  args: {},
};
