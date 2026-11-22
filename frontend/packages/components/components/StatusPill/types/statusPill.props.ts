import type { ReactElement } from 'react';
export interface StatusPillProps { readonly value: string; }
export type StatusPillTone = 'positive' | 'notice' | 'warning' | 'neutral';
export type StatusPillComponent = (props: StatusPillProps) => ReactElement;
