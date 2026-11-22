import type { ForwardRefExoticComponent, InputHTMLAttributes, PropsWithoutRef, RefAttributes } from 'react';

export type InputControlSize = 'sm' | 'md' | 'lg';

export interface InputControlProps extends InputHTMLAttributes<HTMLInputElement> {
  readonly dataTestId?: string;
  readonly unstyled?: boolean;
  readonly controlSize?: InputControlSize;
  readonly invalid?: boolean;
}
export type InputControlComponent = ForwardRefExoticComponent<
  PropsWithoutRef<InputControlProps> & RefAttributes<HTMLInputElement>
>;
