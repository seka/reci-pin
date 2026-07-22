import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { applicationConfig } from '@storybook/angular-vite';
import { Meta, StoryObj, moduleMetadata } from '@storybook/angular-vite';
import { RequestPasswordResetComponent } from './request-password-reset.component';
import { TranslocoModule } from '@jsverse/transloco';
import { RouterModule } from '@angular/router';
import {} from '@angular/platform-browser/animations';
import { AuthService } from '../../../core/services/auth.service';
import { of } from 'rxjs';

const mockAuthService = {
  requestPasswordReset: () => of({ message: 'Reset link sent' }),
};

const meta: Meta<RequestPasswordResetComponent> = {
  title: 'Features/Auth/RequestPasswordReset',
  component: RequestPasswordResetComponent,
  tags: ['autodocs'],
  decorators: [
    applicationConfig({ providers: [provideAnimations(), provideRouter([])] }),
    moduleMetadata({
      imports: [TranslocoModule, RouterModule],
      providers: [{ provide: AuthService, useValue: mockAuthService }],
    }),
  ],
};

export default meta;
type Story = StoryObj<RequestPasswordResetComponent>;

export const Default: Story = {
  args: {},
};
