import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { applicationConfig } from '@storybook/angular-vite';
import { Meta, StoryObj, moduleMetadata } from '@storybook/angular-vite';
import { RecipesComponent } from './recipes.component';
import { TranslocoModule } from '@jsverse/transloco';
import { RouterModule } from '@angular/router';
import {} from '@angular/platform-browser/animations';
import { RecipeService } from '../../core/services/recipe.service';
import { of, delay } from 'rxjs';

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

const mockEmptyRecipeService = {
  getUserRecipes: () => of([]),
  getAllTags: () => of([]),
};

// Never-resolving getUserRecipes() keeps recipesResource in isLoading() === true
// forever, so this story documents/verifies that the "no recipes yet" empty
// state is not shown while the initial fetch is still pending.
const mockLoadingRecipeService = {
  getUserRecipes: () => of([]).pipe(delay(1_000_000)),
  getAllTags: () => of([]).pipe(delay(1_000_000)),
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

export const Empty: Story = {
  args: {},
  decorators: [
    moduleMetadata({
      providers: [{ provide: RecipeService, useValue: mockEmptyRecipeService }],
    }),
  ],
};

export const Loading: Story = {
  args: {},
  decorators: [
    moduleMetadata({
      providers: [{ provide: RecipeService, useValue: mockLoadingRecipeService }],
    }),
  ],
};
