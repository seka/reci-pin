import { Meta, StoryObj, moduleMetadata } from '@storybook/angular';
import { ConfirmDialogComponent } from './confirm-dialog.component';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { TranslocoModule } from '@jsverse/transloco';

const mockDialogRef = {
  close: () => console.log('dialogRef.close'),
};

const meta: Meta<ConfirmDialogComponent> = {
  title: 'Shared/Molecules/ConfirmDialog',
  component: ConfirmDialogComponent,
  tags: ['autodocs'],
  decorators: [
    moduleMetadata({
      imports: [TranslocoModule],
      providers: [
        { provide: MatDialogRef, useValue: mockDialogRef },
        {
          provide: MAT_DIALOG_DATA,
          useValue: {
            title: 'Confirm Action',
            message: 'Are you sure you want to perform this action?',
            confirmText: 'Confirm',
            cancelText: 'Cancel',
          },
        },
      ],
    }),
  ],
};

export default meta;
type Story = StoryObj<ConfirmDialogComponent>;

export const Default: Story = {
  args: {},
};
