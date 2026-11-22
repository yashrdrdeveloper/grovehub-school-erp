import type { ForwardRefExoticComponent, HTMLAttributes, PropsWithoutRef, RefAttributes } from 'react';

export type BoxTag =
  | 'div' | 'span' | 'p' | 'small' | 'strong' | 'figure' | 'figcaption' | 'em' | 'i' | 'time'
  | 'main' | 'section' | 'article' | 'aside' | 'header' | 'footer' | 'nav'
  | 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6';
export interface BoxProps extends HTMLAttributes<HTMLElement> {
  readonly as?: BoxTag;
  readonly dataTestId?: string;
}
export type BoxComponent = ForwardRefExoticComponent<PropsWithoutRef<BoxProps> & RefAttributes<HTMLElement>>;
