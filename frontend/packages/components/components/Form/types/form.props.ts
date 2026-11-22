import type { FormHTMLAttributes, ReactElement } from 'react';
export type FormLayout = 'stack' | 'inline';
export type FormSpacing = 'sm' | 'md' | 'lg';
export interface FormProps extends FormHTMLAttributes<HTMLFormElement> {
  readonly layout?: FormLayout;
  readonly spacing?: FormSpacing;
  readonly dataTestId?: string;
}
export type FormComponent = (props: FormProps) => ReactElement;
