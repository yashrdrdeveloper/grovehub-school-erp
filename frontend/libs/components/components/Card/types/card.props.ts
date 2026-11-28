import type { CARD_PADDING } from '../constants/card.padding';
import type { CARD_VARIANT } from '../constants/card.variant';
import type { HTMLAttributes, ReactElement } from 'react';
export type CardVariant = (typeof CARD_VARIANT)[keyof typeof CARD_VARIANT];
export type CardPadding = (typeof CARD_PADDING)[keyof typeof CARD_PADDING];
export interface CardProps extends HTMLAttributes<HTMLElement> {
  readonly dataTestId?: string;
  readonly variant?: CardVariant;
  readonly padding?: CardPadding;
}
export type CardComponent = (props: CardProps) => ReactElement;
