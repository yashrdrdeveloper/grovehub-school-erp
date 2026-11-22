import { forwardRef } from 'react';
import { cn } from '@Base/utilities/cn';
import type { InputControlComponent, InputControlProps } from './types/inputControl.props';

/** Native leaf input. Unstyled mode preserves existing auth and school-site layouts. */
export const InputControl: InputControlComponent = forwardRef<HTMLInputElement, InputControlProps>(({
  id, name, className, dataTestId, unstyled = false, invalid, 'aria-invalid': ariaInvalid, ...props
}, ref) => (
  <input {...props} ref={ref} id={id} name={name}
    data-testid={dataTestId ?? id ?? name ?? 'ui-input-control'}
    aria-invalid={invalid ?? ariaInvalid}
    className={cn(!unstyled && 'flex h-10 w-full min-w-0 rounded-md border border-input bg-background px-3 py-2 text-sm text-foreground shadow-xs placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50', className)}/>
));
InputControl.displayName = 'InputControl';
