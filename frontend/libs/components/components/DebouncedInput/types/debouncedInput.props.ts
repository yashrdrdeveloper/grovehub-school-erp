import type { InputProps } from '../../Input';

export interface DebouncedInputProps
  extends Omit<InputProps, 'value' | 'defaultValue' | 'onChange'> {
  readonly value: string;
  readonly debounceMs?: number;
  readonly minLength?: number;
  readonly flushOnBlur?: boolean;
  readonly flushOnEnter?: boolean;
  readonly immediateOnClear?: boolean;
  readonly onChange: (value: string) => void;
  readonly onImmediateChange?: (value: string) => void;
}
