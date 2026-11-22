import type { FormHTMLAttributes, ReactElement } from 'react';
export type FormLayout = 'stack' | 'inline';
export type FormSpacing = 'sm' | 'md' | 'lg';
export interface FormProps extends FormHTMLAttributes<HTMLFormElement> {
  readonly layout?: FormLayout;
  readonly spacing?: FormSpacing;
  /** Retain existing feature layout during migration without changing native form semantics. */
  readonly unstyled?: boolean;
  readonly dataTestId?: string;
}
export type FormComponent = (props: FormProps) => ReactElement;
