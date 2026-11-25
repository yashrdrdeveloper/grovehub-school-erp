import type { IconProps } from './types/icon.props';

export const Icon = ({ icon: Glyph, size = 16, label, className, dataTestId }: IconProps) => (
  <Glyph
    size={size}
    aria-label={label}
    aria-hidden={!label}
    className={className}
    data-testid={dataTestId}
  />
);
