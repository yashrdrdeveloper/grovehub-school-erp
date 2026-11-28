import type {
  ForwardRefExoticComponent,
  InputHTMLAttributes,
  PropsWithoutRef,
  RefAttributes,
} from 'react';
import type { InputControlSize } from '../../InputControl/types/inputControl.props';

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  readonly label: string;
  readonly error?: string;
  readonly hint?: string;
  readonly controlSize?: InputControlSize;
  readonly dataTestId?: string;
}
export type InputComponent = ForwardRefExoticComponent<
  PropsWithoutRef<InputProps> & RefAttributes<HTMLInputElement>
>;
