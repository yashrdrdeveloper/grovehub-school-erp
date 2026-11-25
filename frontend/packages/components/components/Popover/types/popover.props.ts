import type { ReactNode } from 'react';

export interface PopoverProps {
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly title: string;
  readonly trigger: ReactNode;
  readonly children: ReactNode;
  readonly width?: number;
  readonly fullWidth?: boolean;
  readonly compact?: boolean;
  readonly dataTestId?: string;
}
