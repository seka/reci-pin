import { Meta, StoryObj, applicationConfig } from '@storybook/angular';
import { InputComponent } from './input.component';
import { provideAnimations } from '@angular/platform-browser/animations';

const meta: Meta<InputComponent> = {
    title: 'Atoms/Input',
    component: InputComponent,
    tags: ['autodocs'],
    decorators: [
        applicationConfig({
            providers: [provideAnimations()],
        }),
    ],
    argTypes: {
        type: {
            control: 'select',
            options: ['text', 'password', 'email', 'number'],
        },
    },
};

export default meta;
type Story = StoryObj<InputComponent>;

export const Default: Story = {
    args: {
        label: 'Username',
        placeholder: 'Enter your username',
        type: 'text',
    },
};

export const Password: Story = {
    args: {
        label: 'Password',
        placeholder: 'Enter password',
        type: 'password',
    },
};

export const Required: Story = {
    args: {
        label: 'Email',
        placeholder: 'Enter email',
        type: 'email',
        required: true,
    },
};

export const WithError: Story = {
    args: {
        label: 'Error State',
        placeholder: 'Invalid input',
        errorMessage: 'This field is required',
    },
};
