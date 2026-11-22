import type { InputHTMLAttributes } from 'react';
export interface PasswordFieldProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type'> {
  readonly label: string;
  readonly hint?: string;
  readonly error?: string;
  readonly initiallyVisible?: boolean;
}
