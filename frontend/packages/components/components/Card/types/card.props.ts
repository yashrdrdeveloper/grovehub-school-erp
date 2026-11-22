import type { HTMLAttributes, ReactElement } from 'react';
export type CardProps = HTMLAttributes<HTMLElement>;
export type CardComponent = (props: CardProps) => ReactElement;
