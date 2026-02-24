import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter } from '@angular/router';
import { applicationConfig } from '@storybook/angular';
import { Meta, StoryObj, moduleMetadata } from '@storybook/angular';
import { SettingsComponent } from './settings.component';
import { TranslateModule } from '@ngx-translate/core';
import { RouterModule } from '@angular/router';
import {  } from '@angular/platform-browser/animations';
import { AuthService } from '../../core/services/auth.service';
import { of } from 'rxjs';

const mockAuthService = {
  currentUser$: of({ id: 1, name: 'Test User', email: 'test@example.com' }),
  withdraw: () => of(null),
};

const meta: Meta<SettingsComponent> = {
  title: 'Features/Settings/SettingsMain',
  component: SettingsComponent,
  tags: ['autodocs'],
  decorators: [
    applicationConfig({ providers: [provideAnimations(), provideRouter([])] }),
    moduleMetadata({
      imports: [
        TranslateModule,
        RouterModule,
        
      ],
      providers: [
        { provide: AuthService, useValue: mockAuthService },
      ],
    }),
  ],
};

export default meta;
type Story = StoryObj<SettingsComponent>;

export const Default: Story = {
  args: {},
};
