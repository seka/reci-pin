import { Meta, StoryObj } from '@storybook/angular-vite';
import { TextComponent } from './text.component';

// Define custom args type to include 'content' which is used in the template but not a component input
type TextStoryArgs = TextComponent & { content: string };

const meta: Meta<TextStoryArgs> = {
  title: 'Atoms/Text',
  component: TextComponent,
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'select',
      options: ['body', 'caption', 'label'],
    },
    content: {
      control: 'text',
      description: 'Text content',
    },
  },
};

export default meta;
type Story = StoryObj<TextStoryArgs>;

export const Body: Story = {
  args: {
    variant: 'body',
    content: 'This is body text. A descriptive paragraph about something interesting.',
  },
  render: (args) => ({
    props: args,
    template: `<app-text [variant]="variant">${args.content}</app-text>`,
  }),
};

export const Caption: Story = {
  args: {
    variant: 'caption',
    content: 'This is a caption text.',
  },
  render: (args) => ({
    props: args,
    template: `<app-text [variant]="variant">${args.content}</app-text>`,
  }),
};

export const Label: Story = {
  args: {
    variant: 'label',
    content: 'Input Label',
  },
  render: (args) => ({
    props: args,
    template: `<app-text [variant]="variant">${args.content}</app-text>`,
  }),
};
