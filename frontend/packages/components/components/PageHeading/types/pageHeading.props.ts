import type { ReactElement, ReactNode } from 'react';
export interface PageHeadingProps {
  readonly eyebrow?: string;
  readonly title: string;
  readonly description?: string;
  readonly action?: ReactNode;
  readonly dataTestId?: string;
}
export type PageHeadingComponent = (props: PageHeadingProps) => ReactElement;
