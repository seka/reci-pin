import { Meta, StoryObj } from '@storybook/angular';
import { ButtonComponent } from './button.component';

// Define custom args type to include 'label' which is used in the template but not a component input
type ButtonStoryArgs = ButtonComponent & { label: string };

const meta: Meta<ButtonStoryArgs> = {
    title: 'Atoms/Button',
    component: ButtonComponent,
    tags: ['autodocs'],
    argTypes: {
        variant: {
            control: 'select',
            options: ['primary', 'secondary', 'outline', 'text', 'warn', 'accent'],
        },
        disabled: {
            control: 'boolean',
        },
        label: {
            control: 'text',
            description: 'Button text content',
        },
    },
};

export default meta;
type Story = StoryObj<ButtonStoryArgs>;

export const Primary: Story = {
    args: {
        variant: 'primary',
        label: 'Primary Button',
    },
    render: (args) => ({
        props: args,
        template: `<app-button [variant]="variant" [disabled]="disabled">${args.label}</app-button>`,
    }),
};

export const Secondary: Story = {
    args: {
        variant: 'secondary',
        label: 'Secondary Button',
    },
    render: (args) => ({
        props: args,
        template: `<app-button [variant]="variant" [disabled]="disabled">${args.label}</app-button>`,
    }),
};

export const Outline: Story = {
    args: {
        variant: 'outline',
        label: 'Outline Button',
    },
    render: (args) => ({
        props: args,
        template: `<app-button [variant]="variant" [disabled]="disabled">${args.label}</app-button>`,
    }),
};

export const Text: Story = {
    args: {
        variant: 'text',
        label: 'Text Button',
    },
    render: (args) => ({
        props: args,
        template: `<app-button [variant]="variant" [disabled]="disabled">${args.label}</app-button>`,
    }),
};

export const Disabled: Story = {
    args: {
        variant: 'primary',
        label: 'Disabled Button',
        disabled: true,
    },
    render: (args) => ({
        props: args,
        template: `<app-button [variant]="variant" [disabled]="disabled">${args.label}</app-button>`,
    }),
};
