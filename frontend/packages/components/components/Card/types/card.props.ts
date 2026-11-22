import type { HTMLAttributes, ReactElement } from 'react';
export type CardVariant = 'default' | 'outline' | 'muted';
export type CardPadding = 'none' | 'sm' | 'md' | 'lg';
export interface CardProps extends HTMLAttributes<HTMLElement> {
  readonly dataTestId?: string;
  readonly variant?: CardVariant;
  readonly padding?: CardPadding;
}
export type CardComponent = (props: CardProps) => ReactElement;
