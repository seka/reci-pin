import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { applicationConfig } from '@storybook/angular';
import { Meta, StoryObj, moduleMetadata } from '@storybook/angular';
import { RecipeFormComponent } from './recipe-form.component';
import { TranslateModule } from '@ngx-translate/core';
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
        TranslateModule,
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
    titleKey: 'RECIPE.NEW_TITLE',
    submitLabelKey: 'RECIPE.SAVE',
    submittingLabelKey: 'RECIPE.SAVING',
    isSubmitting: false,
    initialData: {},
  },
};

export const EditRecipe: Story = {
  args: {
    titleKey: 'RECIPE.EDIT_TITLE',
    submitLabelKey: 'RECIPE.UPDATE',
    submittingLabelKey: 'RECIPE.UPDATING',
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
