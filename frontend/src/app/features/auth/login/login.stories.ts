import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { applicationConfig } from '@storybook/angular';
import { Meta, StoryObj, moduleMetadata } from '@storybook/angular';
import { LoginComponent } from './login.component';
import { TranslocoModule } from '@jsverse/transloco';
import { RouterModule } from '@angular/router';
import {  } from '@angular/platform-browser/animations';
import { AuthService } from '../../../core/services/auth.service';
import { of } from 'rxjs';

const mockAuthService = {
  login: () => of({}),
};

const meta: Meta<LoginComponent> = {
  title: 'Features/Auth/Login',
  component: LoginComponent,
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
type Story = StoryObj<LoginComponent>;

export const Default: Story = {
  args: {},
};
