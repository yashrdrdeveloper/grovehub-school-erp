import type { SelectCommonProps, SelectValue } from '../../../types/select.props';

export interface SelectionPanelProps<T extends SelectValue> extends SelectCommonProps<T> {
  readonly selected: readonly T[];
  readonly multiple: boolean;
  readonly open: boolean;
  readonly search: string;
  readonly onOpenChange: (open: boolean) => void;
  readonly onSearch: (value: string) => void;
  readonly onToggle: (value: T) => void;
  readonly onClear: () => void;
  readonly onApply: () => void;
  readonly hasNext?: boolean;
  readonly onLoadMore?: () => void;
  readonly onRetry?: () => void;
}
