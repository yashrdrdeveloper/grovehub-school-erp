import type { HTMLAttributes, ReactElement } from 'react';
export type TypographyTag = 'h1' | 'h2' | 'h3' | 'h4' | 'p' | 'span' | 'small' | 'strong';
export type TypographyVariant = 'display' | 'heading' | 'subheading' | 'body' | 'caption' | 'label';
export type TypographyTone = 'default' | 'muted' | 'danger' | 'success';
export interface TypographyProps extends HTMLAttributes<HTMLElement> {
  readonly as?: TypographyTag;
  readonly variant?: TypographyVariant;
  readonly tone?: TypographyTone;
}
export type TypographyComponent = (props: TypographyProps) => ReactElement;
