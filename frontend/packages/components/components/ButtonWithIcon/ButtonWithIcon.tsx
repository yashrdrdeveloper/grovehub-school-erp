import { Button } from '../Button';
import { Icon } from '../Icon';

import { BUTTON_ICON_POSITION } from './constants/buttonWithIcon.position';
import type { ButtonWithIconProps } from './types/buttonWithIcon.props';

export const ButtonWithIcon = ({
  icon,
  iconPosition = BUTTON_ICON_POSITION.LEFT,
  iconSize,
  children,
  ...props
}: ButtonWithIconProps) => (
  <Button {...props}>
    {iconPosition === BUTTON_ICON_POSITION.LEFT && <Icon icon={icon} size={iconSize} />}
    {children}
    {iconPosition === BUTTON_ICON_POSITION.RIGHT && <Icon icon={icon} size={iconSize} />}
  </Button>
);
