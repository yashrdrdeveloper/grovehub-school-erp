import type { ReactElement } from 'react';
export interface EmptyStateProps {
  readonly title: string;
  readonly description: string; readonly dataTestId?: string;
}
export type EmptyStateComponent = (props: EmptyStateProps) => ReactElement;
