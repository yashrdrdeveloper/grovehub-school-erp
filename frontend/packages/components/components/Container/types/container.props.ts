import type { HTMLAttributes, ReactElement } from 'react';
export type ContainerTag = 'div' | 'main' | 'section' | 'article';
export type ContainerWidth = 'sm' | 'md' | 'lg' | 'xl' | 'full';
export interface ContainerProps extends HTMLAttributes<HTMLElement> {
  readonly as?: ContainerTag;
  readonly maxWidth?: ContainerWidth;
  readonly padded?: boolean;
  readonly dataTestId?: string;
}
export type ContainerComponent = (props: ContainerProps) => ReactElement;
