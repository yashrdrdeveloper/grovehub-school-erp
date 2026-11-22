import { CARD_VARIANT } from './constants/card.variant';
import { CARD_PADDING } from './constants/card.padding';
import { cn } from '@Base/utilities/cn';
import type { CardComponent } from './types/card.props';

export const Card: CardComponent = ({ className, children, dataTestId, variant = CARD_VARIANT.DEFAULT, padding = CARD_PADDING.MEDIUM, ...props }) => (
  <section className={cn('card gh-card rounded-xl border border-border bg-card text-card-foreground shadow-sm', `gh-card--${variant}`, `gh-card--padding-${padding}`, className)} data-testid={dataTestId} {...props}>{children}</section>
);
