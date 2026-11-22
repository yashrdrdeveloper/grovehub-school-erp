import type { ReactElement, ReactNode } from 'react';
export interface PageHeadingProps {
  readonly eyebrow?: string;
  readonly title: string;
  readonly description?: string;
  readonly action?: ReactNode;
}
export type PageHeadingComponent = (props: PageHeadingProps) => ReactElement;
