import { Meta, StoryObj, applicationConfig } from '@storybook/angular-vite';
import { provideRouter } from '@angular/router';
import { APP_BASE_HREF } from '@angular/common';
import { LogoComponent } from './logo.component';

const meta: Meta<LogoComponent> = {
  title: 'Atoms/Logo',
  component: LogoComponent,
  tags: ['autodocs'],
  decorators: [
    applicationConfig({
      providers: [provideRouter([]), { provide: APP_BASE_HREF, useValue: '/' }],
    }),
  ],
  argTypes: {
    size: {
      control: 'select',
      options: ['small', 'medium', 'large'],
    },
  },
};

export default meta;
type Story = StoryObj<LogoComponent>;

export const Medium: Story = {
  args: {
    size: 'medium',
  },
};

export const Small: Story = {
  args: {
    size: 'small',
  },
};

export const Large: Story = {
  args: {
    size: 'large',
  },
};
