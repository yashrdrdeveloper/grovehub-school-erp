import { cn } from '@Base/utilities/cn';
import type { CardComponent } from './types/card.props';

export const Card: CardComponent = ({ className, children, dataTestId, ...props }) => (
  <section className={cn('card rounded-xl border border-border bg-card text-card-foreground shadow-sm', className)} data-testid={dataTestId} {...props}>{children}</section>
);
