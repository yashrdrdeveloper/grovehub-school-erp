import { SelectionPanel } from '../Select/components/selectionPanel';

import type { SelectionPanelProps } from '../Select/components/selectionPanel';
import type { SelectValue } from '../Select';

export const AsyncSelectView = <T extends SelectValue>(props: SelectionPanelProps<T>) => (
  <SelectionPanel {...props} />
);
