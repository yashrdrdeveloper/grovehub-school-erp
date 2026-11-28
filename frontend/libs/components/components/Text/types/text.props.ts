import type { ReactElement } from 'react';
import type { TEXT_TAG } from '../constants/text.tag';
import type { TypographyProps } from '../../Typography/types/typography.props';
export interface TextProps extends Omit<TypographyProps, 'as'> {
  readonly as?: (typeof TEXT_TAG)[keyof typeof TEXT_TAG];
  readonly dataTestId?: string;
}
export type TextComponent = (props: TextProps) => ReactElement;
