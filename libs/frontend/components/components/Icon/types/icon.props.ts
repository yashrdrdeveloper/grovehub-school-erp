import type { LucideIcon } from '@Components/icons';

export interface IconProps {
  readonly icon: LucideIcon;
  readonly size?: number;
  readonly label?: string;
  readonly className?: string;
  readonly dataTestId?: string;
}
