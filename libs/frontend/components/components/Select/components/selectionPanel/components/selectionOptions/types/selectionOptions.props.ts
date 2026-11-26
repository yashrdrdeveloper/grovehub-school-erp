import type { SelectValue } from '../../../../../types/select.props';
import type { SelectionPanelProps } from '../../../types/selectionPanel.props';

export type SelectionOptionsProps<T extends SelectValue> = Pick<
  SelectionPanelProps<T>,
  | 'options'
  | 'selected'
  | 'multiple'
  | 'label'
  | 'loading'
  | 'error'
  | 'hasNext'
  | 'dataTestId'
  | 'onToggle'
  | 'onLoadMore'
>;
