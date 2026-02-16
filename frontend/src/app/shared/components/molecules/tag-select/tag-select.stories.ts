import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { moduleMetadata, Meta, StoryObj } from '@storybook/angular';
import { TagSelectComponent } from './tag-select.component';

const meta: Meta<TagSelectComponent> = {
  title: 'Molecules/TagSelect',
  component: TagSelectComponent,
  tags: ['autodocs'],
  decorators: [
    moduleMetadata({
      imports: [BrowserAnimationsModule],
    }),
  ],
  argTypes: {
  },
};

export default meta;
type Story = StoryObj<TagSelectComponent>;

const mockTags = [
  { id: 1, name: '和食' },
  { id: 2, name: '洋食' },
  { id: 3, name: '中華' },
  { id: 4, name: 'デザート' },
  { id: 5, name: '時短' },
];

export const Default: Story = {
  args: {
    tags: mockTags,
    selectedTagIds: [],
  },
};

export const WithPreSelection: Story = {
  args: {
    tags: mockTags,
    selectedTagIds: [1, 5],
  },
};

export const Empty: Story = {
  args: {
    tags: [],
    selectedTagIds: [],
  },
};
