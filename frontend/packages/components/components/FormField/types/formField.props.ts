import type { HTMLAttributes, ReactElement, ReactNode } from 'react';
export interface FormFieldProps extends Omit<HTMLAttributes<HTMLDivElement>, 'children'> {
  readonly htmlFor: string;
  readonly label: string;
  readonly hint?: string;
  readonly error?: string;
  readonly required?: boolean;
  readonly children: ReactNode;
}
export type FormFieldComponent = (props: FormFieldProps) => ReactElement;
