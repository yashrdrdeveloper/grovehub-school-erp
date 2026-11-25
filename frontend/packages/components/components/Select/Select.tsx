import { SelectionPanel } from './components/selectionPanel';

import type { SelectValue } from './types/select.props';
import type { SelectionPanelProps } from './components/selectionPanel';

export const SelectView = <T extends SelectValue>(props: SelectionPanelProps<T>) => (
  <SelectionPanel {...props} />
);
