import { provideRouter } from '@angular/router';
import { applicationConfig } from '@storybook/angular-vite';
import { Meta, StoryObj } from '@storybook/angular-vite';
import { LinkComponent } from './link.component';
import { moduleMetadata } from '@storybook/angular-vite';
import { RouterModule } from '@angular/router';

const meta: Meta<LinkComponent> = {
  title: 'Shared/Atoms/Link',
  component: LinkComponent,
  tags: ['autodocs'],
  decorators: [
    applicationConfig({ providers: [provideRouter([])] }),
    moduleMetadata({
      imports: [RouterModule],
    }),
  ],
  argTypes: {
    variant: {
      control: 'select',
      options: ['primary', 'secondary'],
    },
    routerLink: { control: 'text' },
    href: { control: 'text' },
    title: { control: 'text' },
  },
  render: (args) => ({
    props: args,
    template: `<app-link [routerLink]="routerLink" [href]="href" [title]="title" [variant]="variant">Link Text</app-link>`,
  }),
};

export default meta;
type Story = StoryObj<LinkComponent>;

export const Primary: Story = {
  args: {
    variant: 'primary',
    href: '#',
    title: 'Primary Link Title',
  },
};

export const Secondary: Story = {
  args: {
    variant: 'secondary',
    href: '#',
    title: 'Secondary Link Title',
  },
};
