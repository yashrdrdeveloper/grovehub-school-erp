import type { HTMLAttributes } from 'react';

import { cn } from '@Base/utilities';

export const Card = ({ className, children, ...props }: HTMLAttributes<HTMLElement>) => (
  <section className={cn('card', className)} {...props}>{children}</section>
);
