import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { applicationConfig } from '@storybook/angular';
import { Meta, StoryObj, moduleMetadata } from '@storybook/angular';
import { SignupComponent } from './signup.component';
import { TranslocoModule } from '@jsverse/transloco';
import { RouterModule } from '@angular/router';
import {  } from '@angular/platform-browser/animations';
import { AuthService } from '../../../core/services/auth.service';
import { of } from 'rxjs';

const mockAuthService = {
  signup: () => of({}),
};

const meta: Meta<SignupComponent> = {
  title: 'Features/Auth/Signup',
  component: SignupComponent,
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
type Story = StoryObj<SignupComponent>;

export const Default: Story = {
  args: {},
};
