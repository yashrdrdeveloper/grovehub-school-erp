import type { ReactElement } from 'react';
import type { LucideIcon } from '@Components/icons';
export interface EmptyStateProps {
  readonly title: string;
  readonly description: string;
  readonly icon?: LucideIcon;
  readonly compact?: boolean;
  readonly dataTestId?: string;
}
export type EmptyStateComponent = (props: EmptyStateProps) => ReactElement;
