import type { FORM_SPACING } from '../constants/form.spacing';
import type { FORM_LAYOUT } from '../constants/form.layout';
import type { FormHTMLAttributes, ReactElement } from 'react';
export type FormLayout = (typeof FORM_LAYOUT)[keyof typeof FORM_LAYOUT];
export type FormSpacing = (typeof FORM_SPACING)[keyof typeof FORM_SPACING];
export interface FormProps extends FormHTMLAttributes<HTMLFormElement> {
  readonly layout?: FormLayout;
  readonly spacing?: FormSpacing;
  /** Retain existing feature layout during migration without changing native form semantics. */
  readonly unstyled?: boolean;
  readonly dataTestId?: string;
}
export type FormComponent = (props: FormProps) => ReactElement;
