import type { HTMLAttributes, ReactElement } from 'react';
export interface CardProps extends HTMLAttributes<HTMLElement> { readonly dataTestId?: string; }
export type CardComponent = (props: CardProps) => ReactElement;
