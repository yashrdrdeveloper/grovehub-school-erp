import { _filter } from '@Base/utilities/filter';

import { Box } from '../../../Box';
import { Dropdown } from '../../../Dropdown';

import type { TableRowActionsProps } from './types/tableRowActions.props';

export const TableRowActions = <T,>({ row, props }: TableRowActionsProps<T>) => {
  const items = _filter(props.rowActions ?? [], (action) => !action.isVisible || action.isVisible(row));
  if (!items.length) return null;
  return (
    <Box onClick={(event) => event.stopPropagation()}>
      <Dropdown label={__('Actions')} items={items} onSelect={(id) => props.actions?.ROW_ACTION?.({ id, row })} />
    </Box>
  );
};
