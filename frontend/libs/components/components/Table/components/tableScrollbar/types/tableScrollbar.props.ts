import type { RefObject } from 'react';

export interface TableScrollbarProps {
  readonly viewport: RefObject<HTMLElement | null>;
  readonly viewportId: string;
  readonly axis: 'vertical' | 'horizontal';
}
