import type { ForwardRefExoticComponent, InputHTMLAttributes, PropsWithoutRef, RefAttributes } from 'react';

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  readonly label: string;
  readonly error?: string;
  readonly hint?: string;
}
export type InputComponent = ForwardRefExoticComponent<
  PropsWithoutRef<InputProps> & RefAttributes<HTMLInputElement>
>;
