import { provideAnimations } from '@angular/platform-browser/animations';
import { applicationConfig, Meta, StoryObj } from '@storybook/angular-vite';
import { LoadingStateComponent } from './loading-state.component';

const meta: Meta<LoadingStateComponent> = {
  title: 'Shared/Molecules/LoadingState',
  component: LoadingStateComponent,
  tags: ['autodocs'],
  decorators: [applicationConfig({ providers: [provideAnimations()] })],
  argTypes: {
    message: { control: 'text' },
  },
};

export default meta;
type Story = StoryObj<LoadingStateComponent>;

export const Default: Story = {
  args: {},
};

export const WithMessage: Story = {
  args: {
    message: 'Loading...',
  },
};
