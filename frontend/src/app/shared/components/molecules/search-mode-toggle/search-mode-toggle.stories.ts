import { applicationConfig } from '@storybook/angular';
import { provideAnimations } from '@angular/platform-browser/animations';
import { Meta, StoryObj } from '@storybook/angular';
import { SearchModeToggleComponent } from './search-mode-toggle.component';
import {  } from '@angular/platform-browser/animations';
import { moduleMetadata } from '@storybook/angular';

const meta: Meta<SearchModeToggleComponent> = {
  title: 'Shared/Molecules/SearchModeToggle',
  component: SearchModeToggleComponent,
  tags: ['autodocs'],
  decorators: [
    applicationConfig({ providers: [provideAnimations()] }),
    moduleMetadata({
      imports: [],
    }),
  ],
  argTypes: {
    value: {
      control: 'radio',
      options: ['keyword', 'tag'],
    },
    valueChange: { action: 'valueChange' },
  },
};

export default meta;
type Story = StoryObj<SearchModeToggleComponent>;

export const KeywordMode: Story = {
  args: {
    value: 'keyword',
  },
};

export const TagMode: Story = {
  args: {
    value: 'tag',
  },
};
