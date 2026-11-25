import type { ReactNode } from 'react';

export interface TooltipProps {
  readonly label?: string;
  readonly focusable?: boolean;
  readonly fullWidth?: boolean;
  readonly children: ReactNode;
}
