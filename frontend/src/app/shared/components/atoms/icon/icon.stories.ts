import { Meta, StoryObj, moduleMetadata } from '@storybook/angular';
import { MatIconModule } from '@angular/material/icon';
import { IconComponent } from './icon.component';

// Define custom args type to include 'iconName' which is used in the template via ng-content
type IconStoryArgs = IconComponent & { iconName: string };

const meta: Meta<IconStoryArgs> = {
    title: 'Atoms/Icon',
    component: IconComponent,
    tags: ['autodocs'],
    decorators: [
        moduleMetadata({
            imports: [MatIconModule],
        }),
    ],
    argTypes: {
        size: {
            control: 'select',
            options: ['sm', 'md', 'lg'],
        },
        color: {
            control: 'select',
            options: ['inherit', 'primary', 'secondary', 'warn'],
        },
        iconName: {
            control: 'text',
            description: 'Material Icon name',
        },
    },
};

export default meta;
type Story = StoryObj<IconStoryArgs>;

export const Default: Story = {
    args: {
        size: 'md',
        color: 'inherit',
        iconName: 'home',
    },
    render: (args) => ({
        props: args,
        template: `<app-icon [size]="size" [color]="color">${args.iconName}</app-icon>`,
    }),
};

export const Primary: Story = {
    args: {
        size: 'lg',
        color: 'primary',
        iconName: 'favorite',
    },
    render: (args) => ({
        props: args,
        template: `<app-icon [size]="size" [color]="color">${args.iconName}</app-icon>`,
    }),
};

export const Secondary: Story = {
    args: {
        size: 'md',
        color: 'secondary',
        iconName: 'settings',
    },
    render: (args) => ({
        props: args,
        template: `<app-icon [size]="size" [color]="color">${args.iconName}</app-icon>`,
    }),
};
