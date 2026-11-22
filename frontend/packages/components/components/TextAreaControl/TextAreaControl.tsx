import { forwardRef } from 'react';
import type { TextAreaControlProps } from './types/textAreaControl.props';

/** Unstyled native text area for existing forms. */
export const TextAreaControl = forwardRef<HTMLTextAreaElement, TextAreaControlProps>(
  ({ dataTestId, ...props }, ref) => <textarea {...props} ref={ref} data-testid={dataTestId} />,
);
TextAreaControl.displayName = 'TextAreaControl';
