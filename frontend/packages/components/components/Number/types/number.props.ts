import type { InputProps } from '../../Input';

export interface NumberProps extends Omit<InputProps, 'type' | 'value' | 'defaultValue' | 'onChange'> {
  readonly value: number | null;
  readonly onChange: (value: number | null) => void;
}
