import type { ForwardRefExoticComponent, InputHTMLAttributes, PropsWithoutRef, RefAttributes } from 'react';

export interface InputControlProps extends InputHTMLAttributes<HTMLInputElement> {
  readonly dataTestId?: string;
  readonly unstyled?: boolean;
  readonly invalid?: boolean;
}
export type InputControlComponent = ForwardRefExoticComponent<
  PropsWithoutRef<InputControlProps> & RefAttributes<HTMLInputElement>
>;
