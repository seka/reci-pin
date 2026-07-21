import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { applicationConfig } from '@storybook/angular-vite';
import { Meta, StoryObj, moduleMetadata } from '@storybook/angular-vite';
import { RecipeCreateComponent } from './recipe-create.component';
import { TranslocoModule } from '@jsverse/transloco';
import { RouterModule } from '@angular/router';
import {} from '@angular/platform-browser/animations';
import { RecipeService } from '../../../core/services/recipe.service';
import { of } from 'rxjs';

const mockRecipeService = {
  getAllTags: () => of([]),
  createRecipe: () => of({ id: 1 }),
};

const meta: Meta<RecipeCreateComponent> = {
  title: 'Features/Recipes/RecipeCreate',
  component: RecipeCreateComponent,
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
type Story = StoryObj<RecipeCreateComponent>;

export const Default: Story = {
  args: {},
};
