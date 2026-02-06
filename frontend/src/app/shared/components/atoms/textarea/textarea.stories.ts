import { Meta, StoryObj, moduleMetadata, applicationConfig } from '@storybook/angular';
import { provideAnimations } from '@angular/platform-browser/animations';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { TextareaComponent } from './textarea.component';

const meta: Meta<TextareaComponent> = {
    title: 'Atoms/Textarea',
    component: TextareaComponent,
    tags: ['autodocs'],
    decorators: [
        applicationConfig({
            providers: [provideAnimations()],
        }),
        moduleMetadata({
            imports: [ReactiveFormsModule, FormsModule, MatFormFieldModule, MatInputModule],
        }),
    ],
    argTypes: {
        rows: { control: { type: 'number', min: 1, max: 20 } },
    },
};

export default meta;
type Story = StoryObj<TextareaComponent>;

export const Default: Story = {
    args: {
        label: 'メモ',
        placeholder: 'レシピのメモを入力してください...',
        rows: 4,
        required: false,
    },
};

export const Required: Story = {
    args: {
        label: 'コメント',
        placeholder: '必須項目です',
        rows: 3,
        required: true,
    },
};

export const WithError: Story = {
    args: {
        label: '説明',
        placeholder: '入力を間違えています',
        errorMessage: 'この項目は必須です。',
        required: true,
    },
};

export const Disabled: Story = {
    args: {
        label: '無効な入力',
        placeholder: '入力できません',
    },
    render: (args) => ({
        props: args,
        template: `<app-textarea [label]="label" [placeholder]="placeholder" [disabled]="true"></app-textarea>`,
    }),
};
