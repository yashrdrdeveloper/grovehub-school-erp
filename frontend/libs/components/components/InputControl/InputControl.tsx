import { INPUT_CONTROL_SIZE } from './constants/inputControl.size';
import { forwardRef } from 'react';
import { cn } from '@Base/utilities/cn';
import type { InputControlComponent, InputControlProps } from './types/inputControl.props';

/** Native leaf input. Unstyled mode preserves existing auth and school-site layouts. */
export const InputControl: InputControlComponent = forwardRef<HTMLInputElement, InputControlProps>(
  (
    {
      id,
      name,
      className,
      dataTestId,
      unstyled = false,
      controlSize = INPUT_CONTROL_SIZE.MEDIUM,
      invalid,
      'aria-invalid': ariaInvalid,
      ...props
    },
    ref,
  ) => (
    <input
      {...props}
      ref={ref}
      id={id}
      name={name}
      data-testid={dataTestId ?? id ?? name}
      aria-invalid={invalid ?? ariaInvalid}
      className={cn(
        !unstyled &&
          'gh-input-control flex w-full min-w-0 rounded-md border border-input bg-background px-3 py-2 text-sm text-foreground shadow-xs placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-0 disabled:cursor-not-allowed disabled:opacity-50',
        !unstyled && `gh-input-control--${controlSize}`,
        className,
      )}
    />
  ),
);
InputControl.displayName = 'InputControl';
