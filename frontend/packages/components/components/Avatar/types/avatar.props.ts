import type { ReactNode } from 'react';

export interface AvatarProps {
  readonly name: string;
  readonly src?: string;
  readonly fallback?: ReactNode;
  readonly dataTestId?: string;
  /** Compact avatar for dense list rows. */
  readonly compact?: boolean;
}
