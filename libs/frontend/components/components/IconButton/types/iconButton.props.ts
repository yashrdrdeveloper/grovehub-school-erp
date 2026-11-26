import type { ButtonWithIconProps } from '../../ButtonWithIcon';

export interface IconButtonProps extends Omit<ButtonWithIconProps, 'children'> {
  readonly label: string;
  readonly dataTestId?: string;
}
