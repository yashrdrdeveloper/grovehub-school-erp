import type { LucideIcon } from '@Components/icons';

export interface DropdownItem {
  readonly id: string;
  readonly label: string;
  readonly icon?: LucideIcon;
  readonly disabled?: boolean;
  readonly destructive?: boolean;
}

export interface DropdownProps {
  readonly label: string;
  readonly items: readonly DropdownItem[];
  readonly onSelect: (id: string) => void;
  readonly disabled?: boolean;
  readonly dataTestId?: string;
}
