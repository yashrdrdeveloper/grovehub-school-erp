import { forwardRef } from 'react';
import * as LabelPrimitive from '@radix-ui/react-label';
import type { FieldLabelProps } from './types/fieldLabel.props';

/** Accessible label with native htmlFor and a deterministic test selector. */
export const FieldLabel = forwardRef<HTMLLabelElement, FieldLabelProps>(({ dataTestId, ...props }, ref) =>
  <LabelPrimitive.Root {...props} ref={ref} data-testid={dataTestId}/>,
);
FieldLabel.displayName = 'FieldLabel';
