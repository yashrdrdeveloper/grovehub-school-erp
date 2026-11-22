import type {
  SelectHTMLAttributes,
  ForwardRefExoticComponent,
  PropsWithoutRef,
  RefAttributes,
} from 'react';
export interface SelectControlProps extends SelectHTMLAttributes<HTMLSelectElement> {
  readonly dataTestId?: string;
}
export type SelectControlComponent = ForwardRefExoticComponent<
  PropsWithoutRef<SelectControlProps> & RefAttributes<HTMLSelectElement>
>;
