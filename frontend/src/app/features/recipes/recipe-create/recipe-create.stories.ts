import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { applicationConfig } from '@storybook/angular';
import { Meta, StoryObj, moduleMetadata } from '@storybook/angular';
import { RecipeCreateComponent } from './recipe-create.component';
import { TranslateModule } from '@ngx-translate/core';
import { RouterModule } from '@angular/router';
import {  } from '@angular/platform-browser/animations';
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
      imports: [
        TranslateModule,
        RouterModule,
        
      ],
      providers: [
        { provide: RecipeService, useValue: mockRecipeService },
      ],
    }),
  ],
};

export default meta;
type Story = StoryObj<RecipeCreateComponent>;

export const Default: Story = {
  args: {},
};
