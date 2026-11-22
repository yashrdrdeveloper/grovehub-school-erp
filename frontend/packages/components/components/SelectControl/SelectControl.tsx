import { forwardRef } from 'react';
import type { SelectControlProps } from './types/selectControl.props';

/** Unstyled native select for existing school forms. */
export const SelectControl = forwardRef<HTMLSelectElement, SelectControlProps>(({ dataTestId, ...props }, ref) =>
  <select {...props} ref={ref} data-testid={dataTestId}/>,
);
SelectControl.displayName = 'SelectControl';
