import type { INPUT_CONTROL_SIZE } from '../constants/inputControl.size';
import type { ForwardRefExoticComponent, InputHTMLAttributes, PropsWithoutRef, RefAttributes } from 'react';

export type InputControlSize = typeof INPUT_CONTROL_SIZE[keyof typeof INPUT_CONTROL_SIZE];

export interface InputControlProps extends InputHTMLAttributes<HTMLInputElement> {
  readonly dataTestId?: string;
  readonly unstyled?: boolean;
  readonly controlSize?: InputControlSize;
  readonly invalid?: boolean;
}
export type InputControlComponent = ForwardRefExoticComponent<
  PropsWithoutRef<InputControlProps> & RefAttributes<HTMLInputElement>
>;
