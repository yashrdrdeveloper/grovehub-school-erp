import type { ReactElement, ReactNode } from 'react';

export interface DialogProps {
  readonly open: boolean;
  readonly onOpenChange: (value: boolean) => void;
  readonly title: string;
  readonly description?: string;
  readonly children: ReactNode;
  readonly className?: string;
}
export type DialogComponent = (props: DialogProps) => ReactElement;
