import type { SelectCommonProps, SelectValue } from '../../Select';

export interface MultiSelectProps<T extends SelectValue = string> extends SelectCommonProps<T> {
  readonly dataTestId?: string;
  readonly value: readonly T[];
  readonly maxSelections?: number;
  readonly onChange: (value: readonly T[]) => void;
}
