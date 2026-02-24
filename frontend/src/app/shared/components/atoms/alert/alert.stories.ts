import { Meta, StoryObj } from '@storybook/angular';
import { AlertComponent } from './alert.component';

const meta: Meta<AlertComponent> = {
  title: 'Shared/Atoms/Alert',
  component: AlertComponent,
  tags: ['autodocs'],
  argTypes: {
    type: {
      control: 'select',
      options: ['error', 'success', 'info'],
    },
    message: {
      control: 'text',
    },
  },
};

export default meta;
type Story = StoryObj<AlertComponent>;

export const Info: Story = {
  args: {
    type: 'info',
    message: 'This is an information message.',
  },
};

export const Success: Story = {
  args: {
    type: 'success',
    message: 'This is a success message.',
  },
};

export const Error: Story = {
  args: {
    type: 'error',
    message: 'This is an error message.',
  },
};
