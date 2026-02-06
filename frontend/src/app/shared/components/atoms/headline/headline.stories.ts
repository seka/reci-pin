import { Meta, StoryObj } from '@storybook/angular';
import { HeadlineComponent } from './headline.component';

// Define custom args type to include 'content' which is used in the template but not a component input
type HeadlineStoryArgs = HeadlineComponent & { content: string };

const meta: Meta<HeadlineStoryArgs> = {
    title: 'Atoms/Headline',
    component: HeadlineComponent,
    tags: ['autodocs'],
    argTypes: {
        variant: {
            control: 'select',
            options: ['h1', 'h2', 'h3', 'h4', 'h5', 'h6'],
        },
        content: {
            control: 'text',
            description: 'Headline text content',
        },
    },
};

export default meta;
type Story = StoryObj<HeadlineStoryArgs>;

export const Default: Story = {
    args: {
        variant: 'h2',
        content: 'Headline H2',
    },
    render: (args) => ({
        props: args,
        template: `<app-headline [variant]="variant">${args.content}</app-headline>`,
    }),
};

export const H1: Story = {
    args: {
        variant: 'h1',
        content: 'Headline H1',
    },
    render: (args) => ({
        props: args,
        template: `<app-headline [variant]="variant">${args.content}</app-headline>`,
    }),
};

export const H3: Story = {
    args: {
        variant: 'h3',
        content: 'Headline H3',
    },
    render: (args) => ({
        props: args,
        template: `<app-headline [variant]="variant">${args.content}</app-headline>`,
    }),
};
