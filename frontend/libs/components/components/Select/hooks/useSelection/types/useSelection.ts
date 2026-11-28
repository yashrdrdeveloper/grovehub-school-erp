import type { SelectValue } from '../../../types/select.props';

export interface UseSelectionProps<T extends SelectValue> {
  readonly value: readonly T[];
  readonly multiple: boolean;
  readonly onChange: (value: readonly T[]) => void;
  readonly maxSelections?: number;
}
