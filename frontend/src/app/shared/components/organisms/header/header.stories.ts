import { Meta, StoryObj, moduleMetadata, applicationConfig } from '@storybook/angular-vite';
import { provideRouter } from '@angular/router';
import { APP_BASE_HREF } from '@angular/common';
import { of } from 'rxjs';
import { HeaderComponent } from './header.component';
import { AuthService } from '../../../../core/services/auth.service';
import { LogoComponent } from '../../atoms/logo/logo.component';
import { ButtonComponent } from '../../atoms/button/button.component';

// Mock AuthService
const mockAuthService = {
  currentUser$: of({ id: 1, name: 'Tanaka Taro', email: 'tanaka@example.com' }),
  logout: () => console.log('Logged out'),
};

const meta: Meta<HeaderComponent> = {
  title: 'Organisms/Header',
  component: HeaderComponent,
  tags: ['autodocs'],
  decorators: [
    applicationConfig({
      providers: [
        provideRouter([]),
        { provide: AuthService, useValue: mockAuthService },
        { provide: APP_BASE_HREF, useValue: '/' },
      ],
    }),
    moduleMetadata({
      imports: [LogoComponent, ButtonComponent],
    }),
  ],
};

export default meta;
type Story = StoryObj<HeaderComponent>;

export const LoggedIn: Story = {};

export const LoggedOut: Story = {
  decorators: [
    applicationConfig({
      providers: [
        {
          provide: AuthService,
          useValue: { ...mockAuthService, currentUser$: of(null) },
        },
      ],
    }),
  ],
};
