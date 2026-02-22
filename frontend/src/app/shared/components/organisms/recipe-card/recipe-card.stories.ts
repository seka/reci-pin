import { Meta, StoryObj, moduleMetadata } from '@storybook/angular';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { RecipeCardComponent } from './recipe-card.component';
import { Recipe } from '../../../../core/services/recipe.service';

const meta: Meta<RecipeCardComponent> = {
  title: 'Organisms/RecipeCard',
  component: RecipeCardComponent,
  tags: ['autodocs'],
  decorators: [
    moduleMetadata({
      imports: [MatCardModule, MatButtonModule, MatIconModule, MatChipsModule],
    }),
  ],
};

export default meta;
type Story = StoryObj<RecipeCardComponent>;

const mockRecipe: Recipe = {
  id: 1,
  userId: 1,
  name: '手作りハンバーグ',
  url: 'https://example.com/hamburg',
  memo: 'ふっくらジューシーなハンバーグ。玉ねぎはしっかり炒めるのがコツです。',
  images: [
    {
      id: 1,
      recipeId: 1,
      imagePath: '/images/recipe1.jpg',
      imageUrl: 'https://example.com/images/recipe1.jpg',
      createdAt: new Date().toISOString(),
    },
  ],
  createdAt: '2024-02-07T00:00:00Z',
  updatedAt: '2024-02-07T00:00:00Z',
  tags: [
    { id: 1, name: '洋食' },
    { id: 2, name: '夕食' },
  ],
};

export const Default: Story = {
  args: {
    recipe: mockRecipe,
  },
};

export const WithoutTags: Story = {
  args: {
    recipe: {
      ...mockRecipe,
      tags: [],
    },
  },
};

export const WithoutUrl: Story = {
  args: {
    recipe: {
      ...mockRecipe,
      url: '',
    },
  },
};
