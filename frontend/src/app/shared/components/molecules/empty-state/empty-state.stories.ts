import { Meta, StoryObj } from '@storybook/angular-vite';
import { EmptyStateComponent } from './empty-state.component';

const meta: Meta<EmptyStateComponent> = {
  title: 'Shared/Molecules/EmptyState',
  component: EmptyStateComponent,
  tags: ['autodocs'],
  argTypes: {
    icon: { control: 'text' },
    title: { control: 'text' },
    message: { control: 'text' },
  },
};

export default meta;
type Story = StoryObj<EmptyStateComponent>;

export const Default: Story = {
  args: {
    icon: 'inbox',
    title: 'No Items Found',
    message: 'There are no items to display at this moment.',
  },
};

export const CustomMessage: Story = {
  args: {
    icon: 'search_off',
    title: 'No Results',
    message: 'Try adjusting your search filters.',
  },
};
