import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { applicationConfig } from '@storybook/angular';
import { Meta, StoryObj, moduleMetadata } from '@storybook/angular';
import { RecipeFormComponent } from './recipe-form.component';
import { TranslocoModule } from '@jsverse/transloco';
import { RouterModule } from '@angular/router';
import { HttpClientModule } from '@angular/common/http';
import {  } from '@angular/platform-browser/animations';
import { of } from 'rxjs';
import { RecipeService } from '../../../../core/services/recipe.service';

const mockRecipeService = {
  getAllTags: () => of([{ id: 1, name: 'Tag 1' }, { id: 2, name: 'Tag 2' }]),
};

const meta: Meta<RecipeFormComponent> = {
  title: 'Shared/Organisms/RecipeForm',
  component: RecipeFormComponent,
  tags: ['autodocs'],
  decorators: [
    applicationConfig({ providers: [provideAnimations(), provideRouter([])] }),
    moduleMetadata({
      imports: [
        TranslocoModule,
        RouterModule,
        HttpClientModule,
        
      ],
      providers: [
        { provide: RecipeService, useValue: mockRecipeService },
      ],
    }),
  ],
  argTypes: {
    titleKey: { control: 'text' },
    submitLabelKey: { control: 'text' },
    submittingLabelKey: { control: 'text' },
    isSubmitting: { control: 'boolean' },
    initialData: { control: 'object' },
    initialImagePreview: { control: 'text' },
    save: { action: 'save' },
  },
};

export default meta;
type Story = StoryObj<RecipeFormComponent>;

export const CreateRecipe: Story = {
  args: {
    titleKey: 'FEATURES.RECIPES.RECIPE_CREATE.TITLE',
    submitLabelKey: 'COMPONENTS.ORGANISMS.RECIPE_FORM.SAVE',
    submittingLabelKey: 'COMPONENTS.ORGANISMS.RECIPE_FORM.SAVING',
    isSubmitting: false,
    initialData: {},
  },
};

export const EditRecipe: Story = {
  args: {
    titleKey: 'FEATURES.RECIPES.RECIPE_EDIT.TITLE',
    submitLabelKey: 'COMPONENTS.ORGANISMS.RECIPE_FORM.SAVE',
    submittingLabelKey: 'COMPONENTS.ORGANISMS.RECIPE_FORM.SAVING',
    isSubmitting: false,
    initialData: {
      name: 'Sample Recipe',
      url: 'https://example.com/recipe',
      memo: 'This is a sample memo.',
      tagIds: [1],
    },
    initialImagePreview: 'https://via.placeholder.com/150',
  },
};
