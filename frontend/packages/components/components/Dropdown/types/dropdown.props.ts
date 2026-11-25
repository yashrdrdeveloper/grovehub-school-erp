import type { LucideIcon } from '@Components/icons';
import type { DropdownItemType } from '../constants/dropdown.type';

export interface DropdownItem {
  readonly id: string;
  readonly label: string;
  readonly type?: (typeof DropdownItemType)[keyof typeof DropdownItemType];
  readonly icon?: LucideIcon;
  readonly disabled?: boolean;
  readonly destructive?: boolean;
  readonly checked?: boolean;
  readonly children?: readonly DropdownItem[];
}

export interface DropdownProps {
  readonly label: string;
  readonly items: readonly DropdownItem[];
  readonly onSelect: (id: string) => void;
  readonly disabled?: boolean;
  readonly dataTestId?: string;
}
