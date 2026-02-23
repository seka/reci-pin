import { Routes } from '@angular/router';
import { authGuard } from './core/guards/auth.guard';
import { guestGuard } from './core/guards/guest.guard';

export const routes: Routes = [
  { path: '', redirectTo: '/login', pathMatch: 'full' },
  {
    path: 'login',
    canActivate: [guestGuard],
    loadComponent: () =>
      import('./features/auth/login/login.component').then((m) => m.LoginComponent),
  },
  {
    path: 'signup',
    loadComponent: () =>
      import('./features/auth/signup/signup.component').then((m) => m.SignupComponent),
  },
  {
    path: 'password-reset/request',
    loadComponent: () =>
      import('./features/auth/request-password-reset/request-password-reset.component').then(
        (m) => m.RequestPasswordResetComponent,
      ),
  },
  {
    path: 'password-reset',
    loadComponent: () =>
      import('./features/auth/reset-password/reset-password.component').then(
        (m) => m.ResetPasswordComponent,
      ),
  },
  {
    path: 'recipes/new',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./features/recipes/recipe-create/recipe-create.component').then(
        (m) => m.RecipeCreateComponent,
      ),
  },
  {
    path: 'recipes/:id/edit',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./features/recipes/recipe-edit/recipe-edit.component').then(
        (m) => m.RecipeEditComponent,
      ),
  },
  {
    path: 'recipes',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./features/recipes/recipes.component').then((m) => m.RecipesComponent),
  },
  {
    path: 'settings',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./features/settings/settings.component').then((m) => m.SettingsComponent),
  },
  {
    path: 'settings/password',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./features/settings/change-password/change-password.component').then(
        (m) => m.ChangePasswordComponent,
      ),
  },
];
