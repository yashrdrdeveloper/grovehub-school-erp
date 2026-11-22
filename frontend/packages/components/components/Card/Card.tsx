import { cn } from '@Base/utilities';
import type { CardComponent } from './types/card.props';

export const Card: CardComponent = ({ className, children, ...props }) => (
  <section className={cn('card', className)} {...props}>{children}</section>
);
