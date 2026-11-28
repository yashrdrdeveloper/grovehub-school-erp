import { SelectionPanel } from '../Select/components/selectionPanel';

import type { SelectionPanelProps } from '../Select/components/selectionPanel';
import type { SelectValue } from '../Select';

export const MultiSelectView = <T extends SelectValue>({
  dataTestId,
  ...props
}: SelectionPanelProps<T>) => <SelectionPanel {...props} dataTestId={dataTestId} />;
