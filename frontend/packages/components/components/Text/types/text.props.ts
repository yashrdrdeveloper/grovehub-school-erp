import type { ReactElement } from 'react';
import type { TypographyProps } from '../../Typography/types/typography.props';
export interface TextProps extends Omit<TypographyProps, 'as'> {
  readonly as?: 'span' | 'p' | 'small' | 'strong';
  readonly dataTestId?: string;
}
export type TextComponent = (props: TextProps) => ReactElement;
