import { Meta, StoryObj } from '@storybook/angular-vite';
import { SearchBarComponent } from './search-bar.component';

const meta: Meta<SearchBarComponent> = {
  title: 'Shared/Molecules/SearchBar',
  component: SearchBarComponent,
  tags: ['autodocs'],
  argTypes: {
    searchSubmit: { action: 'searchSubmit' },
  },
};

export default meta;
type Story = StoryObj<SearchBarComponent>;

export const Default: Story = {
  args: {},
};
