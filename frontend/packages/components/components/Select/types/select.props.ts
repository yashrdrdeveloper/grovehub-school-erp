export type SelectValue = string | number;

export interface SelectOption<T extends SelectValue = string> {
  readonly value: T;
  readonly label: string;
  readonly description?: string;
  readonly disabled?: boolean;
}

export interface SelectCommonProps<T extends SelectValue> {
  readonly options: readonly SelectOption<T>[];
  readonly label: string;
  readonly hideLabel?: boolean;
  readonly placeholder?: string;
  readonly disabled?: boolean;
  readonly searchable?: boolean;
  readonly clearable?: boolean;
  readonly loading?: boolean;
  readonly error?: string;
  readonly dataTestId?: string;
}

export interface SelectProps<T extends SelectValue = string> extends SelectCommonProps<T> {
  readonly value: T | null;
  readonly onChange: (value: T | null) => void;
}
