import type { ComponentPropsWithoutRef, ForwardRefExoticComponent, PropsWithoutRef, RefAttributes } from 'react';
import * as LabelPrimitive from '@radix-ui/react-label';
export interface FieldLabelProps extends ComponentPropsWithoutRef<typeof LabelPrimitive.Root> {
  readonly dataTestId?: string;
}
export type FieldLabelComponent = ForwardRefExoticComponent<PropsWithoutRef<FieldLabelProps> & RefAttributes<HTMLLabelElement>>;
