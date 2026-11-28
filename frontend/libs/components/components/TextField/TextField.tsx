import { forwardRef } from 'react';
import { Input } from '../Input';
import type { TextFieldProps } from './types/textField.props';

export const TextField = forwardRef<HTMLInputElement, TextFieldProps>((props, ref) => (
  <Input {...props} ref={ref} />
));
TextField.displayName = 'TextField';
