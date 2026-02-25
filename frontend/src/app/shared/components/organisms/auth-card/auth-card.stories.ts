import { Meta, StoryObj, moduleMetadata } from '@storybook/angular';
import { MatCardModule } from '@angular/material/card';
import { AuthCardComponent } from './auth-card.component';
import { HeadlineComponent } from '../../atoms/headline/headline.component';

const meta: Meta<AuthCardComponent> = {
  title: 'Organisms/AuthCard',
  component: AuthCardComponent,
  tags: ['autodocs'],
  decorators: [
    moduleMetadata({
      imports: [MatCardModule, HeadlineComponent],
    }),
  ],
};

export default meta;
type Story = StoryObj<AuthCardComponent>;

export const Login: Story = {
  args: {
    title: 'ログイン',
  },
  render: (args) => ({
    props: args,
    template: `
      <app-auth-card [title]="title">
        <div style="padding: 20px; text-align: center;">
          <p>ここにログインフォームが入ります</p>
          <div style="margin-top: 20px;">
            <button style="width: 100%; padding: 10px; background: #e91e63; color: white; border: none; border-radius: var(--radius-1);">ログイン</button>
          </div>
        </div>
      </app-auth-card>
    `,
  }),
};

export const Signup: Story = {
  args: {
    title: '新規登録',
  },
  render: (args) => ({
    props: args,
    template: `
      <app-auth-card [title]="title">
        <div style="padding: 20px; text-align: center;">
          <p>ここにユーザー登録フォームが入ります</p>
          <div style="margin-top: 20px;">
            <button style="width: 100%; padding: 10px; background: #00bcd4; color: white; border: none; border-radius: var(--radius-1);">登録する</button>
          </div>
        </div>
      </app-auth-card>
    `,
  }),
};
