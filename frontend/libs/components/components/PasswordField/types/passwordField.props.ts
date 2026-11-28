import type { InputHTMLAttributes } from 'react';
import type { InputControlSize } from '../../InputControl/types/inputControl.props';
export interface PasswordFieldProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type'> {
  readonly label: string;
  readonly hint?: string;
  readonly controlSize?: InputControlSize;
  readonly error?: string;
  readonly initiallyVisible?: boolean;
  readonly dataTestId?: string;
}
