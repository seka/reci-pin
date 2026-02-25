import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { applicationConfig } from '@storybook/angular';
import { Meta, StoryObj, moduleMetadata } from '@storybook/angular';
import { ResetPasswordComponent } from './reset-password.component';
import { TranslocoModule } from '@jsverse/transloco';
import { RouterModule } from '@angular/router';
import {  } from '@angular/platform-browser/animations';
import { AuthService } from '../../../core/services/auth.service';
import { of } from 'rxjs';

const mockAuthService = {
  resetPassword: () => of({ message: 'Password reset successful' }),
};

const meta: Meta<ResetPasswordComponent> = {
  title: 'Features/Auth/ResetPassword',
  component: ResetPasswordComponent,
  tags: ['autodocs'],
  decorators: [
    applicationConfig({ providers: [provideAnimations(), provideRouter([])] }),
    moduleMetadata({
      imports: [
        TranslocoModule,
        RouterModule,
        
      ],
      providers: [
        { provide: AuthService, useValue: mockAuthService },
      ],
    }),
  ],
};

export default meta;
type Story = StoryObj<ResetPasswordComponent>;

export const Default: Story = {
  args: {},
};
