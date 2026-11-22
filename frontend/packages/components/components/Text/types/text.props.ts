import type { ReactElement } from 'react';
import type { TYPOGRAPHY_TAG } from '../../Typography/constants/typography.tag';
import type { TypographyProps } from '../../Typography/types/typography.props';
export interface TextProps extends Omit<TypographyProps, 'as'> {
  readonly as?: (typeof TYPOGRAPHY_TAG)['SPAN' | 'P' | 'SMALL' | 'STRONG'];
  readonly dataTestId?: string;
}
export type TextComponent = (props: TextProps) => ReactElement;
