import type { ReactNode } from 'react';

export interface ButtonContentProps {
  readonly loading: boolean;
  readonly loadingText: string;
  readonly children: ReactNode;
}

export type RenderButtonContent = (props: ButtonContentProps) => ReactNode;
