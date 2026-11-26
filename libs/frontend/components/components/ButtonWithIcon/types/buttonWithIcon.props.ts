import type { ReactNode } from 'react';
import type { LucideIcon } from '@Components/icons';
import type { ButtonProps } from '../../Button';
import type { BUTTON_ICON_POSITION } from '../constants/buttonWithIcon.position';

export interface ButtonWithIconProps extends Omit<ButtonProps, 'children'> {
  readonly dataTestId?: string;
  readonly icon: LucideIcon;
  readonly iconPosition?: (typeof BUTTON_ICON_POSITION)[keyof typeof BUTTON_ICON_POSITION];
  readonly iconSize?: number;
  readonly children?: ReactNode;
}
