import type { InputProps } from '../../Input';

export interface DateTimeProps extends Omit<InputProps, 'type' | 'value' | 'defaultValue' | 'onChange'> {
  readonly value: string | null;
  readonly onChange: (value: string | null) => void;
}
